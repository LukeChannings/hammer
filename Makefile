TAGGED_VERSION = $(shell git tag --points-at HEAD)
GIT_HASH = $(shell git rev-parse --short HEAD)
LD_FLAGS = -s -w -X main.gitHash=${GIT_HASH} -X main.taggedVersion=${TAGGED_VERSION}

build: cmd/**/*.go internal/**/*.go
	pkger -o cmd/hammer
	mkdir -p build/{darwin/{amd64},linux/{amd64,386,arm64},windows/{amd64}}

	GOOS=darwin GOARCH=amd64 go build -ldflags="${LD_FLAGS}" -o ./build/darwin/amd64/hammer ./cmd/hammer
	tar -czvf build/hammer-darwin-amd64.tar.gz -C build/darwin/amd64 hammer

	GOOS=linux GOARCH=amd64 go build -ldflags="${LD_FLAGS}" -o ./build/linux/amd64/hammer ./cmd/hammer
	tar -czvf build/hammer-linux-amd64.tar.gz -C build/linux/amd64 hammer

	GOOS=linux GOARCH=386 go build -ldflags="${LD_FLAGS}" -o ./build/linux/386/hammer ./cmd/hammer
	tar -czvf build/hammer-linux-386.tar.gz -C build/linux/386 hammer

	GOOS=linux GOARCH=arm64 go build -ldflags="${LD_FLAGS}" -o ./build/linux/arm64/hammer ./cmd/hammer
	tar -czvf build/hammer-linux-arm64.tar.gz -C build/linux/arm64 hammer

	GOOS=windows GOARCH=amd64 go build -ldflags="${LD_FLAGS}" -o ./build/windows/amd64/hammer ./cmd/hammer
	tar -czvf build/hammer-windows-amd64.tar.gz -C build/windows/amd64 hammer

clean:
	rm -rf build
