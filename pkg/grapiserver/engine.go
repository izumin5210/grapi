package grapiserver

import (
	"context"
	"net"
	"os"
	"os/signal"
	"reflect"
	"sync"
	"syscall"

	"github.com/izumin5210/grapi/pkg/grapiserver/internal"
	"github.com/pkg/errors"
	"golang.org/x/sync/errgroup"
	"google.golang.org/grpc/grpclog"
)

// Engine is the framework instance.
type Engine struct {
	*Config
}

// New creates a server intstance.
func New(opts ...Option) *Engine {
	return &Engine{
		Config: createConfig(opts),
	}
}

// Serve starts gRPC and Gateway servers.
func (e *Engine) Serve() error {
	return errors.WithStack(e.ServeContext(context.Background()))
}

// ServeContext starts gRPC and Gateway servers.
func (e *Engine) ServeContext(ctx context.Context) error {
	var (
		grpcServer, gatewayServer        internal.Server
		grpcLis, gatewayLis, internalLis net.Listener
		cmuxServer                       *cmuxServer
		err                              error
	)

	if e.GrpcAddr != nil && e.GatewayAddr != nil && reflect.DeepEqual(e.GrpcAddr, e.GatewayAddr) {
		lis, err := e.GrpcAddr.createListener()
		if err != nil {
			return errors.Wrap(err, "failed to listen network for servers")
		}
		defer lis.Close()
		cmuxServer = newCmuxServer(lis)
		grpcLis = cmuxServer.GRPCListener()
		gatewayLis = cmuxServer.HTTPListener()
	}

	// Setup servers
	grpcServer = newGRPCServer(e.Config)

	// Setup listeners
	if grpcLis == nil && e.GrpcAddr != nil {
		grpcLis, err = e.GrpcAddr.createListener()
		if err != nil {
			return errors.Wrap(err, "failed to listen network for gRPC server")
		}
		defer grpcLis.Close()
	}

	if e.GatewayAddr != nil {
		gatewayServer = newGatewayServer(e.Config)
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
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()
	eg, ctx := errgroup.WithContext(ctx)

	if internalLis != nil {
		eg.Go(func() error { return grpcServer.Serve(ctx, internalLis) })
	}
	if grpcLis != nil {
		eg.Go(func() error { return grpcServer.Serve(ctx, grpcLis) })
	}
	if gatewayLis != nil {
		eg.Go(func() error { return gatewayServer.Serve(ctx, gatewayLis) })
	}

	var wg sync.WaitGroup

	if cmuxServer != nil {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cmuxServer.Serve()
		}()
	}

	doneCh := make(chan struct{}, 1)
	wg.Add(1)
	go func() {
		defer wg.Done()
		e.watchShutdownSignal(cancel, doneCh)
	}()

	err = errors.WithStack(eg.Wait())
	close(doneCh)
	wg.Wait()

	return err
}

func (e *Engine) watchShutdownSignal(cancel context.CancelFunc, doneCh <-chan struct{}) {
	sigCh := make(chan os.Signal, 1)
	defer close(sigCh)
	defer signal.Stop(sigCh)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	for {
		select {
		case sig := <-sigCh:
			grpclog.Info("received signal: %v", sig)
			cancel()
		case <-doneCh:
			return
		}
	}
}
