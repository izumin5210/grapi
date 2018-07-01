package server

import (
	"context"

	"github.com/izumin5210/grapi/pkg/grapiserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_pb "testapp/api"
)

// NewCorgeServiceServer creates a new CorgeServiceServer instance.
func NewCorgeServiceServer() interface {
	api_pb.CorgeServiceServer
	grapiserver.Server
} {
	return &corgeServiceServerImpl{}
}

type corgeServiceServerImpl struct {
}

