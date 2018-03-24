package main

import (
	"os"

	"testapp/app"
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

