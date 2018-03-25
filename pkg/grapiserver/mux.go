package grapiserver

import (
	"net"
	"sync"

	"github.com/izumin5210/grapi/pkg/grapiserver/internal"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc/grpclog"
)

// MuxServer wraps a conneciton multiplexer and a listener.
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
func (s *MuxServer) Serve(lis net.Listener, wg *sync.WaitGroup) {
	defer wg.Done()
	grpclog.Info("mux is starting")
	err := s.mux.Serve()
	grpclog.Infof("mux is closed: %v", err)
}

// Shutdown implements Server.Shutdown
func (s *MuxServer) Shutdown() {
	err := s.lis.Close()
	if err != nil {
		grpclog.Errorf("failed to close cmux's listener: %v", err)
	}
}
