package grapiserver

import (
	"context"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/izumin5210/grapi/pkg/grapiserver/internal"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

// NewGatewayServer creates GrpcServer instance.
func NewGatewayServer(c *Config) internal.Server {
	return &GatewayServer{
		Config: c,
	}
}

// GatewayServer wraps gRPC gateway server setup process.
type GatewayServer struct {
	server *http.Server
	*Config
}

// Serve implements Server.Shutdown
func (s *GatewayServer) Serve(l net.Listener, wg *sync.WaitGroup) {
	defer wg.Done()

	conn, err := s.createConn()
	if err != nil {
		grpclog.Errorf("failed to create connection with gRPC server: %v", err)
		return
	}
	defer conn.Close()

	s.server, err = s.createServer(conn)
	if err != nil {
		grpclog.Errorf("failed to create gRPC Gateway server: %v", err)
		return
	}

	err = s.server.Serve(l)
	grpclog.Infof("stopped taking more httr(s) requests: %v", err)
}

// Shutdown implements Server.Shutdown
func (s *GatewayServer) Shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err := s.server.Shutdown(ctx)
	grpclog.Info("All http(s) requets finished")
	if err != nil {
		grpclog.Errorf("failed to shutdown gRPC Gateway server: %v", err)
	}
}

func (s *GatewayServer) createConn() (conn *grpc.ClientConn, err error) {
	conn, err = grpc.Dial(s.GrpcInternalAddr.Addr, s.clientOptions()...)
	if err != nil {
		err = errors.Wrap(err, "failed to connect to gRPC server")
	}
	return
}

func (s *GatewayServer) createServer(conn *grpc.ClientConn) (*http.Server, error) {
	mux := runtime.NewServeMux(
		append(
			s.GatewayMuxOptions,
			runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler),
		)...,
	)
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	for _, svr := range s.Servers {
		err := svr.RegisterWithHandler(ctx, mux, conn)
		if err != nil {
			return nil, errors.Wrap(err, "failed to register handler")
		}
	}

	var handler http.Handler = mux

	for i := len(s.GatewayServerMiddlewares) - 1; i >= 0; i-- {
		handler = (s.GatewayServerMiddlewares[i])(handler)
	}

	return &http.Server{
		ReadTimeout:  8 * time.Second,
		WriteTimeout: 8 * time.Second,
		IdleTimeout:  2 * time.Minute,
		Handler:      handler,
	}, nil
}
