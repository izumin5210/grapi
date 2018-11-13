package main

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"testing"
	"time"
)

func Test_Integration(t *testing.T) {
	wd, err := os.Getwd()
	orDie(t, err)

	srcDir := filepath.Dir(wd)

	name := "sample"
	rootPath := filepath.Join(srcDir, name)

	// defer os.RemoveAll(rootPath)

	args := []string{"--debug", "init"}
	if commit, ok := os.LookupEnv("TARGET_REVISION"); ok {
		args = append(args, "--revision="+commit)
	} else {
		args = append(args, "--HEAD")
	}
	args = append(args, name)

	run(t, srcDir, exec.Command("grapi", args...))

	if !exists(t, rootPath) {
		t.Fatalf("%s does not exist: %v", rootPath, err)
	}
	t.Log("Initialize a project successfully")

	run(t, rootPath, exec.Command("grapi", "--debug", "g", "service", "book", "list"))

	svrPath := filepath.Join(rootPath, "app", "server", "book_server.go")
	if !exists(t, svrPath) {
		t.Fatalf("%s does not exist: %v", svrPath, err)
	}
	t.Log("Generate a service successfully")

	port := 15261

	updateRun(t, rootPath, port)
	updateServerImpl(t, rootPath)

	t.Log("Start the server")
	svrCtx, cancel := context.WithCancel(context.Background())
	cmd := exec.CommandContext(svrCtx, "grapi", "--debug", "server")
	cmd.Dir = rootPath
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err = cmd.Start()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	defer func() {
		if cmd.Process != nil && cmd.ProcessState != nil && !cmd.ProcessState.Exited() {
			cmd.Process.Kill()
		}
	}()

	startedAt := time.Now()
	var resp *http.Response
	var retryCnt int

	for {
		func() {
			defer recover()
			resp, err = http.Get(fmt.Sprintf("http://localhost:%d/books", port))
		}()
		if err != nil && time.Since(startedAt) < 120*time.Second {
			time.Sleep(5 * time.Second)
			retryCnt++
		} else {
			break
		}
	}

	if err != nil {
		t.Fatalf("Unexpected error (retry count: %d): %v", retryCnt, err)
	}

	if got, want := resp.StatusCode, 200; got != want {
		t.Errorf("Response status is %d, want %d", got, want)
	}

	t.Log("HTTP Request successfully")

	cancel()
	toCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	select {
	case <-svrCtx.Done():
		t.Log("Shutdown server successfully")
	case <-toCtx.Done():
		t.Log("Deadline exceeded stopping server")
		cmd.Process.Signal(os.Kill)
	}
	err = cmd.Wait()
}

type visitor struct {
	VisitFunc func(ast.Visitor, ast.Node) ast.Visitor
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	return v.VisitFunc(v, node)
}

func updateRun(t *testing.T, rootPath string, port int) {
	data, err := ioutil.ReadFile(filepath.Join(rootPath, "app", "run.go"))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", data, parser.DeclarationErrors)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	ast.Walk(&visitor{
		VisitFunc: func(v ast.Visitor, n ast.Node) ast.Visitor {
			switch n := n.(type) {
			case *ast.GenDecl:
				if n.Tok == token.IMPORT {
					n.Specs = append(n.Specs, &ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: strconv.Quote("sample/app/server"),
						},
					})
				}
			case *ast.CallExpr:
				switch fun := n.Fun.(type) {
				case *ast.SelectorExpr:
					switch fun.Sel.Name {
					case "New":
						n.Args = append(n.Args, &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("grapiserver"),
								Sel: ast.NewIdent("WithGatewayAddr"),
							},
							Args: []ast.Expr{
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: strconv.Quote("tcp"),
								},
								&ast.BasicLit{
									Kind:  token.STRING,
									Value: strconv.Quote(fmt.Sprintf(":%d", port)),
								},
							},
						})
					case "WithServers":
						n.Args = append(n.Args, &ast.CallExpr{
							Fun: &ast.SelectorExpr{
								X:   ast.NewIdent("server"),
								Sel: ast.NewIdent("NewBookServiceServer"),
							},
						})
					}
				}
			}
			return v
		},
	}, f)

	buf := new(bytes.Buffer)
	err = format.Node(buf, token.NewFileSet(), f)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	err = ioutil.WriteFile(filepath.Join(rootPath, "app", "run.go"), buf.Bytes(), 0755)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func updateServerImpl(t *testing.T, rootPath string) {
	data, err := ioutil.ReadFile(filepath.Join(rootPath, "app", "server", "book_server.go"))
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", data, parser.DeclarationErrors)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}

	ast.Walk(&visitor{
		VisitFunc: func(v ast.Visitor, n ast.Node) ast.Visitor {
			switch n := n.(type) {
			case *ast.GenDecl:
				if n.Tok == token.IMPORT {
					n.Specs = []ast.Spec{
						&ast.ImportSpec{
							Path: &ast.BasicLit{
								Kind:  token.STRING,
								Value: strconv.Quote("context"),
							},
						},
						&ast.ImportSpec{
							Path: &ast.BasicLit{
								Kind:  token.STRING,
								Value: strconv.Quote("github.com/izumin5210/grapi/pkg/grapiserver"),
							},
						},
						&ast.ImportSpec{
							Name: &ast.Ident{Name: "api_pb"},
							Path: &ast.BasicLit{
								Kind:  token.STRING,
								Value: strconv.Quote("sample/api"),
							},
						},
					}
				}
			case *ast.FuncDecl:
				if n.Name.Name == "ListBooks" {
					n.Body.List = []ast.Stmt{
						&ast.ReturnStmt{
							Results: []ast.Expr{
								&ast.UnaryExpr{
									X: &ast.CompositeLit{
										Type: &ast.SelectorExpr{
											X:   ast.NewIdent("api_pb"),
											Sel: ast.NewIdent("ListBooksResponse"),
										},
									},
									Op: token.AND,
								},
								&ast.Ident{Name: "nil"},
							},
						},
					}
				}
			}
			return v
		},
	}, f)

	buf := new(bytes.Buffer)
	err = format.Node(buf, token.NewFileSet(), f)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
	err = ioutil.WriteFile(filepath.Join(rootPath, "app", "server", "book_server.go"), buf.Bytes(), 0755)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func run(t *testing.T, dir string, cmd *exec.Cmd) {
	t.Helper()
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf("failed to execute command %v: %v", cmd, err)
	}
}

func orDie(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func exists(t *testing.T, path string) bool {
	t.Helper()
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			return false
		}
		t.Fatalf("failed to check file existence: %v", err)
	}
	return true
}
