package main

import (
	"bytes"
	"context"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"testing"
	"time"

	"github.com/spf13/afero"
)

func Test_Integration(t *testing.T) {
	_, testfilepath, _, _ := runtime.Caller(0)
	wd := filepath.Dir(testfilepath)
	bin := filepath.Join(wd, "..", "..", "bin", "grapi")
	gopath := filepath.Join(wd, "go")
	srcDir := filepath.Join(gopath, "src")

	fs := afero.NewOsFs()
	name := "sample"
	rootPath := filepath.Join(srcDir, name)

	fs.MkdirAll(srcDir, 0755)
	defer fs.RemoveAll(gopath)

	cmd := exec.Command(bin, "--debug", "init", "--HEAD", name)
	cmd.Dir = srcDir
	cmd.Env = append(os.Environ(), "GOPATH="+gopath)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to initialize project: %v\n%s", err, string(out))
	}

	if ok, err := afero.DirExists(fs, rootPath); err != nil || !ok {
		t.Fatalf("%s does not exist: %v", rootPath, err)
	}
	t.Log("Initialize a project successfully")

	cmd = exec.Command(bin, "--debug", "g", "service", "book", "list")
	cmd.Dir = rootPath
	cmd.Env = append(os.Environ(), "GOPATH="+gopath)
	out, err = cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("failed to generate service: %v\n%s", err, string(out))
	}

	svrPath := filepath.Join(rootPath, "app", "server", "book_server.go")
	if ok, err := afero.Exists(fs, svrPath); err != nil || !ok {
		t.Fatalf("%s does not exist: %v", svrPath, err)
	}
	t.Log("Generate a service successfully")

	port := 15261

	updateRun(t, fs, rootPath, port)
	updateServerImpl(t, fs, rootPath)

	t.Log("Start the server")
	svrCtx, cancel := context.WithCancel(context.Background())
	cmd = exec.CommandContext(svrCtx, bin, "--debug", "server")
	cmd.Dir = rootPath
	cmd.Env = append(os.Environ(), "GOPATH="+gopath)
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

func updateRun(t *testing.T, fs afero.Fs, rootPath string, port int) {
	data, err := afero.ReadFile(fs, filepath.Join(rootPath, "app", "run.go"))
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
	err = afero.WriteFile(fs, filepath.Join(rootPath, "app", "run.go"), buf.Bytes(), 0755)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}

func updateServerImpl(t *testing.T, fs afero.Fs, rootPath string) {
	data, err := afero.ReadFile(fs, filepath.Join(rootPath, "app", "server", "book_server.go"))
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
								Value: strconv.Quote("sample/api_pb"),
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
	err = afero.WriteFile(fs, filepath.Join(rootPath, "app", "server", "book_server.go"), buf.Bytes(), 0755)
	if err != nil {
		t.Fatalf("Unexpected error: %v", err)
	}
}
