package main

import (
	"os"

	"github.com/izumin5210/grapi/testing/app"
)

func main() {
	os.Exit(run())
}

func run() int {
	err := app.Run()
	if err != nil {
		return 1
	}
	return 0
}
