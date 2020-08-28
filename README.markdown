# 🔨 Hammer

An HTTP server that transparently transforms TypeScript and JavaScript using [esbuild](https://github.com/evanw/esbuild). Because Go is native, this only adds a few milliseconds to the HTTP request. It mirrors the general switches from [http-server](https://github.com/http-party/http-server), and supports multiple sources.

## Usage

```
  hammer serve <src>... [-p=<port>] [-a=<host>] [--gzip] [--proxy=<url>]
  hammer bundle <entrypoint> <dest> [--minify] [--sourcemap=<external|internal|none>] [--extract-css]
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
```

## Installation

Have a look at [Releases](https://github.com/LukeChannings/hammer/releases/) and download the binary for your architecture and platform

## Example

The example folder shows an extremely simple React project that imports React and ReactDOM from [Pika CDN](https://www.pika.dev/), and import a local CSS file in the TS component.

**What's happening?**

1. Hammer serves static files in a directory you point it to. e.g. `hammer ./example`
2. If the browser loads a file of a particular type: `.ts`, `.tsx`, `.js`, `.jsx`, `.mjs`, then they are transformed with esbuild. Because esbuild is native, the transformation takes only 1 or 2 ms on average per request.
3. If a CSS file is requested and the referer (the resource that triggered the request) is one of the extensions listed above, Hammer wraps the CSS in some JavaScript injection code and returns it as `text/javascript`

The experience the developer gets from this is that the code is loading natively. The transformations are minimal, the browser is still using ECMAScript Modules, but you can use some newer JS syntax, TypeScript, and import CSS.

## But why tho?

If you want a native development tool that allows you quickly build out a TypeScript + React project with no Webpack or Babel or even Node & NPM installed, this tool might be for you.
