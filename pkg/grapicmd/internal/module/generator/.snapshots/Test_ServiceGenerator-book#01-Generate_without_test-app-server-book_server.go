package server

import (
	"context"

	"github.com/izumin5210/grapi/pkg/grapiserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_pb "testapp/api"
)

// NewBookServiceServer creates a new BookServiceServer instance.
func NewBookServiceServer() interface {
	api_pb.BookServiceServer
	grapiserver.Server
} {
	return &bookServiceServerImpl{}
}

type bookServiceServerImpl struct {
}

