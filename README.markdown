# Hammer

Hammer uses [esbuild](https://github.com/evanw/esbuild) to transparently transform TypeScript. Use ECMAScript modules with [Pika](https://www.pika.dev) or [Unpkg](https://unpkg.com) and drop NPM or Yarn from your toolchain.

## Usage

```
Usage:
  hammer bundle <src> <dest> [--minify] [--extract-css]
  hammer serve <src> [--port=<port>] [--host=<host>]
  hammer -h | --help
  hammer --version

Options:
  -h --help            Show this screen.
  --version            Show version.
  --port=<port>        The HTTP Server port [default: 4321]
  --host=<host>        The default IP for the server port [default: 0.0.0.0]
```
