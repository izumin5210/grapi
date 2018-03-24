package foo

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	foo_pb "testapp/api/foo"
)

var (
	// RegisterBarServiceHandler is a function to register card service handler to gRPC Gateway's mux.
	RegisterBarServiceHandler = foo_pb.RegisterBarServiceHandler
)

// RegisterBarServiceServerFactory creates a function to register card service server impl to grpc.Server.
func RegisterBarServiceServerFactory() func(s *grpc.Server) {
	return func(s *grpc.Server) {
		foo_pb.RegisterBarServiceServer(s, New())
	}
}

// New creates a new BarServiceServer instance.
func New() foo_pb.BarServiceServer {
	return &barServiceServerImpl{}
}

type barServiceServerImpl struct {
}

