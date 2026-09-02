package email

import (
	"strings"
)

// only http* links
// no apple.com
func getEmailLinks(v string, opt *Options) (links []string) {
	un := map[string]bool{}
	for _, link := range xurlsMailto.FindAllString(v, -1) {
		link = strings.TrimSpace(link)
		if un[link] {
			continue
		}
		un[link] = true
		links = append(links, link)
	}
	return
}
