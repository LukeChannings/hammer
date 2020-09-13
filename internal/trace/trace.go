package trace

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"sort"
	"strings"

	"github.com/bmatcuk/doublestar/v2"
)

var resolvedNames = make([]string, 0)

// Trace - computes a dependency graph starting with entrypoint.
func Trace(entrypoint string, flat bool, listOrphans bool, outputJSON bool) {
	node := Node{Path: entrypoint}
	getDependencies(&node)
	resolvedNames = append(resolvedNames, entrypoint)

	sort.Strings(resolvedNames)

	basePath := path.Dir(entrypoint)

	if flat {
		for i, name := range resolvedNames {
			resolvedNames[i] = strings.Replace(name, basePath, "", 1)
		}
		if outputJSON {
			output, err := json.MarshalIndent(resolvedNames, "", "  ")
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print(string(output))
		} else {
			for _, name := range resolvedNames {
				fmt.Println(name)
			}
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

		if outputJSON {
			output, err := json.MarshalIndent(orphans, "", "  ")
			if err != nil {
				log.Fatal(err)
			}

			fmt.Print(string(output))
		} else {
			for _, o := range orphans {
				fmt.Println(o)
			}
		}
	} else {
		if outputJSON {
			output, err := json.Marshal(node)
			if err != nil {
				log.Fatalln(err)
			}

			fmt.Println(string(output))
		} else {
			node.Print(0, basePath)
		}
	}
}

func getDependencies(node *Node) {
	data, err := ioutil.ReadFile(node.Path)

	if err != nil {
		node.Exists = false
	} else {
		node.Exists = true
	}

	reader := strings.NewReader(string(data))

	switch path.Ext(node.Path) {
	case ".html":
		traceHTMLDependencies(reader, node)
	case ".js", ".jsx", ".ts", ".tsx":
		traceJavaScriptDependencies(reader, node)
	case ".css", ".scss":
		traceCSSDependencies(reader, node)
	}
}
