package email

import "strings"

// delete links from text
func deleteLinksFromText(s string, links []string) (res string) {
	for _, x := range links {
		s = strings.ReplaceAll(s, x, "")
	}
	return s
}

// delete links from text
func deleteLinkFromText(s string, link string) (res string) {
	return strings.ReplaceAll(s, link, "")
}
