package grapiserver

import (
	"net"
	"sync"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
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
	grpclog.Infof("register %d server impls to gRPC server", len(c.RegisterGrpcServerImplFuncs))
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

	grpclog.Infof("gRPC server is starting %s://%s", s.GrpcInternalAddr.Network, s.GrpcInternalAddr.Addr)
	err := s.server.Serve(l)
	grpclog.Info("gRPC server stopred: %v", err)
}

// Shutdown implements Server.Shutdown
func (s *GrpcServer) Shutdown() {
	s.server.GracefulStop()
}
