# 🔨 Hammer

Develop web apps with React, TypeScript, CSS Modules, without configuring any tooling or compromising speed of development.
Hammer is written in Go, doesn't require Node.js or NPM, and compiles on-the-fly at native speed.

## Usage

```
Usage:
  hammer serve <src>... [-p=<port>] [-a=<host>] [--gzip] [--define=<K>=<V>...] [--proxy=<url>] [--css-modules]
  hammer bundle <src> <dest> [--minify] [--define=<K>=<V>...] [--sourcemap=<external|inline|none>] [--extract-css] [--css-modules]
  hammer trace <src> [--flat|--orphans] [--json]
  hammer -h | --help
  hammer --version

Options:
  -h --help                              Show this screen.
  --version                              Show version.
  -p --port=<port>                       The HTTP Server port [default: 4321]
  -a --addr=<host>                       The default IP for the server port [default: 0.0.0.0]
  -g --gzip                              Compress the output with gzip. Note: Not recommended for local development
  -P --proxy=<url>                       Redirect 404s to a proxy URL
  --sourcemap=<external|inline|none>     Whether or not to include a source map with the bundle
  --css-modules                          Enable CSS Modules
  --minify                               Minify output
  --define=<K>=<V>                       Substitute K with V while parsing

Trace Options:
  --flat                                 List all files in the dependency graph
	--orphans                              Lists all source files in the source path that are never imported in the dependency graph
  --json                                 Output graph as json
```

## Installation

Have a look at [Releases](https://github.com/LukeChannings/hammer/releases/) and download the binary for your architecture and platform

Or install with Go:

```
go get github.com/lukechannings/hammer/cmd/hammer
```
