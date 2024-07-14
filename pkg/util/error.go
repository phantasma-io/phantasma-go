package util

import "strings"

// Checks API response string for errors.
func ErrorDetect(s string) bool {
	if strings.Contains(strings.ToLower(s), "error") {
		return true
	} else {
		return false
	}
}
