package foo

import (
	"context"

	"github.com/izumin5210/grapi/pkg/grapiserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	foo_pb "testapp/api/foo"
)

// NewBarServiceServer creates a new BarServiceServer instance.
func NewBarServiceServer() interface {
	foo_pb.BarServiceServer
	grapiserver.Server
} {
	return &barServiceServerImpl{}
}

type barServiceServerImpl struct {
}

