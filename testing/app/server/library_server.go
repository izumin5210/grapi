package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/izumin5210/grapi/pkg/grapiserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_pb "github.com/izumin5210/grapi/testing/api"
)

// NewLibraryServiceServer creates a new LibraryServiceServer instance.
func NewLibraryServiceServer() interface {
	api_pb.LibraryServiceServer
	grapiserver.Server
} {
	return &libraryServiceServerImpl{}
}

type libraryServiceServerImpl struct {
}

// RegisterWithServer implements grapiserver.Server.RegisterWithServer.
func (s *libraryServiceServerImpl) RegisterWithServer(grpcSvr *grpc.Server) {
	api_pb.RegisterLibraryServiceServer(grpcSvr, s)
}

// RegisterWithHandler implements grapiserver.Server.RegisterWithHandler.
func (s *libraryServiceServerImpl) RegisterWithHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return api_pb.RegisterLibraryServiceHandler(ctx, mux, conn)
}

func (s *libraryServiceServerImpl) ListBooks(ctx context.Context, req *api_pb.ListBooksRequest) (*api_pb.ListBooksResponse, error) {
	return &api_pb.ListBooksResponse{
		Books: []*api_pb.Book{
			{BookId: "The Go Programming Language"},
			{BookId: "Programming Ruby"},
		},
	}, nil
}

func (s *libraryServiceServerImpl) GetBook(ctx context.Context, req *api_pb.GetBookRequest) (*api_pb.Book, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *libraryServiceServerImpl) CreateBook(ctx context.Context, req *api_pb.CreateBookRequest) (*api_pb.Book, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *libraryServiceServerImpl) UpdateBook(ctx context.Context, req *api_pb.UpdateBookRequest) (*api_pb.Book, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *libraryServiceServerImpl) DeleteBook(ctx context.Context, req *api_pb.DeleteBookRequest) (*empty.Empty, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}
