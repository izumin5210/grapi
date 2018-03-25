package grapiserver

import (
	"os"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

// Builder creates an engine.
type Builder interface {
	SetGrpcInternalAddr(network, addr string) Builder
	SetGatewayAddr(network, addr string) Builder
	AddServers(svrs ...Server) Builder
	AddGrpcServerUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) Builder
	AddGrpcServerStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) Builder
	AddGatewayServerUnaryInterceptors(interceptors ...grpc.UnaryClientInterceptor) Builder
	AddGatewayServerStreamInterceptors(interceptors ...grpc.StreamClientInterceptor) Builder
	AddGrpcServerOptions(opts ...grpc.ServerOption) Builder
	AddGatewayDialOptions(opts ...grpc.DialOption) Builder
	AddGatewayMuxOptions(opts ...runtime.ServeMuxOption) Builder
	AddGatewayServerMiddleware(middlewares ...HTTPServerMiddleware) Builder
	AddPassedHeader(decider PassedHeaderDeciderFunc) Builder
	UseDefaultLogger() Builder
	Validate() error
	Build() (*Engine, error)
	Serve() error
}

// New creates a server builder object.
func New() Builder {
	return &builder{
		c: DefaultConfig,
	}
}

type builder struct {
	c *Config
}

func (b *builder) SetGrpcInternalAddr(network, addr string) Builder {
	b.c.GrpcInternalAddr = &Address{
		Network: network,
		Addr:    addr,
	}
	return b
}

func (b *builder) SetGatewayAddr(network, addr string) Builder {
	b.c.GatewayAddr = &Address{
		Network: network,
		Addr:    addr,
	}
	return b
}

func (b *builder) AddServers(svrs ...Server) Builder {
	b.c.Servers = append(b.c.Servers, svrs...)
	return b
}

func (b *builder) AddGrpcServerUnaryInterceptors(interceptors ...grpc.UnaryServerInterceptor) Builder {
	b.c.GrpcServerUnaryInterceptors = append(b.c.GrpcServerUnaryInterceptors, interceptors...)
	return b
}

func (b *builder) AddGrpcServerStreamInterceptors(interceptors ...grpc.StreamServerInterceptor) Builder {
	b.c.GrpcServerStreamInterceptors = append(b.c.GrpcServerStreamInterceptors, interceptors...)
	return b
}

func (b *builder) AddGatewayServerUnaryInterceptors(interceptors ...grpc.UnaryClientInterceptor) Builder {
	b.c.GatewayServerUnaryInterceptors = append(b.c.GatewayServerUnaryInterceptors, interceptors...)
	return b
}

func (b *builder) AddGatewayServerStreamInterceptors(interceptors ...grpc.StreamClientInterceptor) Builder {
	b.c.GatewayServerStreamInterceptors = append(b.c.GatewayServerStreamInterceptors, interceptors...)
	return b
}

func (b *builder) AddGrpcServerOptions(opts ...grpc.ServerOption) Builder {
	b.c.GrpcServerOption = append(b.c.GrpcServerOption, opts...)
	return b
}

func (b *builder) AddGatewayDialOptions(opts ...grpc.DialOption) Builder {
	b.c.GatewayDialOption = append(b.c.GatewayDialOption, opts...)
	return b
}

func (b *builder) AddGatewayMuxOptions(opts ...runtime.ServeMuxOption) Builder {
	b.c.GatewayMuxOptions = append(b.c.GatewayMuxOptions, opts...)
	return b
}

func (b *builder) AddGatewayServerMiddleware(middlewares ...HTTPServerMiddleware) Builder {
	b.c.GatewayServerMiddlewares = append(b.c.GatewayServerMiddlewares, middlewares...)
	return b
}

func (b *builder) AddPassedHeader(decider PassedHeaderDeciderFunc) Builder {
	return b.AddGatewayServerMiddleware(createPassingHeaderMiddleware(decider))
}

func (b *builder) UseDefaultLogger() Builder {
	grpclog.SetLoggerV2(grpclog.NewLoggerV2(os.Stdout, os.Stderr, os.Stderr))
	return b
}

func (b *builder) Validate() error {
	// TODO: not yet implemented
	return nil
}

func (b *builder) Build() (*Engine, error) {
	if err := b.Validate(); err != nil {
		return nil, errors.Wrap(err, "configuration is invalid")
	}
	return &Engine{
		Config: b.c,
	}, nil
}

func (b *builder) Serve() error {
	e, err := b.Build()
	if err != nil {
		return errors.Wrap(err, "failed to build server engine")
	}
	return e.Serve()
}
