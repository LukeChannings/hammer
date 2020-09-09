package trace

import (
	"fmt"
	"path"
	"strings"
)

// Node - represents a module and its dependencies
type Node struct {
	path         string
	exists       bool
	dependencies []*Node
	parentNode   *Node
}

// Print - prints out the dependency tree for the module
func (n *Node) Print(indent int, basePath string) {
	if n.exists {
		fmt.Printf("%s%s\n", strings.Repeat("  ", indent), strings.Replace(n.path, basePath, "", 1))
		for _, dep := range n.dependencies {
			dep.Print(indent+1, basePath)
		}
	}
}

// Resolve - resolves the path if possible and resolves any dependencies
func (n *Node) Resolve() {
	if fileExists(n.path) {
		n.exists = true
	} else {
		var defaultExt string

		if n.parentNode != nil {
			defaultExt = path.Ext(n.parentNode.path)
		} else {
			defaultExt = ".js"
		}

		possibilities := []string{
			n.path + defaultExt,
			n.path + ".js",
			n.path + ".jsx",
			n.path + ".ts",
			n.path + ".tsx",
			n.path + "/index" + defaultExt,
			n.path + "/index.js",
			n.path + "/index.jsx",
			n.path + "/index.ts",
			n.path + "/index.tsx",
		}

		for _, possibility := range possibilities {
			if fileExists(possibility) {
				n.path = possibility
				n.exists = true
			}
		}
	}

	if n.exists {
		for _, resolved := range resolvedNames {
			if resolved == n.path {
				return
			}
		}
		resolvedNames = append(resolvedNames, n.path)

		getDependencies(n)
	}
}
