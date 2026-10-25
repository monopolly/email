package tools

import (
	"strings"
	"unicode"
)

// any text from start until \t\n\space
func Tag(raw string, key string, endtag ...string) (value string) {
	i := strings.Index(raw, key)
	if i == -1 {
		return
	}
	for _, x := range raw[i+len(key):] {
		switch {
		case unicode.IsSpace(x):
			if endtag != nil {
				value = value + string(x)
			} else {
				return
			}
		default:
			if endtag != nil && string(x) == endtag[0] {
				return
			}
			value = value + string(x)
		}
	}

	return
}
