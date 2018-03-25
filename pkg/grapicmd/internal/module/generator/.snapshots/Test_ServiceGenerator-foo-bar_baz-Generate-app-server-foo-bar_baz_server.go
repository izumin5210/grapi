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

// NewBarBazServiceServer creates a new BarBazServiceServer instance.
func NewBarBazServiceServer() interface {
	foo_pb.BarBazServiceServer
	grapiserver.Server
} {
	return &barBazServiceServerImpl{}
}

type barBazServiceServerImpl struct {
}

// RegisterWithServer implements grapiserver.Server.RegisterWithServer.
func (s *barBazServiceServerImpl) RegisterWithServer(grpcSvr *grpc.Server) {
	foo_pb.RegisterBarBazServiceServer(grpcSvr, s)
}

// RegisterWithHandler implements grapiserver.Server.RegisterWithHandler.
func (s *barBazServiceServerImpl) RegisterWithHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return foo_pb.RegisterBarBazServiceHandler(ctx, mux, conn)
}

