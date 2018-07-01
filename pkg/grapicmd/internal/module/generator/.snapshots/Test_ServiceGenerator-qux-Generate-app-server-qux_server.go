package server

import (
	"context"

	"github.com/izumin5210/grapi/pkg/grapiserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_pb "testapp/api"
)

// NewQuxServiceServer creates a new QuxServiceServer instance.
func NewQuxServiceServer() interface {
	api_pb.QuxServiceServer
	grapiserver.Server
} {
	return &quxServiceServerImpl{}
}

type quxServiceServerImpl struct {
}

