package server

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	api_pb "testapp/api_pb"
)

// RegisterWithServer implements grapiserver.Server.RegisterWithServer.
func (s *fooServiceServerImpl) RegisterWithServer(grpcSvr *grpc.Server) {
	api_pb.RegisterFooServiceServer(grpcSvr, s)
}

// RegisterWithHandler implements grapiserver.Server.RegisterWithHandler.
func (s *fooServiceServerImpl) RegisterWithHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return api_pb.RegisterFooServiceHandler(ctx, mux, conn)
}

