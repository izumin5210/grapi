package foo

import (
	"context"

	"github.com/izumin5210/grapi/pkg/grapiserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	foo_pb "testapp/api/foo"
)

// BarBazServiceServer is a composite interface of foo_pb.BarBazServiceServer and grapiserver.Server.
type BarBazServiceServer interface {
	foo_pb.BarBazServiceServer
	grapiserver.Server
}

// NewBarBazServiceServer creates a new BarBazServiceServer instance.
func NewBarBazServiceServer() BarBazServiceServer {
	return &barBazServiceServerImpl{}
}

type barBazServiceServerImpl struct {
}

