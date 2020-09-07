package cssmodule

import (
	"crypto/sha256"
	"fmt"
	"io"
	"path"
	"strings"

	"github.com/lukechannings/hammer/internal/util"
	"github.com/tdewolff/parse/css"
)

// Result contains the processed CSS output with its exports
type Result struct {
	CSS     string
	Exports map[string]string
}

// GetExportsString - a string of ESM exports
func (result Result) GetExportsString() string {
	var out string
	for k, v := range result.Exports {
		out += fmt.Sprintf("export const %s = \"%v\";\n", k, v)
	}
	return out
}

var indentToken = "  "

// Process converts plain CSS into a CSS Module result
func Process(r io.Reader, fullPath string) (Result, error) {
	p := css.NewParser(r, false)
	var out string
	var indent int

	exports := make(map[string]string)

	for {
		gt, _, data := p.Next()
		if gt == css.ErrorGrammar {
			if p.Err().Error() == "EOF" {
				break
			} else {
				return Result{}, p.Err()
			}
		}

		fmt.Println(gt)

		switch gt {
		case css.BeginRulesetGrammar, css.BeginAtRuleGrammar:
			{
				ruleset := strings.Repeat(indentToken, indent) + string(data)
				indent++

				ts := p.Values()
				for i, t := range ts {
					if t.TokenType == css.IdentToken && (i != 0 && ts[i-1].TokenType == css.DelimToken) && !strings.Contains(ruleset, ":global") && !strings.Contains(ruleset, "global(") {
						className, uniqueClassName := getCanonicalClassNames(string(t.Data), fullPath)
						if _, ok := exports[string(t.Data)]; !ok {
							exports[className] = uniqueClassName
						}
						ruleset += exports[className]
					} else {
						ruleset += string(t.Data)
					}
				}

				ruleset += " {\n"
				out += ruleset
			}
		case css.EndRulesetGrammar, css.EndAtRuleGrammar:
			{
				indent--
				out += strings.Repeat(indentToken, indent) + "}\n\n"
			}
		case css.DeclarationGrammar, css.CustomPropertyGrammar:
			{
				out += strings.Repeat(indentToken, indent) + string(data) + ": "
				for _, t := range p.Values() {
					out += string(t.Data)
				}
				out += ";\n"
			}
		}
	}

	return Result{CSS: out, Exports: exports}, nil
}

func getCanonicalClassNames(className string, fullPath string) (string, string) {
	normalClassName := util.ToCamelCase(className)
	baseName := strings.Replace(path.Base(fullPath), path.Ext(fullPath), "", 1)
	return normalClassName, fmt.Sprintf("%s_%s-%.3x", baseName, normalClassName, sha256.Sum256([]byte(fullPath+className)))
}
