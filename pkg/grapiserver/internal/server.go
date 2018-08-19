package internal

import (
	"net"
)

// Server provides an interface for starting and stopping the server.
type Server interface {
	Serve(l net.Listener) error
	Shutdown()
}
