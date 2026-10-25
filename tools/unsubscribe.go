package tools

import (
	"net/url"
	"regexp"
	"strings"
)

var angleLinks = regexp.MustCompile(`<([^>]+)>`)

// List-Unsubscribe: <https://...>, <mailto:...>
// возвращает http(s)-ссылку, если она есть, иначе mailto; nil, если ссылок нет
func ParseUnsubscribeLink(raw string) (link *url.URL) {
	var list []string
	for _, m := range angleLinks.FindAllStringSubmatch(raw, -1) {
		list = append(list, m[1])
	}
	if list == nil {
		list = strings.Split(raw, ",")
	}

	for _, x := range list {
		u, err := url.Parse(strings.TrimSpace(x))
		if err != nil || u.Scheme == "" {
			continue
		}
		switch {
		case u.Scheme == "http" || u.Scheme == "https":
			return u
		case link == nil:
			link = u
		}
	}
	return
}
