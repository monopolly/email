package tools

import "unicode/utf8"

// ---------------------------
func HasRepeats(s string) (res bool) {
	if s == "" {
		return false
	}

	first, _ := utf8.DecodeRuneInString(s)
	for _, x := range s {
		if x != first {
			return
		}
	}

	return true
}
