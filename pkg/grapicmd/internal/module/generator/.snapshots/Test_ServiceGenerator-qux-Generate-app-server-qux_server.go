package server

import (
	"context"

	"github.com/izumin5210/grapi/pkg/grapiserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_pb "testapp/api"
)

// QuxServiceServer is a composite interface of api_pb.QuxServiceServer and grapiserver.Server.
type QuxServiceServer interface {
	api_pb.QuxServiceServer
	grapiserver.Server
}

// NewQuxServiceServer creates a new QuxServiceServer instance.
func NewQuxServiceServer() QuxServiceServer {
	return &quxServiceServerImpl{}
}

type quxServiceServerImpl struct {
}

