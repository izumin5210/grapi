package main

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/mattn/go-colorable"

	"github.com/izumin5210/grapi/pkg/cli"
	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/cmd"
)

var (
	inReader            = os.Stdin
	outWriter io.Writer = os.Stdout
	errWriter io.Writer = os.Stderr
)

func init() {
	if runtime.GOOS == "windows" {
		outWriter = colorable.NewColorableStdout()
		errWriter = colorable.NewColorableStderr()
	}
}

func main() {
	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	err = cmd.NewGrapiCommand(&grapicmd.Ctx{
		IO: &cli.IO{
			In:  inReader,
			Out: outWriter,
			Err: errWriter,
		},
		RootDir: cli.RootDir(cwd),
		Build: grapicmd.BuildConfig{
			AppName:   name,
			Version:   version,
			Revision:  revision,
			BuildDate: buildDate,
			Prebuilt:  prebuilt,
		},
	}).Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
