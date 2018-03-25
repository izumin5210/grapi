package server

import (
	"context"
	"testing"

	"github.com/golang/protobuf/ptypes/empty"

	api_pb "testapp/api"
)

func Test_LibraryServiceServer_ListLibraries(t *testing.T) {
	svr := NewLibraryServiceServer()

	ctx := context.Background()
	req := &api_pb.ListLibrariesResponse{}

	resp, err := svr.ListLibraries(ctx, req)

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		at.Error("response should not nil")
	}
}

func Test_LibraryServiceServer_GetLibrary(t *testing.T) {
	svr := NewLibraryServiceServer()

	ctx := context.Background()
	req := &api_pb.Library{}

	resp, err := svr.GetLibrary(ctx, req)

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		at.Error("response should not nil")
	}
}

func Test_LibraryServiceServer_CreateLibrary(t *testing.T) {
	svr := NewLibraryServiceServer()

	ctx := context.Background()
	req := &api_pb.Library{}

	resp, err := svr.CreateLibrary(ctx, req)

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		at.Error("response should not nil")
	}
}

func Test_LibraryServiceServer_UpdateLibrary(t *testing.T) {
	svr := NewLibraryServiceServer()

	ctx := context.Background()
	req := &api_pb.Library{}

	resp, err := svr.UpdateLibrary(ctx, req)

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		at.Error("response should not nil")
	}
}

func Test_LibraryServiceServer_DeleteLibrary(t *testing.T) {
	svr := NewLibraryServiceServer()

	ctx := context.Background()
	req := &empty.Empty{}

	resp, err := svr.DeleteLibrary(ctx, req)

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		at.Error("response should not nil")
	}
}

