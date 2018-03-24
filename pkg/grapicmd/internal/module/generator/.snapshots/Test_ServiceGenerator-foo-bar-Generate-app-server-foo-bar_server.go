package foo

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/izumin5210/grapi/pkg/grapiserver"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	foo_pb "testapp/api/foo"
)

// New creates a new BarServiceServer instance.
func NewBarServiceServer() interface {
	foo_pb.BarServiceServer
	grapiserver.Server
} {
	return &barServiceServerImpl{}
}

type barServiceServerImpl struct {
}

// RegisterWithServer implements grapiserver.Server.RegisterWithServer.
func (s *barServiceServerImpl) RegisterWithServer(grpcSvr *grpc.Server) {
	foo_pb.RegisterBarServiceServer(grpcSvr, s)
}

// RegisterWithHandler implements grapiserver.Server.RegisterWithHandler.
func (s *barServiceServerImpl) RegisterWithHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return foo_pb.RegisterBarServiceHandler(ctx, mux, conn)
}

