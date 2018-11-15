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


## :warning: Migrate 0.2.x -> 0.3.x :warning:
grapi v0.3.0 has some breaking changes. If you have a grapi project <=v0.2.x, you should migrate it.

<details>
<summary>:memo: How to migrate</summary>

0. Bump grapi version
    - If you use [dep](https://golang.github.io/dep/), update `Gopkg.toml`
      ```diff
       [[constraint]]
         name = "github.com/izumin5210/grapi"
      -  version = "0.2.2"
      +  version = "0.3.0"
      ```
    - and run `dep ensure`
1. Introduce [gex](https://github.com/izumin5210/gex)
    - ```
      go get github.com/izumin5210/gex/cmd/gex
      ```
1. Add defualt generator plugins:
    - ```
      gex \
        --add github.com/izumin5210/grapi/cmd/grapi \
        --add github.com/izumin5210/grapi/cmd/grapi-gen-command \
        --add github.com/izumin5210/grapi/cmd/grapi-gen-service \
        --add github.com/izumin5210/grapi/cmd/grapi-gen-scaffold-service \
        --add github.com/izumin5210/grapi/cmd/grapi-gen-type
      ```
1. Add protoc plugins via gex
    - ```
      gex \
        --add github.com/golang/protobuf/protoc-gen-go \
        --add github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway \
        --add github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger
      ```
    - Remove protoc plugins from `Gopkg.toml`
      ```diff
      -required = [
      -  "github.com/golang/protobuf/protoc-gen-go",
      -  "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway",
      -  "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger",
      -]
      ```
1. Update `grapi.toml`
    - ```diff
      +package = "yourcompany.yourappname"
      +
       [grapi]
       server_dir = "./app/server"

       [protoc]
       protos_dir = "./api/protos"
       out_dir = "./api"
       import_dirs = [
      +  "./api/protos",
         "./vendor/github.com/grpc-ecosystem/grpc-gateway",
         "./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis",
       ]

         [[protoc.plugins]]
      -  path = "./vendor/github.com/golang/protobuf/protoc-gen-go"
         name = "go"
         args = { plugins = "grpc", paths = "source_relative" }

         [[protoc.plugins]]
      -  path = "./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
         name = "grpc-gateway"
      -  args = { logtostderr = true }
      +  args = { logtostderr = true, paths = "source_relative" }

         [[protoc.plugins]]
      -  path = "./vendor/github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger"
         name = "swagger"
         args = { logtostderr = true }
      ```

</details>


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
