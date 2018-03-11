package grapiserver

import (
	"context"
	"net"
	"os"
	"path/filepath"
	"time"

	"github.com/grpc-ecosystem/go-grpc-middleware"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
)

var (
	// DefaultConfig is a default configuration.
	DefaultConfig = &Config{
		GrpcInternalAddr: &Address{
			Network: "unix",
			Addr:    "tmp/server.sock",
		},
		GatewayAddr: &Address{
			Network: "tcp",
			Addr:    ":3000",
		},
		MaxConcurrentStreams: 1000,
		Logger:               DefaultLogger,
	}
)

// RegisterGatewayHandlerFunc represents gRPC gateway's register handler functions.
type RegisterGatewayHandlerFunc func(context.Context, *runtime.ServeMux, *grpc.ClientConn) error

// RegisterGrpcServerImplFunc should register gRPC service server implementations to *grpc.Server.
type RegisterGrpcServerImplFunc func(s *grpc.Server)

// Address represents a network end point address.
type Address struct {
	Network string
	Addr    string
}

func (a *Address) createListener() (net.Listener, error) {
	if a.Network == "unix" {
		dir := filepath.Dir(a.Addr)
		f, err := os.Stat(dir)
		if err != nil {
			if err = os.MkdirAll(dir, 0755); err != nil {
				return nil, errors.Wrap(err, "failed to create the directory")
			}
		} else if !f.IsDir() {
			return nil, errors.Errorf("file %q already exists", dir)
		}
	}
	lis, err := net.Listen(a.Network, a.Addr)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to listen %s %s", a.Network, a.Addr)
	}
	return lis, nil
}

// Config contains configurations of gRPC and Gateway server.
type Config struct {
	GrpcInternalAddr                *Address
	GatewayAddr                     *Address
	RegisterGrpcServerImplFuncs     []RegisterGrpcServerImplFunc
	RegisterGatewayHandlerFuncs     []RegisterGatewayHandlerFunc
	GrpcServerUnaryInterceptors     []grpc.UnaryServerInterceptor
	GrpcServerStreamInterceptors    []grpc.StreamServerInterceptor
	GatewayServerUnaryInterceptors  []grpc.UnaryClientInterceptor
	GatewayServerStreamInterceptors []grpc.StreamClientInterceptor
	GatewayMuxOptions               []runtime.ServeMuxOption
	MaxConcurrentStreams            uint32
	Logger                          Logger
	GatewayServerMiddlewares        []HTTPServerMiddleware
}

func (c *Config) serverOptions() []grpc.ServerOption {
	return []grpc.ServerOption{
		grpc_middleware.WithUnaryServerChain(c.GrpcServerUnaryInterceptors...),
		grpc_middleware.WithStreamServerChain(c.GrpcServerStreamInterceptors...),
		grpc.MaxConcurrentStreams(c.MaxConcurrentStreams),
	}
}

func (c *Config) clientOptions() []grpc.DialOption {
	return []grpc.DialOption{
		grpc.WithInsecure(),
		grpc.WithDialer(func(a string, t time.Duration) (net.Conn, error) {
			return net.Dial(c.GrpcInternalAddr.Network, a)
		}),
		grpc.WithUnaryInterceptor(
			grpc_middleware.ChainUnaryClient(c.GatewayServerUnaryInterceptors...),
		),
		grpc.WithStreamInterceptor(
			grpc_middleware.ChainStreamClient(c.GatewayServerStreamInterceptors...),
		),
	}
}
