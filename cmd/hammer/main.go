package main

import (
	"fmt"
	"log"
	"os"

	"github.com/evanw/esbuild/pkg/api"

	"github.com/docopt/docopt-go"
	"github.com/lukechannings/hammer/internal/bundle"
	"github.com/lukechannings/hammer/internal/serve"
	"github.com/lukechannings/hammer/internal/trace"
)

var taggedVersion string
var gitHash string

func main() {
	usage := `hammer.

Usage:
  hammer serve <src>... [-p=<port>] [-a=<host>] [--gzip] [--proxy=<url>] [--css-modules]
  hammer bundle <entrypoint> <dest> [--minify] [--sourcemap=<external|inline|none>] [--extract-css] [--css-modules]
  hammer dependency-graph <entrypoint> [--flat|--list-orphans]
  hammer -h | --help
  hammer --version

Options:
  -h --help                              Show this screen.
  --version                              Show version.
  -p --port=<port>                       The HTTP Server port [default: 4321]
  -a --addr=<host>                       The default IP for the server port [default: 0.0.0.0]
  -g --gzip                              Compress the output with gzip. Note: Not recommended for local development.
  -P --proxy=<url>                       Redirect 404s to a proxy URL
  --sourcemap=<external|inline|none>     Whether or not to include a source map with the bundle
  --css-modules                          Enable CSS Modules

Dependency Graph Options:
  --flat                                 A flat list of all files in the dependency graph
  --list-orphans                         Lists all files that are never imported
`

	var version string = gitHash
	if taggedVersion != "" {
		version = taggedVersion
	}

	args, docoptErr := docopt.ParseArgs(usage, os.Args[1:], version)

	if docoptErr != nil {
		log.Fatal(docoptErr.Error())
	}

	cssModules, _ := args.Bool("--css-modules")

	if shouldBundle, _ := args.Bool("bundle"); shouldBundle {
		entrypoint, _ := args.String("<entrypoint>")
		dest, _ := args.String("<dest>")
		minify, _ := args.Bool("--minify")
		extractCSS, _ := args.Bool("--extract-css")
		smap, _ := args.String("--sourcemap")
		var sourcemap api.SourceMap
		if smap == "external" {
			sourcemap = api.SourceMapExternal
		} else if smap == "inline" {
			sourcemap = api.SourceMapInline
		} else {
			sourcemap = api.SourceMapNone
		}
		bundle.Bundle(entrypoint, dest, minify, sourcemap, extractCSS, cssModules)
		os.Exit(0)
	} else if shouldServe, _ := args.Bool("serve"); shouldServe {
		srcs, _ := args["<src>"].([]string)
		addr, _ := args.String("--addr")
		port, _ := args.String("--port")
		proxy, _ := args.String("--proxy")
		zgip, _ := args.Bool("--gzip")

		var compress serve.Compress = serve.CompressNone

		if zgip {
			compress = serve.CompressGzip
		}

		serve.Serve(srcs, fmt.Sprintf("%s:%s", addr, port), compress, proxy, cssModules)
		os.Exit(0)
	} else if dependencyGraph, _ := args.Bool("dependency-graph"); dependencyGraph {
		entrypoint, _ := args.String("<entrypoint>")
		flat, _ := args.Bool("--flat")
		listOrphans, _ := args.Bool("--list-orphans")
		trace.Trace(entrypoint, flat, listOrphans)
	} else {
		fmt.Print(usage)
		os.Exit(1)
	}
}
