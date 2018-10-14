package foo

import (
	"context"

	"github.com/izumin5210/grapi/pkg/grapiserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	foo_pb "testapp/api/foo"
)

// BarServiceServer is a composite interface of foo_pb.BarServiceServer and grapiserver.Server.
type BarServiceServer interface {
	foo_pb.BarServiceServer
	grapiserver.Server
}

// NewBarServiceServer creates a new BarServiceServer instance.
func NewBarServiceServer() BarServiceServer {
	return &barServiceServerImpl{}
}

type barServiceServerImpl struct {
}

