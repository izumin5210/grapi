package grapiserver

import (
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"

	"github.com/izumin5210/grapi/pkg/grapiserver/internal"
)

// GrpcServer wraps grpc.Server setup process.
type GrpcServer struct {
	server *grpc.Server
	*Config
}

// NewGrpcServer creates GrpcServer instance.
func NewGrpcServer(c *Config) internal.Server {
	s := grpc.NewServer(c.serverOptions()...)
	reflection.Register(s)
	for _, svr := range c.Servers {
		svr.RegisterWithServer(s)
	}
	return &GrpcServer{
		server: s,
		Config: c,
	}
}

// Serve implements Server.Shutdown
func (s *GrpcServer) Serve(l net.Listener, wg *sync.WaitGroup) {
	defer wg.Done()

	grpclog.Infof("gRPC server is starting %s://%s", s.GrpcInternalAddr.Network, s.GrpcInternalAddr.Addr)
	err := s.server.Serve(l)
	grpclog.Infof("gRPC server stopred: %v", err)
}

// Shutdown implements Server.Shutdown
func (s *GrpcServer) Shutdown() {
	s.server.GracefulStop()
}
