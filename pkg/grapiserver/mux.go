package grapiserver

import (
	"net"

	"github.com/soheilhy/cmux"
	"google.golang.org/grpc/grpclog"
)

// MuxServer wraps a connection multiplexer and a listener.
type MuxServer struct {
	mux cmux.CMux
	lis net.Listener
}

// NewMuxServer creates MuxServer instance.
func NewMuxServer(lis net.Listener) *MuxServer {
	return &MuxServer{
		mux: cmux.New(lis),
		lis: lis,
	}
}

// Serve implements Server.Serve
func (s *MuxServer) Serve() {
	grpclog.Info("mux is starting %s", s.lis.Addr())

	err := s.mux.Serve()

	grpclog.Infof("mux is closed: %v", err)
}

func (s *MuxServer) GRPCListener() net.Listener {
	return s.mux.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
}

func (s *MuxServer) HTTPListener() net.Listener {
	return s.mux.Match(cmux.HTTP2(), cmux.HTTP1Fast())
}
