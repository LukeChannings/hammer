# Hammer

An HTTP server that transparently transforms TypeScript and JavaScript using [esbuild](https://github.com/evanw/esbuild). Because Go is native, this only adds a few milliseconds to the HTTP request. It mirrors the general switches from [http-server](https://github.com/http-party/http-server), and supports multiple sources.

## Usage

```
Usage:
  hammer <src>... [-p=<port>] [-a=<host>] [--gzip] [--proxy=<url>]
  hammer bundle <entrypoint> <dest> [--minify] [--extract-css]
  hammer -h | --help
  hammer --version

Options:
  -h --help            Show this screen.
  --version            Show version.
  -p --port=<port>     The HTTP Server port [default: 4321]
  -a --addr=<host>     The default IP for the server port [default: 0.0.0.0]
  -g --gzip            Compress the output with gzip. Note: Not recommended for local development.
  -P --proxy=<url>     Redirect 404s to a proxy URL
```

## Installation

Have a look at [Releases](https://github.com/LukeChannings/hammer/releases/) and download the binary for your architecture and platform
