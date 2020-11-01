package livereload

import (
	"fmt"
	"io/ioutil"
	"log"
	"strings"

	"github.com/evanw/esbuild/pkg/api"
	"github.com/markbates/pkger"
)

// Inject - Inserts the livereload JS into an HTML page
func Inject(html string) []byte {
	fd, err := pkger.Open("/internal/livereload/livereload.ts")

	if err != nil {
		log.Fatalf("error %v", err)
	}

	data, err := ioutil.ReadAll(fd)

	if err != nil {
		fmt.Println("Error injecting live reload script into HTML", err.Error())
		return []byte(html)
	}

	result := api.Transform(string(data), api.TransformOptions{
		Loader:            api.LoaderTS,
		MinifyWhitespace:  true,
		MinifyIdentifiers: true,
		MinifySyntax:      true,
		Define:           map[string]string{"globalThis.__LIVE_RELOAD_PATH__": fmt.Sprintf("'%s'", LiveReloadPath)},
	})

	injectedHTML := strings.Replace(html, "</body>", fmt.Sprintf("<script>%s</script></body>", result.Code), 1)

	return []byte(injectedHTML)
}
