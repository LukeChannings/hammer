package util

import (
	"testing"
)

func TestToCamelCase(t *testing.T) {
	var fixtures = map[string]string{
		"not-snake-case":               "notSnakeCase",
		"how about this?":              "howAboutThis",
		"we11_then":                    "we11Then",
		"How_About_These_Underscores?": "howAboutTheseUnderscores",
		"alreadyCamelCase":             "alreadyCamelCase",
	}

	for a, expected := range fixtures {
		if result := ToCamelCase(a); result != expected {
			t.Fatalf("Expected ToCamelCase(%s) to be %s, got '%s'", a, expected, result)
		}
	}
}
