package grapiserver

import (
	"context"
	"net"

	"github.com/izumin5210/grapi/pkg/grapiserver/internal"
	"github.com/pkg/errors"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc/grpclog"
)

// MuxServer wraps a connection multiplexer and a listener.
type MuxServer struct {
	mux cmux.CMux
	lis net.Listener
}

// NewMuxServer creates MuxServer instance.
func NewMuxServer(mux cmux.CMux, lis net.Listener) internal.Server {
	return &MuxServer{
		mux: mux,
		lis: lis,
	}
}

// Serve implements Server.Serve
func (s *MuxServer) Serve(ctx context.Context, _ net.Listener) error {
	grpclog.Info("mux is starting %s", s.lis.Addr())

	err := internal.StartServer(ctx, s.mux.Serve, s.Shutdown)

	grpclog.Infof("mux is closed: %v", err)

	return errors.Wrap(err, "failed to serve cmux server")
}

// Shutdown implements Server.Shutdown
func (s *MuxServer) Shutdown() {
	err := s.lis.Close()
	if err != nil {
		grpclog.Errorf("failed to close cmux's listener: %v", err)
	}
}
