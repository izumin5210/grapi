package grapiserver

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/izumin5210/grapi/pkg/grapiserver/internal"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/grpclog"
)

func newGatewayServer(c *Config) internal.Server {
	return &gatewayServer{
		Config: c,
	}
}

// gatewayServer wraps gRPC gateway server setup process.
type gatewayServer struct {
	server *http.Server
	*Config
}

// Serve implements Server.Server
func (s *gatewayServer) Serve(ctx context.Context, l net.Listener) error {
	conn, err := s.createConn()
	if err != nil {
		return errors.Wrap(err, "failed to create connection with grpc-gateway server")
	}
	defer conn.Close()

	s.server, err = s.createServer(conn)
	if err != nil {
		return errors.Wrap(err, "failed to create gRPC Gateway server: %v")
	}

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	go func() {
		<-ctx.Done()
		s.shutdown()
	}()

	grpclog.Infof("grpc-gateway server is starting %s", l.Addr())

	err = s.server.Serve(l)

	grpclog.Infof("stopped taking more http(s) requests: %v", err)

	if err != http.ErrServerClosed {
		return errors.Wrap(err, "failed to serve grpc-gateway server")
	}

	return nil
}

func (s *gatewayServer) shutdown() {
	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()
	err := s.server.Shutdown(ctx)
	grpclog.Info("All http(s) requests finished")
	if err != nil {
		grpclog.Errorf("failed to shutdown grpc-gateway server: %v", err)
	}
}

func (s *gatewayServer) createConn() (conn *grpc.ClientConn, err error) {
	conn, err = grpc.Dial(s.GrpcInternalAddr.Addr, s.clientOptions()...)
	if err != nil {
		err = errors.Wrap(err, "failed to connect to gRPC server")
	}
	return
}

func (s *gatewayServer) createServer(conn *grpc.ClientConn) (*http.Server, error) {
	mux := runtime.NewServeMux(
		append(
			[]runtime.ServeMuxOption{runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler)},
			s.GatewayMuxOptions...,
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

	svr := &http.Server{
		Handler: handler,
	}
	if cfg := s.GatewayServerConfig; cfg != nil {
		cfg.applyTo(svr)
	}

	return svr, nil
}
