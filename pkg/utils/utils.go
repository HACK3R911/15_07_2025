package utils

import "strings"

func ParseExtensions(extStr string) map[string]bool {
	result := make(map[string]bool)
	for _, ext := range strings.Split(extStr, ",") {
		result[strings.ToLower(strings.TrimSpace(ext))] = true
	}
	return result
}
