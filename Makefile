TAGGED_VERSION = $(shell git tag --points-at HEAD)
GIT_HASH = $(shell git rev-parse --short HEAD)
LD_FLAGS = -s -w -X main.gitHash=${GIT_HASH} -X main.taggedVersion=${TAGGED_VERSION}

build: build/linux-i386/hammer build/linux-amd64/hammer

build/linux-i386/hammer:
	mkdir -p build/linux-i386
	go build -ldflags="${LD_FLAGS}" -o ./build/linux-i386/hammer ./cmd/hammer

build/linux-amd64/hammer:
	mkdir -p build/linux-amd64
	go build -ldflags="${LD_FLAGS}" -o ./build/linux-amd64/hammer ./cmd/hammer

clean:
	rm -rf build
