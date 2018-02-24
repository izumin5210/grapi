package main

import (
	"fmt"
	"os"

	"github.com/izumin5210/grapi/pkg/grapicmd/cmd"
	"github.com/spf13/afero"
)

func main() {
	err := cmd.NewGrapiCommand(afero.NewOsFs(), os.Stdin, os.Stdout, os.Stderr).Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
