package app

import (
	"github.com/izumin5210/grapi/pkg/grapiserver"
)

// Run starts the grapiserver.
func Run() error {
	s := grapiserver.NewServer(
		grapiserver.WithDefaultLogger(),
		grapiserver.WithServers(
		// TODO
		),
	)
	return s.Serve()
}
