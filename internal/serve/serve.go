package serve

import (
	"compress/gzip"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"
	"time"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/lukechannings/hammer/internal/esbuild"
)

var extensions []string = []string{".ts", ".tsx", ".js", ".jsx", ".mjs"}

func getLoaderForExtension(extension string) api.Loader {
	switch extension {
	case ".ts":
		return api.LoaderTS
	case ".tsx":
		return api.LoaderTSX
	case ".js":
	case ".mjs":
		return api.LoaderJS
	case ".jsx":
		return api.LoaderJSX
	}

	return api.LoaderNone
}

func hasExtension(path string, exts ...string) bool {
	for _, ext := range exts {
		if strings.HasSuffix(path, ext) {
			return true
		}
	}
	return false
}

type Compress uint8

const (
	CompressNone Compress = iota
	CompressGzip
)

// Serve - starts an HTTP server to serve static sources and transform TS and JS files on-the-fly.
func Serve(srcs []string, addr string, compress Compress, proxy string) {

	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		ref := r.Referer()

		if p == "/" {
			p = "/index.html"
		}

		var content []byte

		for _, src := range srcs {
			fullpath := path.Join(src, p)
			data, err := ioutil.ReadFile(fullpath)

			if err == nil {
				content = data
				break
			} else if path.Ext(fullpath) == "" {
				fmt.Printf("Couldn't find %v. Trying extensions\n", p)
				for _, ext := range extensions {
					data, err := ioutil.ReadFile(fullpath + ext)
					fmt.Printf("Looking for %v\n", fullpath+ext)
					if err == nil {
						content = data
						p += ext
						break
					}
				}
			}
		}

		var zw *gzip.Writer = nil

		if compress == CompressGzip && strings.Contains(r.Header.Get("accept-encoding"), "gzip") {
			zw = gzip.NewWriter(w)
			w.Header().Set("content-encoding", "gzip")
		}

		if len(content) == 0 {
			if len(proxy) != 0 {
				location := proxy + r.URL.Path
				w.Header().Set("Location", location)
				w.WriteHeader(http.StatusTemporaryRedirect)
				return
			}

			w.WriteHeader(http.StatusNotFound)
			if zw != nil {
				zw.Write([]byte(fmt.Sprintf("404 - File not found")))
				zw.Close()
			} else {
				w.Write([]byte(fmt.Sprintf("404 - File not found")))
			}
		} else {
			if hasExtension(p, extensions...) {
				result := api.Transform(string(content), api.TransformOptions{
					Target:     api.ES2018,
					Sourcemap:  api.SourceMapInline,
					Loader:     getLoaderForExtension(path.Ext(p)),
					Sourcefile: p,
				})

				w.Header().Set("content-type", "text/javascript")

				if len(result.Errors) != 0 {
					for _, err := range result.Errors {
						fmt.Printf("Error: %s\n", err.Text)
					}
				}
				if len(result.Warnings) != 0 {
					for _, warn := range result.Warnings {
						fmt.Printf("Warning: %s\n", warn.Text)
					}
				}

				if zw != nil {
					zw.Write(result.JS)
					zw.Close()
				} else {
					w.Write(result.JS)
				}
			} else if strings.HasSuffix(p, ".css") {
				if hasExtension(ref, extensions...) || path.Ext(ref) == "" {
					w.Header().Set("content-type", "application/javascript")
					fmt.Fprintf(w, esbuild.WrapCSSForJSInjection(string(content), p))
				} else {
					w.Header().Set("content-type", "text/css")
					if zw != nil {
						zw.Write(content)
						zw.Close()
					} else {
						w.Write(content)
					}
				}
			} else {
				if zw != nil {
					zw.Write(content)
					zw.Close()
				} else {
					w.Write(content)
				}
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

	fmt.Println("Listening on ", addr)

	if err := s.ListenAndServe(); err != nil {
		log.Fatalf("An error occurred trying to listen on %v - %v", addr, err.Error())
	}
}
