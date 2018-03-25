package foo

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
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

func (s *barBazServiceServerImpl) ListBarBazs(ctx context.Context, req *foo_pb.ListBarBazsRequest) (*foo_pb.ListBarBazsResponse, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *barBazServiceServerImpl) CreateBarBaz(ctx context.Context, req *foo_pb.CreateBarBazRequest) (*foo_pb.BarBaz, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

func (s *barBazServiceServerImpl) DeleteBarBaz(ctx context.Context, req *foo_pb.DeleteBarBazRequest) (*empty.Empty, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

