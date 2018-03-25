package grapiserver

import (
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/izumin5210/grapi/pkg/grapiserver/internal"
	"github.com/pkg/errors"
	"google.golang.org/grpc/grpclog"
)

// Engine is the framework instance.
type Engine struct {
	*Config
	sdCh chan os.Signal
}

// Serve starts gRPC and Gateway servers.
func (e *Engine) Serve() error {
	// Setup gRPC server
	grpcServer := NewGrpcServer(e.Config)
	grpcLis, err := e.GrpcInternalAddr.createListener()
	if err != nil {
		return errors.Wrap(err, "failed to listen network for gRPC server")
	}
	defer grpcLis.Close()

	// Setup gRPC gateway server
	gatewayServer := NewGatewayServer(e.Config)
	gatewayLis, err := e.GatewayAddr.createListener()
	if err != nil {
		return errors.Wrap(err, "failed to listen network for gateway server")
	}
	defer gatewayLis.Close()

	// Start servers
	var wg sync.WaitGroup
	wg.Add(3)

	go grpcServer.Serve(grpcLis, &wg)
	go gatewayServer.Serve(gatewayLis, &wg)
	go e.watchShutdownSignal(&wg, gatewayServer, grpcServer)

	wg.Wait()

	return nil
}

// Shutdown closes servers.
func (e *Engine) Shutdown() {
	e.sdCh <- os.Interrupt
}

func (e *Engine) watchShutdownSignal(wg *sync.WaitGroup, servers ...internal.Server) {
	defer wg.Done()
	e.sdCh = make(chan os.Signal, 1)
	defer close(e.sdCh)
	defer signal.Reset()
	signal.Notify(e.sdCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	sig := <-e.sdCh
	grpclog.Infof("terminating now...: %v", sig)
	for _, svr := range servers {
		svr.Shutdown()
	}
}
