package server

import (
	"context"

	"github.com/izumin5210/grapi/pkg/grapiserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_pb "testapp/api_pb"
)

// NewFooServiceServer creates a new FooServiceServer instance.
func NewFooServiceServer() interface {
	api_pb.FooServiceServer
	grapiserver.Server
} {
	return &fooServiceServerImpl{}
}

type fooServiceServerImpl struct {
}

