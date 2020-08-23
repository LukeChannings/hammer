package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"github.com/docopt/docopt-go"
	"github.com/evanw/esbuild/pkg/api"
	"github.com/lukechannings/hammer/internal/esbuild"
)

func main() {
	usage := `Hammer.

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
`

	arguments, docoptErr := docopt.ParseDoc(usage)

	if docoptErr != nil {
		log.Fatal(docoptErr.Error())
	}

	if showVersion, _ := arguments.Bool("--version"); showVersion {
		fmt.Println("v1.0.0")
		os.Exit(0)
	}

	if shouldBundle, _ := arguments.Bool("bundle"); shouldBundle {
		src, _ := arguments.String("<src>")
		dest, _ := arguments.String("<dest>")
		minify, _ := arguments.Bool("--minify")
		extractCSS, _ := arguments.Bool("--extract-css")
		bundle(src, dest, minify, extractCSS)
		os.Exit(0)
	}

	if shouldServe, _ := arguments.Bool("serve"); shouldServe {
		src, _ := arguments.String("<src>")
		host, _ := arguments.String("--host")
		port, _ := arguments.String("--port")
		addr := fmt.Sprintf("%s:%s", host, port)
		serve(src, addr)
		os.Exit(0)
	}

	fmt.Println(usage)
}

func serve(src string, addr string) {

	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		path := r.URL.Path

		if path == "/" {
			path = "/index.html"
		}

		content, err := ioutil.ReadFile(src + path)

		if err != nil {
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(fmt.Sprintf("404 - File not found (%v)", err.Error())))
		} else {
			if strings.HasSuffix(path, ".ts") || strings.HasSuffix(path, ".tsx") {
				result := api.Transform(string(content), api.TransformOptions{
					Target:    api.ES2018,
					Sourcemap: api.SourceMapInline,
					Loader:    api.LoaderTSX,
				})

				w.Header().Set("content-type", "text/javascript")

				if len(result.Errors) != 0 {
					for _, err := range result.Errors {
						fmt.Printf("ERROR: %s\n", err.Text)
					}
				}
				if len(result.Warnings) != 0 {
					for _, warn := range result.Warnings {
						fmt.Printf("WARNING: %s\n", warn.Text)
					}
				}

				w.Write(result.JS)
			} else if strings.HasSuffix(path, ".css") {
				ref := r.Referer()
				if strings.HasSuffix(ref, ".tsx") || strings.HasSuffix(ref, ".ts") || strings.HasSuffix(ref, ".js") {
					w.Header().Set("content-type", "application/javascript")
					fmt.Fprintf(w, "const s = document.createElement('style'); s.innerHTML = \n`%s`; document.head.appendChild(s)", string(content))
				} else {
					w.Header().Set("content-type", "text/css")
					w.Write(content)
				}
			} else {
				w.Write(content)
			}
		}
	}

	s := &http.Server{
		Addr:           addr,
		Handler:        handler,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	fmt.Printf("Listening on %s\n", addr)
	s.ListenAndServe()
}

func bundle(src string, dest string, minify bool, extractCSS bool) {
	fmt.Printf("Bundling %s", src)
	result := api.Build(api.BuildOptions{
		EntryPoints:       []string{src},
		Bundle:            true,
		Write:             true,
		LogLevel:          api.LogLevelInfo,
		Outfile:           dest,
		Format:            api.FormatESModule,
		MinifyIdentifiers: minify,
		MinifySyntax:      minify,
		MinifyWhitespace:  minify,
		Plugins:           []func(api.Plugin){esbuild.HTTPPlugin(true), esbuild.CSSPlugin(extractCSS)},
	})

	if len(result.Warnings) != 0 {
		for _, warning := range result.Warnings {
			log.Printf("Warning: %v\n", warning.Text)
		}
	}

	if len(result.Errors) != 0 {
		for _, error := range result.Errors {
			log.Printf("Error: %v\n", error.Text)
		}
		os.Exit(1)
	}

	if extractCSS {
		css := esbuild.GetExtractedCSS()
		cssPath := strings.Replace(dest, path.Ext(dest), ".css", 1)
		ioutil.WriteFile(cssPath, []byte(css), 0x740)
	}
}
