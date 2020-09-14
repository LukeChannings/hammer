package serve

import (
	"compress/gzip"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"path"
	"strings"

	"github.com/lukechannings/hammer/internal/cssmodule"
	"github.com/lukechannings/hammer/internal/esbuild"
	"github.com/lukechannings/hammer/internal/livereload"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/logrusorgru/aurora/v3"
)

var extensions []string = []string{".ts", ".tsx", ".js", ".jsx", ".mjs"}

func getLoaderForExtension(extension string) api.Loader {
	switch extension {
	case ".ts":
		return api.LoaderTS
	case ".tsx":
		return api.LoaderTSX
	case ".js", ".mjs":
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
func Serve(srcs []string, addr string, compress Compress, proxy string, cssModules bool, liveReload bool, defines map[string]string) {

	if liveReload {
		livereload.InitWatcher()
		defer livereload.CloseWatcher()
	}

	var handler http.HandlerFunc = func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		ref := r.Referer()

		if p == "/" {
			p = "/index.html"
		}

		resolvedPath, content, err := loadStaticFile(srcs, p, &extensions)

		var body []byte
		statusCode := http.StatusNotFound

		var zw *gzip.Writer = nil

		if compress == CompressGzip && strings.Contains(r.Header.Get("accept-encoding"), "gzip") {
			zw = gzip.NewWriter(w)
			w.Header().Set("content-encoding", "gzip")
		}

		if err == nil {
			if liveReload {
				livereload.AddWatcher(*resolvedPath)
			}
			statusCode = http.StatusOK
			body = *content
			switch path.Ext(*resolvedPath) {
			case ".ts", ".tsx", ".js", ".jsx", ".mjs":
				{
					result := api.Transform(string(*content), api.TransformOptions{
						Target:     api.ES2018,
						Sourcemap:  api.SourceMapInline,
						Loader:     getLoaderForExtension(path.Ext(*resolvedPath)),
						Sourcefile: *resolvedPath,
						Defines:    defines,
					})

					if len(result.Errors) != 0 {
						statusCode = http.StatusInternalServerError
						for _, err := range result.Errors {
							fmt.Printf("%s in %s:%v:%v\n> %s\n", aurora.Red(err.Text), err.Location.File, err.Location.Line, err.Location.Column, aurora.Bold(err.Location.LineText))
						}
					} else {
						if len(result.Warnings) != 0 {
							for _, warn := range result.Warnings {
								fmt.Printf("%s in %s:%v:%v\n> %s\n", aurora.Red(warn.Text), warn.Location.File, warn.Location.Line, warn.Location.Column, aurora.Bold(warn.Location.LineText))
							}
						}
						body = result.JS
					}

					w.Header().Set("content-type", "text/javascript")
				}
			case ".css", ".scss":
				{
					if hasExtension(ref, extensions...) || path.Ext(ref) == "" {
						w.Header().Set("content-type", "application/javascript")
						if cssModules {
							module, err := cssmodule.Process(strings.NewReader(string(*content)), *resolvedPath)
							if err != nil {
								statusCode = http.StatusInternalServerError
								body = []byte(fmt.Sprintf("500 - %s", err.Error()))
							}
							body = []byte(esbuild.WrapCSSForJSInjection(module.CSS, *resolvedPath, module.GetExportsString()))
						} else {
							body = []byte(esbuild.WrapCSSForJSInjection(string(*content), *resolvedPath, ""))
						}
					} else {
						w.Header().Set("content-type", "text/css")
					}
				}
			case ".html":
				{
					if liveReload {
						tmp := livereload.Inject(string(*content))
						content = &tmp
						body = *content
					}
					w.Header().Set("content-type", "text/html")
				}
			}
		} else {
			fmt.Println(err.Error())
			if len(proxy) != 0 {
				location := proxy + r.URL.Path
				w.Header().Set("Location", location)
				w.WriteHeader(http.StatusTemporaryRedirect)
			}
		}

		w.WriteHeader(statusCode)

		if zw != nil {
			zw.Write(body)
			zw.Close()
		} else {
			w.Write(body)
		}
	}

	mux := http.NewServeMux()
	if liveReload {
		mux.HandleFunc(livereload.LiveReloadPath, livereload.Handler)
	}

	mux.HandleFunc("/", handler)

	fmt.Println("Listening on ", "http://"+addr)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatalf("An error occurred trying to listen on %v - %v", addr, err.Error())
	}
}

// Searches the filesystem for a file based on the srcs and the path.
// If a duplicate file exists in multiple srcs, the first src will be used.
// If no file is found an error is returned
func loadStaticFile(srcs []string, p string, exts *[]string) (resolvedPath *string, content *[]byte, err error) {
	var searchPath string
	var found bool
	for _, src := range srcs {
		searchPath = path.Join(src, p)
		data, notFoundErr := ioutil.ReadFile(searchPath)

		if notFoundErr != nil && path.Ext(searchPath) == "" && exts != nil {
			for _, ext := range *exts {
				extData, extNotFoundErr := ioutil.ReadFile(searchPath + ext)
				if extNotFoundErr == nil {
					searchPath = searchPath + ext
					data = extData
					notFoundErr = nil
					break
				}
			}
		}

		if notFoundErr == nil {
			resolvedPath = &searchPath
			content = &data
			found = true
			return
		}
	}

	if !found {
		err = errors.New("No file found for " + p)
	}

	return
}
