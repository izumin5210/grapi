package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"go/ast"
	"go/format"
	"go/parser"
	"go/token"
	"io/ioutil"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/bradleyjkemp/cupaloy/v2"
	"golang.org/x/tools/go/packages/packagestest"
)

var (
	grapiCmd = flag.String("grapi", "grapi", "path of grapi command")
	revision = flag.String("revision", "", "target revision")
)

func TestE2E_withModules(t *testing.T) {
	invokeE2ETest(t, packagestest.Modules)
}

func TestE2E_withDep(t *testing.T) {
	t.SkipNow()
	invokeE2ETest(t, packagestest.GOPATH)
}

func invokeE2ETest(t *testing.T, exporter packagestest.Exporter) {
	t.Helper()

	exported := packagestest.Export(t, exporter, []packagestest.Module{
		{Name: "sampleapp", Files: map[string]interface{}{".keep": ""}},
	})
	defer exported.Cleanup()

	rootPath := exported.Config.Dir
	exported.Config.Dir = filepath.Dir(rootPath)
	checkNoErr(t, os.RemoveAll(rootPath))

	// init
	{
		args := []string{"--debug", "init", "--package", "testuser.sampleapp"}
		if *revision != "" {
			args = append(args, "--revision="+*revision)
		} else {
			args = append(args, "--HEAD")
		}
		if exporter.Name() == "GOPATH" {
			args = append(args, "--use-dep")
		}
		args = append(args, filepath.Base(rootPath))
		invoke(t, exported, exec.Command(*grapiCmd, args...))
		checkExistence(t, rootPath)
		t.Log("Initialize a project successfully")
	}

	exported.Config.Dir = rootPath

	ignoreFiles := map[string]struct{}{
		"/go.mod": struct{}{},
		"/go.sum": struct{}{},
	}

	err := filepath.Walk(rootPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if info.IsDir() {
			return nil
		}
		rel := strings.TrimPrefix(path, rootPath)
		if _, ok := ignoreFiles[rel]; ok {
			return nil
		}
		if strings.HasPrefix(rel, "/bin/") {
			return nil
		}

		t.Run(rel, func(t *testing.T) {
			data, err := ioutil.ReadFile(path)
			if err != nil {
				t.Errorf("failed to open %s: %v", path, err)
			}

			cupaloy.SnapshotT(t, string(data))
		})

		return nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// generate service
	{
		invoke(t, exported, exec.Command(*grapiCmd, "--debug", "g", "service", "book", "list"))
		checkExistence(t, filepath.Join(rootPath, "app", "server", "book_server.go"))
		t.Log("Generate a service successfully")
	}

	port := getFreePort(t)
	updateRun(t, rootPath, port)
	updateServerImpl(t, rootPath)

	// run server
	{
		t.Log("Start the server")
		cmd := exec.Command(*grapiCmd, "--debug", "server")
		sdCh := make(chan struct{}, 1)
		go func() {
			defer close(sdCh)
			invoke(t, exported, cmd)
		}()

		startedAt := time.Now()
		var (
			resp     *http.Response
			retryCnt int
			err      error
		)

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

		sendSignal(t, cmd, os.Interrupt)
		toCtx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		select {
		case <-sdCh:
			t.Log("Shutdown server successfully")
		case <-toCtx.Done():
			t.Log("Deadline exceeded stopping server")
			sendSignal(t, cmd, os.Kill)
			<-sdCh
		}
	}
}

func invoke(t *testing.T, exported *packagestest.Exported, cmd *exec.Cmd) {
	t.Helper()
	cmd.Dir = exported.Config.Dir
	for _, kv := range exported.Config.Env {
		if strings.HasPrefix(kv, "GOPROXY=") {
			continue
		}
		cmd.Env = append(cmd.Env, kv)
	}
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Run()
	if err != nil {
		t.Fatalf("failed to execute command %v: %v", cmd, err)
	}
}

func checkNoErr(t *testing.T, err error) {
	t.Helper()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func checkExistence(t *testing.T, path string) {
	t.Helper()
	_, err := os.Stat(path)
	if err != nil {
		if os.IsNotExist(err) {
			t.Fatalf("%s does not exist: %v", path, err)
		}
		t.Fatalf("failed to check file existence: %v", err)
	}
}

func getFreePort(t *testing.T) int {
	t.Helper()
	lis, err := net.Listen("tcp", ":0")
	checkNoErr(t, err)
	defer lis.Close()

	return lis.Addr().(*net.TCPAddr).Port
}

func sendSignal(t *testing.T, cmd *exec.Cmd, sig os.Signal) {
	t.Helper()
	checkNoErr(t, cmd.Process.Signal(sig))
}

type visitor struct {
	VisitFunc func(ast.Visitor, ast.Node) ast.Visitor
}

func (v *visitor) Visit(node ast.Node) ast.Visitor {
	return v.VisitFunc(v, node)
}

func updateRun(t *testing.T, rootPath string, port int) {
	data, err := ioutil.ReadFile(filepath.Join(rootPath, "cmd", "server", "run.go"))
	checkNoErr(t, err)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", data, parser.DeclarationErrors)
	checkNoErr(t, err)

	ast.Walk(&visitor{
		VisitFunc: func(v ast.Visitor, n ast.Node) ast.Visitor {
			switch n := n.(type) {
			case *ast.GenDecl:
				if n.Tok == token.IMPORT {
					n.Specs = append(n.Specs, &ast.ImportSpec{
						Path: &ast.BasicLit{
							Kind:  token.STRING,
							Value: strconv.Quote("sampleapp/app/server"),
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
	checkNoErr(t, err)
	err = ioutil.WriteFile(filepath.Join(rootPath, "cmd", "server", "run.go"), buf.Bytes(), 0755)
	checkNoErr(t, err)
}

func updateServerImpl(t *testing.T, rootPath string) {
	data, err := ioutil.ReadFile(filepath.Join(rootPath, "app", "server", "book_server.go"))
	checkNoErr(t, err)

	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", data, parser.DeclarationErrors)
	checkNoErr(t, err)

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
								Value: strconv.Quote("sampleapp/api"),
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
	checkNoErr(t, err)
	err = ioutil.WriteFile(filepath.Join(rootPath, "app", "server", "book_server.go"), buf.Bytes(), 0755)
	checkNoErr(t, err)
}
