package strings

import "strings"

// TrimSlash trims all "/" and "\\" from both prefix and suffix of the string
func TrimSlash(str string) string {
	for strings.HasPrefix(str, "/") || strings.HasPrefix(str, "\\") {
		str = strings.TrimPrefix(str, "/")
		str = strings.TrimPrefix(str, "\\")
	}

	for strings.HasSuffix(str, "/") || strings.HasSuffix(str, "\\") {
		str = strings.TrimSuffix(str, "/")
		str = strings.TrimSuffix(str, "\\")
	}

	return str
}
