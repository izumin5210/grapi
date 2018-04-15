package main

import (
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/cmd"
	"github.com/mattn/go-colorable"
)

var (
	name      string
	version   string
	revision  string
	buildDate string

	releaseType string

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
	err = cmd.NewGrapiCommand(grapicmd.NewConfig(
		cwd,
		name,
		version,
		revision,
		buildDate,
		releaseType,
		inReader,
		outWriter,
		errWriter,
	)).Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
