package server

import (
	"context"

	"github.com/izumin5210/grapi/pkg/grapiserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	out_pb "testapp/api/out"
)

// QuuxServiceServer is a composite interface of out_pb.QuuxServiceServer and grapiserver.Server.
type QuuxServiceServer interface {
	out_pb.QuuxServiceServer
	grapiserver.Server
}

// NewQuuxServiceServer creates a new QuuxServiceServer instance.
func NewQuuxServiceServer() QuuxServiceServer {
	return &quuxServiceServerImpl{}
}

type quuxServiceServerImpl struct {
}

