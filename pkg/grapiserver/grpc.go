package grapiserver

import (
	"context"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
	"google.golang.org/grpc/reflection"

	"github.com/izumin5210/grapi/pkg/grapiserver/internal"
	"github.com/pkg/errors"
)

// grpcServer wraps grpc.Server setup process.
type grpcServer struct {
	server *grpc.Server
	*Config
}

func newGRPCServer(c *Config) internal.Server {
	s := grpc.NewServer(c.serverOptions()...)
	reflection.Register(s)
	for _, svr := range c.Servers {
		svr.RegisterWithServer(s)
	}
	return &grpcServer{
		server: s,
		Config: c,
	}
}

// Serve implements Server.Server
func (s *grpcServer) Serve(ctx context.Context, l net.Listener) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		<-ctx.Done()
		s.server.GracefulStop()
	}()

	grpclog.Infof("gRPC server is starting %s", l.Addr())

	err := s.server.Serve(l)

	grpclog.Infof("gRPC server stopped: %v", err)

	return errors.Wrap(err, "failed to serve gRPC server")
}
