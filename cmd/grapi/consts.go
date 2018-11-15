package main

const (
	name    = "grapi"
	version = "v0.3.0"
)

var (
	prebuilt bool

	// set via ldflags
	revision  string
	buildDate string
)

func init() {
	prebuilt = revision != ""
}
