package server

import (
	"context"

	"github.com/izumin5210/grapi/pkg/grapiserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_pb "testapp/api"
)

// FooServiceServer is a composite interface of api_pb.FooServiceServer and grapiserver.Server.
type FooServiceServer interface {
	api_pb.FooServiceServer
	grapiserver.Server
}

// NewFooServiceServer creates a new FooServiceServer instance.
func NewFooServiceServer() FooServiceServer {
	return &fooServiceServerImpl{}
}

type fooServiceServerImpl struct {
}

