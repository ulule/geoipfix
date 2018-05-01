ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
VERSION=$(awk '/Version/ { gsub("\"", ""); print $NF }' ${ROOT_DIR}/application/constants.go)

BIN_DIR = $(ROOT_DIR)/bin
SHARE_DIR = $(ROOT_DIR)/share
APP_DIR = /go/src/github.com/ulule/geoipfix

branch = $(shell git rev-parse --abbrev-ref HEAD)
commit = $(shell git log --pretty=format:'%h' -n 1)
now = $(shell date "+%Y-%m-%d %T UTC%z")
compiler = $(shell go version)

test: unit

generate:
	protoc --gofast_out=plugins=grpc:. proto/geoipfix.proto

dependencies:
	dep ensure -v

run:
	GEOIPFIX_CONF=`pwd`/config.json ./bin/geoipfix

live:
	@modd

unit:
	@(go list ./... | xargs -n1 go test -v -cover)

all: geoipfix
	@(mkdir -p $(BIN_DIR))

build:
	@(echo "-> Compiling geoipfix binary")
	@(mkdir -p $(BIN_DIR))
	@(go build -o $(BIN_DIR)/geoipfix ./cmd/main.go)
	@(echo "-> geoipfix binary created")

build-static:
	@(echo "-> Creating statically linked binary...")
	@(mkdir -p $(BIN_DIR))
	@(CGO_ENABLED=0 go build -ldflags "\
		-X 'geoipfix.Branch=$(branch)' \
		-X 'geoipfix.Revision=$(commit)' \
		-X 'geoipfix.BuildTime=$(now)' \
		-X 'geoipfix.Compiler=$(compiler)'" -a -installsuffix cgo -o $(BIN_DIR)/geoipfix ./cmd/main.go)


docker-build-geoip:
	@(echo "-> Preparing builder...")
	@(docker build -t geoipfix-geoip -f Dockerfile.geoip .)
	@(mkdir -p $(SHARE_DIR))
	@(docker run --rm -v $(SHARE_DIR):/usr/share/geoip/ geoipfix-geoip)

docker-build:
	@(echo "-> Preparing builder...")
	@(docker build -t geoipfix-builder -f Dockerfile.build .)
	@(mkdir -p $(BIN_DIR))
	@(echo "-> Running geoipfix builder...")
	@(docker run --rm -v $(BIN_DIR):$(APP_DIR)/bin geoipfix-builder)

format:
	@(go fmt ./...)
	@(go vet ./...)

.PNONY: all test format
