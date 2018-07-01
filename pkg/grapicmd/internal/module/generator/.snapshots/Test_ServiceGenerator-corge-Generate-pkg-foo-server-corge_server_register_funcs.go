package server

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	api_pb "testapp/api"
)

// RegisterWithServer implements grapiserver.Server.RegisterWithServer.
func (s *corgeServiceServerImpl) RegisterWithServer(grpcSvr *grpc.Server) {
	api_pb.RegisterCorgeServiceServer(grpcSvr, s)
}

// RegisterWithHandler implements grapiserver.Server.RegisterWithHandler.
func (s *corgeServiceServerImpl) RegisterWithHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return api_pb.RegisterCorgeServiceHandler(ctx, mux, conn)
}

