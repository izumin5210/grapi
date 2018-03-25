package grapiserver

import (
	"net"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"

	"github.com/izumin5210/grapi/pkg/grapiserver/internal"
	"github.com/pkg/errors"
	"github.com/soheilhy/cmux"
	"google.golang.org/grpc/grpclog"
)

// Engine is the framework instance.
type Engine struct {
	*Config
	sdCh chan os.Signal
}

// Serve starts gRPC and Gateway servers.
func (e *Engine) Serve() error {
	var (
		grpcServer, gatewayServer, muxServer internal.Server
		grpcLis, gatewayLis, internalLis     net.Listener
		err                                  error
	)

	if e.GrpcAddr != nil && e.GatewayAddr != nil && reflect.DeepEqual(e.GrpcAddr, e.GatewayAddr) {
		lis, err := e.GrpcAddr.createListener()
		if err != nil {
			return errors.Wrap(err, "failed to listen network for servers")
		}
		mux := cmux.New(lis)
		muxServer = NewMuxServer(mux, lis)
		grpcLis = mux.Match(cmux.HTTP2HeaderField("content-type", "application/grpc"))
		gatewayLis = mux.Match(cmux.HTTP2(), cmux.HTTP1Fast())
	}

	// Setup servers
	grpcServer = NewGrpcServer(e.Config)

	// Setup listeners
	if grpcLis == nil && e.GrpcAddr != nil {
		grpcLis, err = e.GrpcAddr.createListener()
		if err != nil {
			return errors.Wrap(err, "failed to listen network for gRPC server")
		}
		defer grpcLis.Close()
	}

	if e.GatewayAddr != nil {
		gatewayServer = NewGatewayServer(e.Config)
		internalLis, err = e.GrpcInternalAddr.createListener()
		if err != nil {
			return errors.Wrap(err, "failed to listen network for gRPC server internal")
		}
		defer internalLis.Close()
	}

	if gatewayLis == nil && e.GatewayAddr != nil {
		gatewayLis, err = e.GatewayAddr.createListener()
		if err != nil {
			return errors.Wrap(err, "failed to listen network for gateway server")
		}
		defer gatewayLis.Close()
	}

	// Start servers
	var wg sync.WaitGroup

	if internalLis != nil {
		wg.Add(1)
		grpclog.Infof("gRPC server is starting %s://%s", e.GrpcInternalAddr.Network, e.GrpcInternalAddr.Addr)
		go grpcServer.Serve(internalLis, &wg)
	}
	if grpcLis != nil {
		wg.Add(1)
		grpclog.Infof("gRPC server is starting %s://%s", e.GrpcAddr.Network, e.GrpcAddr.Addr)
		go grpcServer.Serve(grpcLis, &wg)
	}
	if gatewayLis != nil {
		wg.Add(1)
		grpclog.Infof("gRPC Gateway server is starting: %s://%s", e.GatewayAddr.Network, e.GatewayAddr.Addr)
		go gatewayServer.Serve(gatewayLis, &wg)
	}
	if muxServer != nil {
		wg.Add(1)
		go muxServer.Serve(nil, &wg)
	}

	wg.Add(1)
	go e.watchShutdownSignal(&wg, gatewayServer, grpcServer, muxServer)

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
		if svr != nil {
			svr.Shutdown()
		}
	}
}
