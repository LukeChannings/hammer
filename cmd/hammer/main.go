package main

import (
	"fmt"
	"log"
	"os"

	"github.com/docopt/docopt-go"
	"github.com/lukechannings/hammer/internal/bundle"
	"github.com/lukechannings/hammer/internal/serve"
)

var taggedVersion string
var gitHash string

func main() {
	usage := `Hammer.

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
`

	args, docoptErr := docopt.ParseDoc(usage)

	if docoptErr != nil {
		log.Fatal(docoptErr.Error())
	}

	if showVersion, _ := args.Bool("--version"); showVersion {
		if taggedVersion != "" {
			fmt.Println(taggedVersion)
		} else {
			fmt.Println(gitHash)
		}
		os.Exit(0)
	} else if shouldBundle, _ := args.Bool("bundle"); shouldBundle {
		entrypoint, _ := args.String("<entrypoint>")
		dest, _ := args.String("<dest>")
		minify, _ := args.Bool("--minify")
		extractCSS, _ := args.Bool("--extract-css")
		bundle.Bundle(entrypoint, dest, minify, extractCSS)
		os.Exit(0)
	} else {
		srcs, _ := args["<src>"].([]string)
		addr, _ := args.String("--addr")
		port, _ := args.String("--port")
		proxy, _ := args.String("--proxy")
		zgip, _ := args.Bool("--gzip")

		var compress serve.Compress = serve.CompressNone

		if zgip {
			compress = serve.CompressGzip
		}

		serve.Serve(srcs, fmt.Sprintf("%s:%s", addr, port), compress, proxy)
		os.Exit(0)
	}
}
