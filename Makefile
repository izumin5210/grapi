.DEFAULT_GOAL := all

VERSION_MAJOR ?= 0
VERSION_MINOR ?= 2
VERSION_BUILD ?= 2

VERSION ?= v$(VERSION_MAJOR).$(VERSION_MINOR).$(VERSION_BUILD)
REVISION ?= $(shell git describe --always)
BUILD_DATE ?= $(shell date +'%Y-%m-%dT%H:%M:%SZ')
RELEASE_TYPE ?= $(if $(shell git tag --contains $(REVISION) | grep $(VERSION)),stable,canary)

ORG := github.com/izumin5210
PROJECT := grapi
ROOT_PKG ?= $(ORG)/$(PROJECT)

TEMPLATE_PKG := pkg/grapicmd/internal/module/generator/template
MOCK_PKG := pkg/grapicmd/internal/module/testing
GENERATED_PKGS := $(TEMPLATE_PKG)
GENERATED_PKGS += $(MOCK_PKG)

SRC_FILES := $(shell git ls-files --cached --others --exclude-standard | grep -E "\.go$$" | grep -v ".snapshot")
GOFMT_TARGET := $(SRC_FILES)
$(foreach pkg,$(GENERATED_PKGS),$(eval GOFMT_TARGET := $(filter-out $(pkg)/%,$(GOFMT_TARGET))))
GOLINT_TARGET := $(shell go list ./... | grep -v "$(pkg)/testing")
$(foreach pkg,$(GENERATED_PKGS),$(eval GOLINT_TARGET := $(filter-out $(ROOT_PKG)/$(pkg),$(GOLINT_TARGET))))

GO_BUILD_FLAGS := -v
GO_TEST_FLAGS := -v -timeout 30s
GO_TEST_INTEGRATION_FLAGS := -v
GO_COVER_FLAGS := -coverpkg $(shell echo $(GOLINT_TARGET) | tr ' ' ',') -coverprofile coverage.txt -covermode atomic

XC_ARCH := 386 amd64
XC_OS := darwin linux windows

PATH := ${PWD}/bin:${PATH}
export PATH

#  Utils
#----------------------------------------------------------------
define section
  @printf "\e[34m--> $1\e[0m\n"
endef


#  App
#----------------------------------------------------------------
BIN_DIR := ./bin/
OUT_DIR := ./dist
GENERATED_BINS :=
PACKAGES :=
CMDS := $(wildcard ./cmd/*)

define cmd-tmpl

$(eval NAME := $(notdir $(1)))
$(eval OUT := $(addprefix $(BIN_DIR),$(NAME)))
$(eval LDFLAGS := -ldflags "-X main.name=$(NAME) -X main.version=$(VERSION) -X main.revision=$(REVISION) -X main.buildDate=$(BUILD_DATE) -X main.releaseType=$(RELEASE_TYPE)")
$(eval GENERATED_BINS += $(OUT))
$(OUT): $(SRC_FILES)
	$(call section,Building $(OUT))
	@go build $(GO_BUILD_FLAGS) $(LDFLAGS) -o $(OUT) $(1)

.PHONY: $(NAME)
$(NAME): $(OUT)

$(eval PACKAGES += $(NAME)-package)

.PHONY: $(NAME)-package
$(NAME)-package: $(NAME) $(BIN_DIR)/gox
	gox \
		$(LDFLAGS) \
		-os="$(XC_OS)" \
		-arch="$(XC_ARCH)" \
		-output="$(OUT_DIR)/$(NAME)_{{.OS}}_{{.Arch}}" \
		$(1)
endef

$(foreach src,$(CMDS),$(eval $(call cmd-tmpl,$(src))))

.PHONY: all
all: $(GENERATED_BINS)


#  Commands
#----------------------------------------------------------------
$(BIN_DIR)/mockgen $(BIN_DIR)/go-assets-builder $(BIN_DIR)/gox: $(BIN_DIR)/gex
	gex --build

.PHONY: setup
setup:
	@go get github.com/izumin5210/gex/cmd/gex
ifeq ($(shell go env GOMOD),)
ifdef CI
	curl https://raw.githubusercontent.com/golang/dep/master/install.sh | sh
endif
	dep ensure -v -vendor-only
else
	go get -v ./...
endif

.PHONY: clean
clean:
	rm -rf $(BIN_DIR)/*

.PHONY: gen
gen: $(BIN_DIR)/mockgen $(BIN_DIR)/go-assets-builder
	go generate ./...

.PHONY: lint
lint:
	$(call section,Linting)
ifdef CI
	gex reviewdog -reporter=github-pr-review
else
	gex reviewdog -diff="git diff master"
endif

.PHONY: test
test:
	$(call section,Testing)
	@go test $(GO_TEST_FLAGS) ./...

.PHONY: cover
cover:
	$(call section,Testing with coverage)
	@go test $(GO_TEST_FLAGS) $(GO_COVER_FLAGS) ./...

.PHONY: test-integration
test-integration:
	$(call section,Integration Testing)
	cd _tests && go test $(GO_TEST_INTEGRATION_FLAGS) ./...

.PHONY: packages
packages: $(PACKAGES)
