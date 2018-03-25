package testing

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strconv"
	"sync"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/izumin5210/grapi/pkg/grapiserver"
	"github.com/izumin5210/grapi/testing/app/server"
)

func Test_server(t *testing.T) {
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
