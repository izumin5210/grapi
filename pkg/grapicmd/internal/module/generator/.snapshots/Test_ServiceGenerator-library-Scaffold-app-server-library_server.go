package server

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/izumin5210/grapi/pkg/grapiserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_pb "testapp/api"
)

// New creates a new LibraryServiceServer instance.
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

func (s *libraryServiceServerImpl) ListLibraries(ctx context.Context, req *api_pb.ListLibrariesRequest) (*api_pb.ListLibrariesResponse, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *libraryServiceServerImpl) GetLibrary(ctx context.Context, req *api_pb.GetLibraryRequest) (*api_pb.Library, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *libraryServiceServerImpl) CreateLibrary(ctx context.Context, req *api_pb.CreateLibraryRequest) (*api_pb.Library, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *libraryServiceServerImpl) UpdateLibrary(ctx context.Context, req *api_pb.UpdateLibraryRequest) (*api_pb.Library, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *libraryServiceServerImpl) DeleteLibrary(ctx context.Context, req *api_pb.DeleteLibraryRequest) (*empty.Empty, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

