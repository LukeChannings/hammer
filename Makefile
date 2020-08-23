TAGGED_VERSION = $(shell git tag --points-at HEAD)
GIT_HASH = $(shell git rev-parse --short HEAD)
LD_FLAGS = -s -w -X main.gitHash=${GIT_HASH} -X main.taggedVersion=${TAGGED_VERSION}

build:
	go build -ldflags="${LD_FLAGS}" ./cmd/hammer
