# ![grapi](./grapi.png)
[![CI](https://github.com/izumin5210/grapi/workflows/CI/badge.svg)](https://github.com/izumin5210/grapi/actions?workflow=CI)
[![GoDoc](https://godoc.org/github.com/izumin5210/grapi/pkg/grapiserver?status.svg)](https://godoc.org/github.com/izumin5210/grapi/pkg/grapiserver)
[![Go Report Card](https://goreportcard.com/badge/github.com/izumin5210/grapi)](https://goreportcard.com/report/github.com/izumin5210/grapi)
[![GitHub release (latest SemVer)](https://img.shields.io/github/v/release/izumin5210/grapi)](http://github.com/izumin5210/grapi/releases/latest)
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

## :warning: Migrate 0.4.x -> 0.5.x :warning:
[grapiserver](https://godoc.org/github.com/izumin5210/grapi/pkg/grapiserver) will not handle os signals from v0.5.x.
We recommend to use [`appctx.Global()`](https://godoc.org/github.com/srvc/appctx#Global) if you want to handle them.

<details>
<summary>:memo: How to migrate</summary>

0. Bump grapi version
    - `go get -u github.com/izuimn5210/grapi@v0.5'
1. Update `cmd/server/run.go`
    - ```diff
       	// Application context
      -	ctx := context.Background()
      +	ctx := appctx.Global()
      ```
    - ```diff
      -	return s.ServeContext(ctx)
      +	return s.Serve(ctx)
      ```

</details>


## :warning: Migrate 0.3.x -> 0.4.x :warning:
Some tools that are depended by grapi are updated. If you have a grapi project <=v0.3.x, you should migrate it.

<details>
<summary>:memo: How to migrate</summary>

0. Bump grapi version
    - If you use [dep](https://golang.github.io/dep/), update `Gopkg.toml`
      ```diff
       [[constraint]]
         name = "github.com/izumin5210/grapi"
      -  version = "0.3.0"
      +  version = "0.4.0"
      ```
    - and run `dep ensure`
1. Update [gex](https://github.com/izumin5210/gex) and `tools.go`
    - ```
      go get -u github.com/izumin5210/gex/cmd/gex
      gex --regen
      ```
1. Initialize [Go Modules](https://github.com/golang/go/wiki/Modules)
    - ```
      go mod init
      go mod tidy
      ```
1. Update `grapi.toml`
    - ```diff
      package = "yourcompany.yourappname"
      
      [grapi]
      server_dir = "./app/server"
   
      [protoc]
      protos_dir = "./api/protos"
      out_dir = "./api"
      import_dirs = [
        "./api/protos",
      -  "./vendor/github.com/grpc-ecosystem/grpc-gateway",
      -  "./vendor/github.com/grpc-ecosystem/grpc-gateway/third_party/googleapis",
      +  '{{ module "github.com/grpc-ecosystem/grpc-gateway" }}',
      +  '{{ module "github.com/grpc-ecosystem/grpc-gateway" }}/third_party/googleapis',
      ]
   
        [[protoc.plugins]]
        name = "go"
        args = { plugins = "grpc", paths = "source_relative" }
   
        [[protoc.plugins]]
        name = "grpc-gateway"
        args = { logtostderr = true, paths = "source_relative" }
   
        [[protoc.plugins]]
        name = "swagger"
        args = { logtostderr = true }
      ```
1. Drop dep
    - ```
      rm Gopkg.*
      ```

	
</details>

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

## Getting Started

### Create a new application
```
$ grapi init awesome-app
```

### Create a new service
First you need to move to the application.

```
$ cd awesome-app/
```

Then you can generate service.

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
