package internal

import (
	"context"
	"net"
	"sync"
)

// Server provides an interface for starting and stopping the server.
type Server interface {
	Serve(context.Context, net.Listener) error
	Shutdown()
}

// StartServer is a helper function to start the server and handle context to shutdown it.
func StartServer(ctx context.Context, serve func() error, shutdown func()) error {
	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		<-ctx.Done()
		shutdown()
	}()

	err := serve()

	wg.Wait()

	return err
}
