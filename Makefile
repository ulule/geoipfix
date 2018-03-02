ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
VERSION=$(awk '/Version/ { gsub("\"", ""); print $NF }' ${ROOT_DIR}/application/constants.go)

BIN_DIR = $(ROOT_DIR)/bin
APP_DIR = /go/src/github.com/ulule/ipfix

branch = $(shell git rev-parse --abbrev-ref HEAD)
commit = $(shell git log --pretty=format:'%h' -n 1)
now = $(shell date "+%Y-%m-%d %T UTC%z")
compiler = $(shell go version)

test: unit

dependencies:
	dep ensure -v

run:
	IPFIX_CONF=`pwd`/config.json ./bin/ipfix

live:
	@modd

unit:
	@(go list ./... | xargs -n1 go test -v -cover)

all: ipfix
	@(mkdir -p $(BIN_DIR))

build:
	@(echo "-> Compiling ipfix binary")
	@(mkdir -p $(BIN_DIR))
	@(go build -o $(BIN_DIR)/ipfix ./cmd/main.go)
	@(echo "-> ipfix binary created")

build-static:
	@(echo "-> Creating statically linked binary...")
	@(mkdir -p $(BIN_DIR))
	@(CGO_ENABLED=0 go build -ldflags "-X 'main.branch=$(branch)' -X 'main.sha=$(commit)'  -X 'main.now=$(now)' -X 'main.compiler=$(compiler)'" -a -installsuffix cgo -o $(BIN_DIR)/ipfix)

docker-build:
	@(echo "-> Preparing builder...")
	@(docker build -t ipfix-builder -f Dockerfile.build .)
	@(mkdir -p $(BIN_DIR))
	@(echo "-> Running ipfix builder...")
	@(docker run --rm -v $(BIN_DIR):$(APP_DIR)/bin ipfix-builder)

format:
	@(go fmt ./...)
	@(go vet ./...)

.PNONY: all test format
