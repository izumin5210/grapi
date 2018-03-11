package grapiserver

import (
	"fmt"
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// GrpcServer wraps grpc.Server setup process.
type GrpcServer struct {
	server *grpc.Server
	*Config
}

// NewGrpcServer creates GrpcServer instance.
func NewGrpcServer(c *Config) Server {
	s := grpc.NewServer(c.serverOptions()...)
	reflection.Register(s)
	c.Logger.Info(fmt.Sprintf("Register %d server impls to gRPC server", len(c.RegisterGrpcServerImplFuncs)), LogFields{})
	for _, register := range c.RegisterGrpcServerImplFuncs {
		register(s)
	}
	return &GrpcServer{
		server: s,
		Config: c,
	}
}

// Serve implements Server.Shutdown
func (s *GrpcServer) Serve(l net.Listener, wg *sync.WaitGroup) {
	defer wg.Done()

	s.Logger.Info("gRPC server is starting", LogFields{"network": s.GrpcInternalAddr.Network, "addr": s.GrpcInternalAddr.Addr})
	err := s.server.Serve(l)
	s.Logger.Info("gRPC server stopred", LogFields{"error": err})
}

// Shutdown implements Server.Shutdown
func (s *GrpcServer) Shutdown() {
	s.server.GracefulStop()
}
