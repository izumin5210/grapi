package server

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_pb "testapp/api"
)

var (
	// RegisterFooServiceHandler is a function to register card service handler to gRPC Gateway's mux.
	RegisterFooServiceHandler = api_pb.RegisterFooServiceHandler
)

// RegisterFooServiceServerFactory creates a function to register card service server impl to grpc.Server.
func RegisterFooServiceServerFactory() func(s *grpc.Server) {
	return func(s *grpc.Server) {
		api_pb.RegisterFooServiceServer(s, New())
	}
}

// New creates a new FooServiceServer instance.
func New() api_pb.FooServiceServer {
	return &fooServiceServerImpl{}
}

type fooServiceServerImpl struct {
}

func (s *fooServiceServerImpl) GetFoo(ctx context.Context, req *api_pb.GetFooRequest) (*api_pb.GetFooResponse, error) {
	// TODO: Not yet implemented.
	return nil, status.Error(codes.Unimplemented, "TODO: You should implement it!")
}

