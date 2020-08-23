package esbuild

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

// HTTPPlugin - handles external imports in esbuild
func HTTPPlugin(bundle bool) func(api.Plugin) {
	return func(plugin api.Plugin) {
		plugin.SetName("http-plugin")

		plugin.AddResolver(api.ResolverOptions{Filter: "^(https?:)?//"},
			func(args api.ResolverArgs) (api.ResolverResult, error) {
				return api.ResolverResult{Path: args.Path, Namespace: "http", External: !bundle}, nil
			})

		plugin.AddLoader(api.LoaderOptions{Filter: "^https?://", Namespace: "http"},
			func(args api.LoaderArgs) (api.LoaderResult, error) {
				resp, err := http.Get(args.Path)
				var content string
				if err != nil {
					content = ""
				} else {
					defer resp.Body.Close()
					body, err := ioutil.ReadAll(resp.Body)
					if err != nil {
						content = ""
					} else {
						content = strings.ReplaceAll(string(body), "from '/", "from 'https://cdn.skypack.dev/")
					}
				}
				return api.LoaderResult{Contents: &content, Loader: api.LoaderNone}, nil
			})
	}
}
