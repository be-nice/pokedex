package utils

import "strings"

func CleanInput(s string) []string {
	return strings.Fields(strings.ToLower(s))
}
