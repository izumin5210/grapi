package server

import (
	"context"

	"github.com/izumin5210/grapi/pkg/grapiserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	out_pb "testapp/api/out"
)

// NewQuuxServiceServer creates a new QuuxServiceServer instance.
func NewQuuxServiceServer() interface {
	out_pb.QuuxServiceServer
	grapiserver.Server
} {
	return &quuxServiceServerImpl{}
}

type quuxServiceServerImpl struct {
}

