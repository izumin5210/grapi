package di

import (
	"github.com/izumin5210/grapi/pkg/gencmd"
	"github.com/izumin5210/grapi/pkg/protoc"
)

type CreateAppFunc func(*gencmd.Command) (*App, error)

type App struct {
	Protoc protoc.Wrapper
}
