package foo

import (
	"context"

	"github.com/izumin5210/grapi/pkg/grapiserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	foo_pb "testapp/api/foo_pb"
)

// NewBarBazServiceServer creates a new BarBazServiceServer instance.
func NewBarBazServiceServer() interface {
	foo_pb.BarBazServiceServer
	grapiserver.Server
} {
	return &barBazServiceServerImpl{}
}

type barBazServiceServerImpl struct {
}

