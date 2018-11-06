package server

import (
	"context"

	"github.com/izumin5210/grapi/pkg/grapiserver"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	api_pb "testapp/api"
)

// LibraryServiceServer is a composite interface of api_pb.LibraryServiceServer and grapiserver.Server.
type LibraryServiceServer interface {
	api_pb.LibraryServiceServer
	grapiserver.Server
}

// NewLibraryServiceServer creates a new LibraryServiceServer instance.
func NewLibraryServiceServer() LibraryServiceServer {
	return &libraryServiceServerImpl{}
}

type libraryServiceServerImpl struct {
}

