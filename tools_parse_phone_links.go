package email

import (
	"strings"
)

// tel:+8429829842t5
func getPhoneLinks(v string, opt *Options) (links []string) {
	un := map[string]bool{}
	for _, link := range xurlsTel.FindAllString(v, -1) {
		link = strings.TrimSpace(link)
		if un[link] {
			continue
		}
		un[link] = true
		links = append(links, link)
	}
	return
}
