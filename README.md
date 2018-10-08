# ![grapi](./grapi.png)
[![Build Status](https://travis-ci.org/izumin5210/grapi.svg?branch=master)](https://travis-ci.org/izumin5210/grapi)
[![GoDoc](https://godoc.org/github.com/izumin5210/grapi/pkg/grapiserver?status.svg)](https://godoc.org/github.com/izumin5210/grapi/pkg/grapiserver)
[![Go Report Card](https://goreportcard.com/badge/github.com/izumin5210/grapi)](https://goreportcard.com/report/github.com/izumin5210/grapi)
[![Go project version](https://badge.fury.io/go/github.com%2Fizumin5210%2Fgrapi.svg)](https://badge.fury.io/go/github.com%2Fizumin5210%2Fgrapi)
[![license](https://img.shields.io/github/license/izumin5210/grapi.svg)](./LICENSE)

:open_mouth: A surprisingly easy API server and generator in gRPC and Go

## Features
- You can develop and deploy API servers blazingly fast :zap:
- Easy code generator
	- application  (inspired by `rails new` and `create-react-app`)
	- gRPC services and their implementations (inspired by `rails g (scaffold_)controller`)
- User-friendly `protoc` wrapper (inspired by [protoeasy](https://github.com/peter-edge/protoeasy-go))
- Provides gRPC and HTTP JSON API  with single implementation by using [grpc-gateway](https://github.com/grpc-ecosystem/grpc-gateway)
- Generates codes based on [google's API design guideline](https://cloud.google.com/apis/design/)

[![asciicast](https://asciinema.org/a/176280.png)](https://asciinema.org/a/176280)

## Getting Started

### Create a new application
```
$ grapi init awesome-app
```

### Create a new service
```
$ grapi g service books
```

Or, if you need full [standard methods](https://cloud.google.com/apis/design/standard_methods), you can get them with following command:

```
$ grapi g scaffold-service books
```

And you should register generated services to the `grapiserver.Engine` instance:

```diff
 // app/run.go
 
 // Run starts the grapiserver.
 func Run() error {
 	s := grapiserver.New(
 		grapiserver.WithDefaultLogger(),
 		grapiserver.WithServers(
+			server.NewBookServiceServer(),
-		// TODO
 		),
 	)
 	return s.Serve()
 }
```

If you updated service definition, you can re-generate `.pb.go` and `.pb.gw.go` with following command:

```
$ grapi protoc
```

### Start server

```
$ grapi server
```

### User-defined commands

```
$ grapi g command import-books
$ vim cmd/import-books/run.go  # implements the command
$ grapi import-books  # run the command
```

### Build commands (including server)

```
$ grapi build
```

## Installation

1. **grapi**
    - Linux
        - `curl -Lo grapi https://github.com/izumin5210/grapi/releases/download/v0.2.2/grapi_linux_amd64 && chmod +x grapi && sudo mv grapi /usr/local/bin`
    - macOS
        - `brew install izumin5210/tools/grapi`
    - others
        - `go get github.com/izumin5210/grapi/cmd/grapi`
1. **dep** or **Modules**
    - [dep](https://golang.github.io/dep/)
        - macOS
            - `brew install dep`
        - others
            - See [Installation · dep](https://golang.github.io/dep/docs/installation.html)
            - `curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh`
    - [Modules](https://github.com/golang/go/wiki/Modules) (experimental)
        - Use Go 1.11 and set `GO111MODULE=on` your env vars
1. **protoc**
    - macOS
        - `brew install protobuf`
    - others
        - Download and install from [google/protobuf](https://github.com/google/protobuf)
