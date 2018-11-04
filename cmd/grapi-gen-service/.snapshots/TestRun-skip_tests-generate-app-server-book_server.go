package server

import (
	"context"

	"github.com/izumin5210/grapi/pkg/grapiserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_pb "testapp/api"
)

// BookServiceServer is a composite interface of api_pb.BookServiceServer and grapiserver.Server.
type BookServiceServer interface {
	api_pb.BookServiceServer
	grapiserver.Server
}

// NewBookServiceServer creates a new BookServiceServer instance.
func NewBookServiceServer() BookServiceServer {
	return &bookServiceServerImpl{}
}

type bookServiceServerImpl struct {
}

