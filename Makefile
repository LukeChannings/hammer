TAGGED_VERSION = $(shell git tag --points-at HEAD)
GIT_HASH = $(shell git rev-parse --short HEAD)
LD_FLAGS = -s -w -X main.gitHash=${GIT_HASH} -X main.taggedVersion=${TAGGED_VERSION}

build: cmd/**/*.go internal/**/*.go
	mkdir -p build/{darwin,linux,linux-386,windows}
	GOOS=darwin go build -ldflags="${LD_FLAGS}" -o ./build/darwin/hammer ./cmd/hammer
	GOOS=linux GOARCH=386 go build -ldflags="${LD_FLAGS}" -o ./build/linux-386/hammer ./cmd/hammer
	GOOS=linux GOARCH=amd64 go build -ldflags="${LD_FLAGS}" -o ./build/linux/hammer ./cmd/hammer
	GOOS=windows go build -ldflags="${LD_FLAGS}" -o ./build/windows/hammer ./cmd/hammer

clean:
	rm -rf build
