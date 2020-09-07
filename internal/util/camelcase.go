package util

import "unicode"

// ToCamelCase - converts a string to camel case
func ToCamelCase(input string) string {
	var output []rune
	uppercaseNextChar := false

	for i, char := range input {
		if (char >= 65 && char <= 90) || (char >= 97 && char <= 122) || (char >= 48 && char <= 57) {
			if uppercaseNextChar {
				output = append(output, unicode.ToUpper(char))
				uppercaseNextChar = false
			} else if i == 0 {
				output = append(output, unicode.ToLower(char))
			} else {
				output = append(output, char)
			}
		} else {
			uppercaseNextChar = true
		}
	}

	return string(output)
}
