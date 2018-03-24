package foo

import (
	"context"

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

