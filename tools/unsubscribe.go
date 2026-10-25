package email

import (
	"net/url"
	"strings"
)

func ParseUnsubscribeLink(raw string) (link *url.URL) {
	raw = strings.ReplaceAll(raw, "<", "")
	raw = strings.ReplaceAll(raw, ">", "")
	link, _ = url.Parse(raw)
	return
}
