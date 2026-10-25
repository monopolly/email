package tools

import (
	"strings"
)

func IsHTML(s string) (res bool) {

	switch {
	case
		strings.Contains(s, "<div"),
		strings.Contains(s, "<h1"),
		strings.Contains(s, "<h2"),
		strings.Contains(s, "<h3"),
		strings.Contains(s, "<h4"),
		strings.Contains(s, "<h5"),
		strings.Contains(s, "<h6"),
		strings.Contains(s, "<img"),
		strings.Contains(s, "<b"),
		strings.Contains(s, "<br"),
		strings.Contains(s, "<hr"),
		strings.Contains(s, "<center"),
		strings.Contains(s, "<table"),
		strings.Contains(s, "<code"),
		strings.Contains(s, "<body"),
		strings.Contains(s, "<html"):
		return true

	}

	return
}
