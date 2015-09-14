ROOT_DIR := $(shell dirname $(realpath $(lastword $(MAKEFILE_LIST))))
VERSION=$(awk '/Version/ { gsub("\"", ""); print $NF }' ${ROOT_DIR}/application/constants.go)

BIN_DIR = $(ROOT_DIR)/bin

test: unit

unit:
	@(go list ./... | xargs -n1 go test -v -cover)

deps:
	godep save ./...

all: ipfix
	@(mkdir -p $(BIN_DIR))

build:
	@(go get github.com/tools/godep)
	@(godep restore)
	@(echo "-> Compiling ipfix binary")
	@(mkdir -p $(BIN_DIR))
	@(godep go build -o $(BIN_DIR)/ipfix)
	@(echo "-> ipfix binary created")

build-static:
	@(echo "-> Creating statically linked binary...")
	@(go get github.com/tools/godep)
	mkdir -p $(BIN_DIR)
	CGO_ENABLED=0 godep go build -a -installsuffix cgo -o $(BIN_DIR)/ulule-api

format:
	@(go fmt ./...)
	@(go vet ./...)

.PNONY: all test format

