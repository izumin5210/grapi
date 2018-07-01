package server

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"

	out_pb "testapp/api/out"
)

// RegisterWithServer implements grapiserver.Server.RegisterWithServer.
func (s *quuxServiceServerImpl) RegisterWithServer(grpcSvr *grpc.Server) {
	out_pb.RegisterQuuxServiceServer(grpcSvr, s)
}

// RegisterWithHandler implements grapiserver.Server.RegisterWithHandler.
func (s *quuxServiceServerImpl) RegisterWithHandler(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
	return out_pb.RegisterQuuxServiceHandler(ctx, mux, conn)
}

