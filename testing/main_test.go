package testing

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	"google.golang.org/grpc"

	"github.com/izumin5210/grapi/pkg/grapiserver"
	"github.com/izumin5210/grapi/testing/api"
	"github.com/izumin5210/grapi/testing/app/server"
)

func Test_server_onlyGateway(t *testing.T) {
	var port int64 = 15261
	app, err := grapiserver.New().
		SetGatewayAddr("tcp", ":"+strconv.FormatInt(port, 10)).
		AddServers(
			server.NewLibraryServiceServer(),
		).
		Build()

	if err != nil {
		t.Fatalf("failed to build grapserver.Engine: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.Serve(); err != nil {
			t.Errorf("Engine.Serve returned an error: %v", err)
		}
	}()

	time.Sleep(2)
	resp, err := http.Get(fmt.Sprintf("http://localhost:%d/books", port))

	if err != nil {
		t.Fatalf("failed to fetch book resources: %v", err)
	}
	defer resp.Body.Close()

	if got, want := resp.StatusCode, 200; got != want {
		t.Fatalf("Response status is %d, want %d", got, want)
	}

	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("failed to read response body: %v", err)
	}

	got := map[string]interface{}{}
	err = json.Unmarshal(data, &got)
	if err != nil {
		t.Fatalf("failed to parse response body: %v", err)
	}
	want := map[string]interface{}{
		"books": []interface{}{
			map[string]interface{}{"book_id": "The Go Programming Language"},
			map[string]interface{}{"book_id": "Programming Ruby"},
		},
	}

	if diff := cmp.Diff(got, want); diff != "" {
		t.Errorf("Received body differs: (-got +want)\n%s", diff)
	}

	app.Shutdown()
	wg.Wait()
}

func Test_server_samePort(t *testing.T) {
	var port int64 = 15261
	addr := ":" + strconv.FormatInt(port, 10)
	app, err := grapiserver.New().
		UseDefaultLogger().
		SetGrpcAddr("tcp", addr).
		SetGatewayAddr("tcp", addr).
		AddServers(
			server.NewLibraryServiceServer(),
		).
		Build()

	if err != nil {
		t.Fatalf("failed to build grapserver.Engine: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.Serve(); err != nil {
			t.Errorf("Engine.Serve returned an error: %v", err)
		}
	}()

	time.Sleep(2)

	t.Run("http", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/books", port))

		if err != nil {
			t.Fatalf("failed to fetch book resources: %v", err)
		}
		defer resp.Body.Close()

		if got, want := resp.StatusCode, 200; got != want {
			t.Fatalf("Response status is %d, want %d", got, want)
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read response body: %v", err)
		}

		got := map[string]interface{}{}
		err = json.Unmarshal(data, &got)
		if err != nil {
			t.Fatalf("failed to parse response body: %v", err)
		}
		want := map[string]interface{}{
			"books": []interface{}{
				map[string]interface{}{"book_id": "The Go Programming Language"},
				map[string]interface{}{"book_id": "Programming Ruby"},
			},
		}

		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("Received body differs: (-got +want)\n%s", diff)
		}
	})

	t.Run("gRPC", func(t *testing.T) {
		conn, err := grpc.Dial(addr, grpc.WithInsecure())
		if err != nil {
			t.Fatalf("failed to connect with gRPC server: %v", err)
		}
		defer conn.Close()

		cli := api_pb.NewLibraryServiceClient(conn)
		resp, err := cli.ListBooks(context.Background(), &api_pb.ListBooksRequest{})

		if err != nil {
			t.Fatalf("failed to fetch book resources: %v", err)
		}

		want := &api_pb.ListBooksResponse{
			Books: []*api_pb.Book{
				{BookId: "The Go Programming Language"},
				{BookId: "Programming Ruby"},
			},
		}
		if diff := cmp.Diff(resp, want); diff != "" {
			t.Errorf("Received body differs: (-got +want)\n%s", diff)
		}
	})

	app.Shutdown()
	wg.Wait()
}

func Test_server_differentPort(t *testing.T) {
	var (
		grpcPort int64 = 15261
		httpPort int64 = 15262
	)

	grpcAddr := ":" + strconv.FormatInt(grpcPort, 10)
	httpAddr := ":" + strconv.FormatInt(httpPort, 10)

	app, err := grapiserver.New().
		SetGrpcAddr("tcp", grpcAddr).
		SetGatewayAddr("tcp", httpAddr).
		AddServers(
			server.NewLibraryServiceServer(),
		).
		Build()

	if err != nil {
		t.Fatalf("failed to build grapserver.Engine: %v", err)
	}

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := app.Serve(); err != nil {
			t.Errorf("Engine.Serve returned an error: %v", err)
		}
	}()

	time.Sleep(2)

	t.Run("http", func(t *testing.T) {
		resp, err := http.Get(fmt.Sprintf("http://localhost:%d/books", httpPort))

		if err != nil {
			t.Fatalf("failed to fetch book resources: %v", err)
		}
		defer resp.Body.Close()

		if got, want := resp.StatusCode, 200; got != want {
			t.Fatalf("Response status is %d, want %d", got, want)
		}

		data, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			t.Fatalf("failed to read response body: %v", err)
		}

		got := map[string]interface{}{}
		err = json.Unmarshal(data, &got)
		if err != nil {
			t.Fatalf("failed to parse response body: %v", err)
		}
		want := map[string]interface{}{
			"books": []interface{}{
				map[string]interface{}{"book_id": "The Go Programming Language"},
				map[string]interface{}{"book_id": "Programming Ruby"},
			},
		}

		if diff := cmp.Diff(got, want); diff != "" {
			t.Errorf("Received body differs: (-got +want)\n%s", diff)
		}
	})

	t.Run("gRPC", func(t *testing.T) {
		conn, err := grpc.Dial(grpcAddr, grpc.WithInsecure())
		if err != nil {
			t.Fatalf("failed to connect with gRPC server: %v", err)
		}
		defer conn.Close()

		cli := api_pb.NewLibraryServiceClient(conn)
		resp, err := cli.ListBooks(context.Background(), &api_pb.ListBooksRequest{})

		if err != nil {
			t.Fatalf("failed to fetch book resources: %v", err)
		}

		want := &api_pb.ListBooksResponse{
			Books: []*api_pb.Book{
				{BookId: "The Go Programming Language"},
				{BookId: "Programming Ruby"},
			},
		}
		if diff := cmp.Diff(resp, want); diff != "" {
			t.Errorf("Received body differs: (-got +want)\n%s", diff)
		}
	})

	app.Shutdown()
	wg.Wait()
}
