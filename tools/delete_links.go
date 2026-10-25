package tools

import "strings"

// delete links from text
func DeleteLinks(s string, links []string) (res string) {
	for _, x := range links {
		s = strings.ReplaceAll(s, x, "")
	}
	return s
}

// delete links from text
func DeleteLink(s string, link string) (res string) {
	return strings.ReplaceAll(s, link, "")
}
