package server

import (
	"context"

	"github.com/izumin5210/grapi/pkg/grapiserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_pb "testapp/api"
)

// CorgeServiceServer is a composite interface of api_pb.CorgeServiceServer and grapiserver.Server.
type CorgeServiceServer interface {
	api_pb.CorgeServiceServer
	grapiserver.Server
}

// NewCorgeServiceServer creates a new CorgeServiceServer instance.
func NewCorgeServiceServer() CorgeServiceServer {
	return &corgeServiceServerImpl{}
}

type corgeServiceServerImpl struct {
}

