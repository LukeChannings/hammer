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
						childNode := Node{ParentNode: node}
						importValue := strings.Trim(string(t.Data), "\"")
						if isURL(importValue) {
							childNode.Path = importValue
						} else {
							childNode.Path = path.Join(path.Dir(node.Path), importValue)
						}
						childNode.Resolve()
						node.Dependencies = append(node.Dependencies, &childNode)
					}
				}
			}
		}
	}
}
