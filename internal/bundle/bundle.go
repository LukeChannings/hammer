package bundle

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/lukechannings/hammer/internal/esbuild"
)

// Bundle - combines all dependencies of entrypoint into a single file
func Bundle(entrypoint string, dest string, minify bool, sourcemap api.SourceMap, extractCSS bool) {
	fmt.Printf("Bundling %s", entrypoint)
	result := api.Build(api.BuildOptions{
		EntryPoints:       []string{entrypoint},
		Bundle:            true,
		Write:             true,
		LogLevel:          api.LogLevelInfo,
		Outfile:           dest,
		Format:            api.FormatESModule,
		MinifyIdentifiers: minify,
		MinifySyntax:      minify,
		MinifyWhitespace:  minify,
		Sourcemap:         sourcemap,
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
