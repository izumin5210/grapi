package foo

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	foo_pb "testapp/api/foo"
)

var (
	// RegisterBarBazServiceHandler is a function to register card service handler to gRPC Gateway's mux.
	RegisterBarBazServiceHandler = foo_pb.RegisterBarBazServiceHandler
)

// RegisterBarBazServiceServerFactory creates a function to register card service server impl to grpc.Server.
func RegisterBarBazServiceServerFactory() func(s *grpc.Server) {
	return func(s *grpc.Server) {
		foo_pb.RegisterBarBazServiceServer(s, New())
	}
}

// New creates a new BarBazServiceServer instance.
func New() foo_pb.BarBazServiceServer {
	return &barBazServiceServerImpl{}
}

type barBazServiceServerImpl struct {
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

