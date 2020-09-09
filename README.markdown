# 🔨 Hammer

Develop web apps with React, TypeScript, CSS Modules, without configuring any tooling or compromising speed of development.
Hammer is written in Go, doesn't require Node.js or NPM, and compiles on-the-fly at native speed.

## Usage

```
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
```

## Installation

Have a look at [Releases](https://github.com/LukeChannings/hammer/releases/) and download the binary for your architecture and platform

Or install with Go:

```
go get github.com/lukechannings/hammer/cmd/hammer
```
