package trace

import (
	"io"
	"path"
	"strings"

	"github.com/tdewolff/parse/html"
)

func traceHTMLDependencies(r io.Reader, node *Node) {
	l := html.NewLexer(r)

main:
	for {
		tt, _ := l.Next()
		switch tt {
		case html.ErrorToken:
			break main
		case html.StartTagToken:
			if string(l.Text()) == "script" {
				for {
					ttAttr, _ := l.Next()
					if ttAttr != html.AttributeToken {
						break
					}

					if string(l.Text()) == "src" {
						path := path.Join(path.Dir(node.path), strings.Trim(string(l.AttrVal()), "\""))
						childNode := Node{path: path, parentNode: node}
						childNode.Resolve()
						node.dependencies = append(node.dependencies, &childNode)
					}
				}
			} else {
				continue
			}
		}
	}
}
