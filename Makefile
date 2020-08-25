TAGGED_VERSION = $(shell git tag --points-at HEAD)
GIT_HASH = $(shell git rev-parse --short HEAD)
LD_FLAGS = -s -w -X main.gitHash=${GIT_HASH} -X main.taggedVersion=${TAGGED_VERSION}

build: build/linux-i386/hammer build/linux-amd64/hammer build/darwin/hammer

build/linux-i386/hammer:
	mkdir -p build/linux-i386
	GOARCH=386 GOOS=linux go build -ldflags="${LD_FLAGS}" -o ./build/linux-i386/hammer ./cmd/hammer

build/linux-amd64/hammer:
	mkdir -p build/linux-amd64
	GOARCH=amd64 GOOS=linux go build -ldflags="${LD_FLAGS}" -o ./build/linux-amd64/hammer ./cmd/hammer

build/darwin/hammer:
	mkdir -p build/darwin
	GOOS=darwin go build -ldflags="${LD_FLAGS}" -o ./build/darwin/hammer ./cmd/hammer


clean:
	rm -rf build
