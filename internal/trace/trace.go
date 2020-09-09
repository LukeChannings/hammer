package trace

import (
	"fmt"
	"io/ioutil"
	"path"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v2"
)

var resolvedNames = make([]string, 0)

// Trace - computes a dependency graph starting with entrypoint.
func Trace(entrypoint string, flat bool, listOrphans bool) {
	node := Node{path: entrypoint}
	getDependencies(&node)
	resolvedNames = append(resolvedNames, entrypoint)

	sort.Strings(resolvedNames)

	basePath := path.Dir(entrypoint)

	if flat {
		for _, name := range resolvedNames {
			fmt.Println(strings.Replace(name, basePath, "", 1))
		}
	} else if listOrphans {

		orphans := make([]string, 0)

		matches, _ := doublestar.Glob(fmt.Sprintf("%s/**/*.{tsx,ts,jsx,js,css,scss}", basePath))

		for _, m := range matches {
			i := sort.SearchStrings(resolvedNames, m)
			if m != resolvedNames[i] {
				orphans = append(orphans, m)
			}
		}

		for _, o := range orphans {
			fmt.Println(o)
		}
	} else {
		node.Print(0, basePath)
	}
}

func getDependencies(node *Node) {
	data, err := ioutil.ReadFile(node.path)

	if err != nil {
		node.exists = false
	} else {
		node.exists = true
	}

	reader := strings.NewReader(string(data))

	switch path.Ext(node.path) {
	case ".html":
		traceHTMLDependencies(reader, node)
	case ".js", ".jsx", ".ts", ".tsx":
		traceJavaScriptDependencies(reader, node)
	case ".css", ".scss":
		traceCSSDependencies(reader, node)
	}
}
