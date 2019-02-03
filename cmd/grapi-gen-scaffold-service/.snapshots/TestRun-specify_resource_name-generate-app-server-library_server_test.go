package server

import (
	"context"
	"testing"

	api_pb "testapp/api"
)

func Test_LibraryServiceServer_ListBooks(t *testing.T) {
	svr := NewLibraryServiceServer()

	ctx := context.Background()
	req := &api_pb.ListBooksRequest{}

	resp, err := svr.ListBooks(ctx, req)

	t.SkipNow()

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		t.Error("response should not nil")
	}
}

func Test_LibraryServiceServer_GetBook(t *testing.T) {
	svr := NewLibraryServiceServer()

	ctx := context.Background()
	req := &api_pb.GetBookRequest{}

	resp, err := svr.GetBook(ctx, req)

	t.SkipNow()

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		t.Error("response should not nil")
	}
}

func Test_LibraryServiceServer_CreateBook(t *testing.T) {
	svr := NewLibraryServiceServer()

	ctx := context.Background()
	req := &api_pb.CreateBookRequest{}

	resp, err := svr.CreateBook(ctx, req)

	t.SkipNow()

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		t.Error("response should not nil")
	}
}

func Test_LibraryServiceServer_UpdateBook(t *testing.T) {
	svr := NewLibraryServiceServer()

	ctx := context.Background()
	req := &api_pb.UpdateBookRequest{}

	resp, err := svr.UpdateBook(ctx, req)

	t.SkipNow()

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		t.Error("response should not nil")
	}
}

func Test_LibraryServiceServer_DeleteBook(t *testing.T) {
	svr := NewLibraryServiceServer()

	ctx := context.Background()
	req := &api_pb.DeleteBookRequest{}

	resp, err := svr.DeleteBook(ctx, req)

	t.SkipNow()

	if err != nil {
		t.Errorf("returned an error %v", err)
	}

	if resp == nil {
		t.Error("response should not nil")
	}
}

