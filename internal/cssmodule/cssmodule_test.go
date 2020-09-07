package cssmodule

import (
	"strings"
	"testing"
)

func TestCanonicalClassName(t *testing.T) {
	exportedName, replacedName := getCanonicalClassNames("test", "foo/bar/test.css")

	expectedExportedName := "test"
	expectedReplacedName := "test_test-969528"

	if exportedName != expectedExportedName {
		t.Fatalf("Expected %s, got %s\n", expectedExportedName, exportedName)
	}

	if replacedName != expectedReplacedName {
		t.Fatalf("Expected %s, got %s\n", expectedReplacedName, replacedName)
	}
}

func TestCanonicalClassNameCamelCase(t *testing.T) {
	exportedName, replacedName := getCanonicalClassNames("this-is-snake-case", "foo/bar/test.css")

	expectedExportedName := "thisIsSnakeCase"
	expectedReplacedName := "test_thisIsSnakeCase-4d7e5a"

	if exportedName != expectedExportedName {
		t.Fatalf("Expected %s, got %s\n", expectedExportedName, exportedName)
	}

	if replacedName != expectedReplacedName {
		t.Fatalf("Expected %s, got %s\n", expectedReplacedName, replacedName)
	}
}

func TestPrettyPrint(t *testing.T) {
	css := ":root { --test: red; background-color: var(--red) } .foo { color: var(--test); }"
	result, _ := Process(strings.NewReader(css), "test")
	expected := `:root {
  --test: red;
  background-color: var(--red);
}

.test_foo-6e03ac {
  color: var(--test);
}`
	if expected != result.CSS {
		t.Fatalf("Expected %s, got %s.", expected, result.CSS)
	}
}
