PATH := ${PWD}/bin:${PATH}
export PATH

.DEFAULT_GOAL := build

.PHONY: tools
tools:
	go genereate ./tools.go

.PHONY: clean
clean:
	rm -rf ./bin/*

.PHONY: gen
gen: tools
	go generate ./...

.PHONY: grapi
build:
	go build -v -o ./bin/grapi ./cmd/grapi

.PHONY: lint
lint: tools
	reviewdog -diff="git diff master"

.PHONY: test
test:
	go test -v ./...

.PHONY: cover
cover:
	go test -v -coverprofile coverage.txt -covermode atomic ./...

.PHONY: test-e2e
test-e2e: build
	go test -v -timeout 4m ./_tests/e2e --grapi=$$PWD/bin/grapi --revision="$(TARGET_REVISION)"
