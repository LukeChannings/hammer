package trace

import (
	"io"
	"path"
	"strings"

	"github.com/tdewolff/parse/js"
)

func traceJavaScriptDependencies(r io.Reader, node *Node) {
	l := js.NewLexer(r)

main:
	for {
		tt, text := l.Next()
		switch tt {
		case js.ErrorToken:
			break main
		case js.IdentifierToken:
			if string(text) == "import" {
				var from string
				for {
					tt, text := l.Next()
					if tt == js.ErrorToken {
						break main
					}
					if string(text) == "type" {
						break
					}
					if tt == js.StringToken {
						from = strings.Trim(string(text), "\"'")
						break
					}
				}

				childNode := Node{ParentNode: node}
				if isURL(from) {
					childNode.Path = from
				} else {
					childNode.Path = path.Join(path.Dir(node.Path), from)
				}

				childNode.Resolve()
				node.Dependencies = append(node.Dependencies, &childNode)
			}
		}
	}
}
