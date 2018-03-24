package server

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/izumin5210/grapi/pkg/grapiserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_pb "testapp/api"
)

// New creates a new FooServiceServer instance.
func NewFooServiceServer() interface {
	api_pb.FooServiceServer
	grapiserver.Server
} {
	return &fooServiceServerImpl{}
}

type fooServiceServerImpl struct {
}

// RegisterWithServer implements grapiserver.Server.RegisterWithServer.
func (s *fooServiceServerImpl) RegisterWithServer(grpcSvr *grpc.Server) {
	api_pb.RegisterFooServiceServer(grpcSvr, s)
}

// RegisterWithHandler implements grapiserver.Server.RegisterWithHandler.
func (s *fooServiceServerImpl) RegisterWithHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return api_pb.RegisterFooServiceHandler(ctx, mux, conn)
}

