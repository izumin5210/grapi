package foo

import (
	"context"

	"github.com/golang/protobuf/ptypes/empty"
	"github.com/izumin5210/grapi/pkg/grapiserver"
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

