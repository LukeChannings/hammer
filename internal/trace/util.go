package trace

import (
	"os"
	"strings"
)

func fileExists(path string) bool {
	stat, err := os.Stat(path)

	if err == nil && !stat.IsDir() {
		return true
	}

	if os.IsNotExist(err) {
		return false
	}

	return false
}

func isURL(path string) bool {
	return strings.HasPrefix(path, "http:") || strings.HasPrefix(path, "https:") || strings.HasPrefix(path, "//")
}
