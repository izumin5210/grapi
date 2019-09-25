// +build tools

package tools

// tool dependencies
import (
	_ "github.com/golang/protobuf/protoc-gen-go"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/protoc-gen-swagger"
	_ "github.com/izumin5210/grapi/cmd/grapi"
	_ "github.com/izumin5210/grapi/cmd/grapi-gen-command"
	_ "github.com/izumin5210/grapi/cmd/grapi-gen-scaffold-service"
	_ "github.com/izumin5210/grapi/cmd/grapi-gen-service"
	_ "github.com/izumin5210/grapi/cmd/grapi-gen-type"
)
