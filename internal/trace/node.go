package trace

import (
	"fmt"
	"path"
	"strings"
)

// Node - represents a module and its dependencies
type Node struct {
	Path         string  `json:"path"`
	Exists       bool    `json:"exists"`
	Dependencies []*Node `json:"dependencies,omitempty"`
	ParentNode   *Node   `json:"parentNode,omitempty"`
}

// Print - prints out the dependency tree for the module
func (n *Node) Print(indent int, basePath string) {
	if n.Exists {
		fmt.Printf("%s%s\n", strings.Repeat("  ", indent), strings.Replace(n.Path, basePath, "", 1))
		for _, dep := range n.Dependencies {
			dep.Print(indent+1, basePath)
		}
	}
}

// Resolve - resolves the path if possible and resolves any dependencies
func (n *Node) Resolve() {
	if fileExists(n.Path) {
		n.Exists = true
	} else {
		var defaultExt string

		if n.ParentNode != nil {
			defaultExt = path.Ext(n.ParentNode.Path)
		} else {
			defaultExt = ".js"
		}

		possibilities := []string{
			n.Path + defaultExt,
			n.Path + ".js",
			n.Path + ".jsx",
			n.Path + ".ts",
			n.Path + ".tsx",
			n.Path + "/index" + defaultExt,
			n.Path + "/index.js",
			n.Path + "/index.jsx",
			n.Path + "/index.ts",
			n.Path + "/index.tsx",
		}

		for _, possibility := range possibilities {
			if fileExists(possibility) {
				n.Path = possibility
				n.Exists = true
			}
		}
	}

	if n.Exists {

		for _, resolved := range resolvedNames {
			if resolved == n.Path {
				return
			}
		}

		resolvedNames = append(resolvedNames, n.Path)

		getDependencies(n)
	}
}
