package internal

import (
	"net"
	"sync"
)

// Server provides an interface for starting and stopping the server.
type Server interface {
	Serve(l net.Listener, wg *sync.WaitGroup)
	Shutdown()
}
