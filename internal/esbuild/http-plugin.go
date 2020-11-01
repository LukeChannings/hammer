package esbuild

import (
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
)

// HTTPPlugin - handles external imports in esbuild
func HTTPPlugin(bundle bool) api.Plugin {
	return api.Plugin{
		Name: "http-plugin",
		Setup: func(b api.PluginBuild) {
			b.OnResolve(api.OnResolveOptions{Namespace: "http", Filter: "^(https?:)?//"}, func(args api.OnResolveArgs) (api.OnResolveResult, error) {
				return api.OnResolveResult{Path: args.Path, Namespace: "http", External: !bundle}, nil
			})
			b.OnLoad(api.OnLoadOptions{Namespace: "http", Filter: "^(https?:)?//"}, func(args api.OnLoadArgs) (api.OnLoadResult, error) {
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
				return api.OnLoadResult{Contents: &content, Loader: api.LoaderNone}, nil
			})
		},
	}
}
