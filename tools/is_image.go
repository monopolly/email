package email

import "strings"

func isImageLink(link string) (res bool) {
	if !strings.Contains(link, "http") {
		return
	}
	p := strings.Split(link, ".")
	if len(p) < 3 {
		return
	}
	switch p[len(p)-1] {
	case "jpg", "jpeg", "png", "gif":
		return true
	}
	return
}
