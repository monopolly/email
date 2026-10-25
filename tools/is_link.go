package email

import "strings"

func isLink(link string) bool {
	return strings.Contains(link, "http")
}
