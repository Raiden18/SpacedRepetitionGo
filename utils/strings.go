package utils

import "strings"

func SubstringAfter(s, sep string) string {
	idx := strings.Index(s, sep)
	if idx == -1 {
		return ""
	}
	return s[idx+len(sep):]
}

func SubstringBefore(s, sep string) string {
	idx := strings.Index(s, sep)
	if idx == -1 {
		return s
	}
	return s[:idx]
}
