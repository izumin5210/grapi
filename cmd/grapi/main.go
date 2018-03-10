package main

import (
	"fmt"
	"os"

	"github.com/izumin5210/grapi/pkg/grapicmd"
	"github.com/izumin5210/grapi/pkg/grapicmd/cmd"
)

var (
	name      string
	version   string
	revision  string
	buildDate string

	releaseType string

	inReader  = os.Stdin
	outWriter = os.Stdout
	errWriter = os.Stderr
)

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
