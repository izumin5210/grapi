package main

const (
	name    = "grapi"
	version = "v0.2.2"
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
