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
					if tt == js.StringToken {
						from = strings.Trim(string(text), "\"'")
						break
					}
				}

				childNode := Node{parentNode: node}
				if isURL(from) {
					childNode.path = from
				} else {
					childNode.path = path.Join(path.Dir(node.path), from)
				}

				childNode.Resolve()
				node.dependencies = append(node.dependencies, &childNode)
			}
		}
	}
}
