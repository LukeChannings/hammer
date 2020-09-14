package main

import (
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/docopt/docopt-go"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/lukechannings/hammer/internal/bundle"
	"github.com/lukechannings/hammer/internal/serve"
	"github.com/lukechannings/hammer/internal/template"
	"github.com/lukechannings/hammer/internal/trace"
)

var taggedVersion string
var gitHash string

func main() {
	usage := `hammer.

Usage:
  hammer new <name> [-t=<framework>] [--typescript]
  hammer serve <src>... [-p=<port>] [-a=<host>] [--gzip] [--define=<K>=<V>...] [--proxy=<url>] [--css-modules] [--live-reload]
  hammer bundle <src> <dest> [--minify] [--define=<K>=<V>...] [--sourcemap=<external|inline|none>] [--extract-css] [--css-modules]
  hammer trace <src> [--flat|--orphans] [--json]
  hammer -h | --help
  hammer --version

Options:
  -h --help                              Show this screen.
  --version                              Show version.
  -p --port=<port>                       The HTTP Server port [default: 4321]
  -a --addr=<host>                       The default IP for the server port [default: 0.0.0.0]
  --sourcemap=<external|inline|none>     Whether or not to include a source map with the bundle
  --css-modules                          Enable CSS Modules
  --minify                               Minify output
  --define=<K>=<V>                       Substitute K with V while parsing

New Options:
  -t --framework=<framework>             A JavaScript framework to use. Options: none, react [default: none]
  --ts, --typescript                     Use TypeScript

Serve Options:
  -g --gzip                              Compress the output with gzip. Note: Not recommended for local development
  -P --proxy=<url>                       Redirect 404s to a proxy URL
  -w --live-reload                       Watches for changes and sends the changes to the browser

Trace Options:
  --flat                                 List all files in the dependency graph
  --orphans                              Lists all source files in the source path that are never imported in the dependency graph
  --json                                 Output graph as json
`

	var version string = gitHash
	if taggedVersion != "" {
		version = taggedVersion
	}

	args, docoptErr := docopt.ParseArgs(usage, os.Args[1:], version)

	if docoptErr != nil {
		log.Fatal(docoptErr.Error())
	}

	var command string

	if ok, _ := args.Bool("new"); ok {
		command = "new"
	}
	if ok, _ := args.Bool("serve"); ok {
		command = "serve"
	}
	if ok, _ := args.Bool("bundle"); ok {
		command = "bundle"
	}
	if ok, _ := args.Bool("trace"); ok {
		command = "trace"
	}

	srcs, _ := args["<src>"].([]string)
	dest, _ := args.String("<dest>")
	cssModules, _ := args.Bool("--css-modules")
	minify, _ := args.Bool("--minify")
	extractCSS, _ := args.Bool("--extract-css")
	var sourceMap api.SourceMap
	switch flag, _ := args.String("--sourcemap"); flag {
	case "external":
		{
			sourceMap = api.SourceMapExternal
		}
	case "inline":
		{
			sourceMap = api.SourceMapInline
		}
	case "linked":
		{
			sourceMap = api.SourceMapLinked
		}
	default:
		{
			sourceMap = api.SourceMapNone
		}
	}
	port, _ := args.String("--port")
	addr, _ := args.String("--addr")
	name, _ := args.String("<name>")
	framework, _ := args.String("<framework>")
	typescript, _ := args.Bool("--typescript")
	gzip, _ := args.Bool("--gzip")
	proxy, _ := args.String("--proxy")
	flat, _ := args.Bool("--flat")
	orphans, _ := args.Bool("--orphans")
	outputJSON, _ := args.Bool("--json")
	liveReload, _ := args.Bool("--live-reload")
	define := args["--define"].([]string)
	var defines map[string]string = make(map[string]string)

	if define != nil && len(define) != 0 {
		for _, def := range define {
			pair := strings.Split(def, "=")
			if len(pair) == 2 {
				defines[pair[0]] = pair[1]
			}
		}
	}

	var compress serve.Compress

	if gzip {
		compress = serve.CompressGzip
	} else {
		compress = serve.CompressNone
	}

	switch command {
	case "new":
		template.New(name, framework, typescript)
	case "serve":
		serve.Serve(srcs, fmt.Sprintf("%s:%s", addr, port), compress, proxy, cssModules, liveReload, defines)
	case "bundle":
		bundle.Bundle(srcs[0], dest, minify, sourceMap, extractCSS, cssModules, defines)
	case "trace":
		trace.Trace(srcs[0], flat, orphans, outputJSON)
	}
}
