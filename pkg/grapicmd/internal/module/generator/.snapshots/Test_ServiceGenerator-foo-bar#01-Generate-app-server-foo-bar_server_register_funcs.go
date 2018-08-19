package foo

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	foo_pb "testapp/api/foo"
)

// RegisterWithServer implements grapiserver.Server.RegisterWithServer.
func (s *barServiceServerImpl) RegisterWithServer(grpcSvr *grpc.Server) {
	foo_pb.RegisterBarServiceServer(grpcSvr, s)
}

// RegisterWithHandler implements grapiserver.Server.RegisterWithHandler.
func (s *barServiceServerImpl) RegisterWithHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return foo_pb.RegisterBarServiceHandler(ctx, mux, conn)
}

