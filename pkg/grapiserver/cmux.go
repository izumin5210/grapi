package grapiserver

import (
	"net"

	"github.com/soheilhy/cmux"
	"google.golang.org/grpc/grpclog"
)

type cmuxServer struct {
	mux cmux.CMux
	lis net.Listener
}

func newCmuxServer(lis net.Listener) *cmuxServer {
	return &cmuxServer{
		mux: cmux.New(lis),
		lis: lis,
	}
}

// Serve implements Server.Serve
func (s *cmuxServer) Serve() {
	grpclog.Info("mux is starting %s", s.lis.Addr())

	err := s.mux.Serve()

	grpclog.Infof("mux is closed: %v", err)
}

func (s *cmuxServer) GRPCListener() net.Listener {
	return s.mux.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
}

func (s *cmuxServer) HTTPListener() net.Listener {
	return s.mux.Match(cmux.HTTP2(), cmux.HTTP1Fast())
}
