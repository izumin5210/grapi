package internal

import (
	"context"
	"net"
)

// Server provides an interface for starting and stopping the server.
type Server interface {
	Serve(context.Context, net.Listener) error
}
