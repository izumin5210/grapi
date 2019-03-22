ARG GO_VERSION

#  Build
#----------------------------------------------------------------
FROM golang:$GO_VERSION-alpine as builder

RUN apk add -U --no-cache curl build-base git
ENV GO111MODULE on

WORKDIR /go/src/github.com/izumin5210/grapi
COPY . .

RUN make setup && make && cp ./bin/* /go/bin/


#  Test
#----------------------------------------------------------------
FROM golang:$GO_VERSION-alpine

RUN apk add -U --no-cache \
  build-base \
  curl \
  git \
  protobuf \
  protobuf-dev
COPY --from=builder /go/bin/* /go/bin/

WORKDIR /go/src/e2e
