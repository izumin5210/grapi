package foo

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	foo_pb "testapp/api/foo"
)

// RegisterWithServer implements grapiserver.Server.RegisterWithServer.
func (s *barBazServiceServerImpl) RegisterWithServer(grpcSvr *grpc.Server) {
	foo_pb.RegisterBarBazServiceServer(grpcSvr, s)
}

// RegisterWithHandler implements grapiserver.Server.RegisterWithHandler.
func (s *barBazServiceServerImpl) RegisterWithHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return foo_pb.RegisterBarBazServiceHandler(ctx, mux, conn)
}

