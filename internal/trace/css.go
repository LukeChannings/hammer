package trace

import (
	"io"
	"path"
	"strings"

	"github.com/tdewolff/parse/css"
)

func traceCSSDependencies(r io.Reader, node *Node) {
	l := css.NewParser(r, false)

main:
	for {
		gt, _, text := l.Next()
		switch gt {
		case css.ErrorGrammar:
			break main
		case css.AtRuleGrammar:
			if string(text) == "@import" {
				ts := l.Values()
				for _, t := range ts {
					if t.TokenType == css.StringToken {
						childNode := Node{parentNode: node}
						importValue := strings.Trim(string(t.Data), "\"")
						if isURL(importValue) {
							childNode.path = importValue
						} else {
							childNode.path = path.Join(path.Dir(node.path), importValue)
						}
						childNode.Resolve()
						node.dependencies = append(node.dependencies, &childNode)
					}
				}
			}
		}
	}
}
