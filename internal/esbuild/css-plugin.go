package esbuild

import (
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/lukechannings/hammer/internal/cssmodule"
)

var extracted string

// CSSPlugin - handles CSS imports
func CSSPlugin(extract bool, modules bool) func(api.Plugin) {
	return func(plugin api.Plugin) {
		plugin.SetName("css-plugin")

		plugin.AddResolver(api.ResolverOptions{Filter: ".css$"},
			func(args api.ResolverArgs) (api.ResolverResult, error) {
		return api.ResolverResult{Path: filepath.Join(args.ResolveDir, args.Path), Namespace: "css", External: false}, nil
			})

		plugin.AddLoader(api.LoaderOptions{Filter: ".css$", Namespace: "css"},
			func(args api.LoaderArgs) (api.LoaderResult, error) {
				data, err := ioutil.ReadFile(args.Path)

				if err != nil {
					log.Fatalf("\nCould not load %s. Error: %s\n", args.Path, err.Error())
				}
				var content string

				if modules {
					module, err := cssmodule.Process(strings.NewReader(string(data)), args.Path)

					if err != nil {
						return api.LoaderResult{}, err
					}

					content = module.GetExportsString()

					if extract {
						extracted += module.CSS
					} else {
						content = WrapCSSForJSInjection(module.CSS, args.Path, content)
					}
				} else {
					if extract {
						extracted += string(data)
					} else {
						content = WrapCSSForJSInjection(string(data), args.Path, "")
					}
				}

				return api.LoaderResult{Contents: &content, Loader: api.LoaderNone}, nil
			})
	}
}

// GetExtractedCSS - returns extracted stylesheet
func GetExtractedCSS() string {
	return extracted
}

// WrapCSSForJSInjection - wraps a CSS string in JavaScript that injects it into the DOM
func WrapCSSForJSInjection(css string, path string, exports string) string {
	return fmt.Sprintf(`
const style = document.createElement('style')
style.setAttribute('data-path', "%s")
style.innerHTML = `+"`%s`"+`
document.head.appendChild(style)
%s`, path, css, exports)
}
