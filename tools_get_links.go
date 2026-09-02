package email

import (
	"net/url"
	"strings"
)

// first link
func Link(s string) (link string) {
	for _, x := range strings.Split(s, "http") {
		if !strings.Contains(x, "://") {
			return
		}
		x = "http" + strings.TrimSpace(x)

		x = strings.TrimSuffix(x, ".")
		x = strings.TrimSuffix(x, "!")
		x = strings.TrimSuffix(x, "?")
		x = strings.TrimSuffix(x, "=")
		x = strings.TrimSuffix(x, "+")
		x = strings.TrimSuffix(x, "-")
		x = strings.TrimSuffix(x, "&")
		x = strings.TrimSuffix(x, "@")
		x = strings.TrimSuffix(x, "%")

		x = strings.ReplaceAll(x, "<", " ")
		x = strings.ReplaceAll(x, ">", " ")

		p := strings.Fields(x)
		if len(p) == 0 {
			continue
		}
		_, err := url.ParseRequestURI(p[0])
		if err != nil {
			continue
		}
		return p[0]

	}
	return
}

func Links(source string) (res []string) {

	//var rawlinks []string //some https://a.com?url=https://b.com?redirect_to_site
	for _, x := range strings.Split(source, "\n") {
		if !strings.Contains(x, "://") {
			continue
		}
		x = strings.TrimSpace(x)
		list := xurlsStrict.FindAllString(x, -1)
		if list != nil {
			res = append(res, list...)
			continue
		}
	}
	return
}

func LinksMap(s string) (maps map[string]bool) {
	maps = make(map[string]bool)
	for _, x := range Links(s) {
		maps[x] = true
	}
	return
}

// only http* links
// no apple.com
// func getAllLinks(v string) (links []string) {
// 	un := map[string]bool{}
// 	xurlsStrict, _ := xurls.StrictMatchingScheme("http")
// 	for _, link := range xurlsStrict.FindAllString(v, -1) {
// 		link = strings.TrimSpace(link)
// 		if un[link] {
// 			continue
// 		}
// 		links = append(links, link)
// 		un[link] = true
// 	}
// 	return
// }

// func getAllLinksManualInMarkdown(v string) (links []string) {
// 	un := getAllLinksMap(v)

// 	for x := range un {
// 		v = strings.ReplaceAll(v, x, "")
// 	}

// 	for _, x := range strings.Fields(v) {
// 		x = strings.ReplaceAll(x, "[](", " ")
// 		x = strings.ReplaceAll(x, "](", " ")
// 		x = strings.ReplaceAll(x, "![]", " ")
// 		x = strings.ReplaceAll(x, "![", " ")
// 		x = strings.ReplaceAll(x, "[", " ")
// 		x = strings.ReplaceAll(x, "]", " ")
// 		x = strings.ReplaceAll(x, "(", " ")
// 		x = strings.ReplaceAll(x, ")", " ")
// 		x = strings.ReplaceAll(x, "!", " ")
// 		x = strings.ReplaceAll(x, ">", " ")
// 		x = strings.ReplaceAll(x, "<", " ")

// 		x = strings.TrimPrefix(x, "-")
// 		x = strings.TrimPrefix(x, ",")
// 		x = strings.TrimPrefix(x, "=")
// 		x = strings.TrimPrefix(x, ",")

// 		x = strings.TrimSuffix(x, ".")
// 		x = strings.TrimSuffix(x, "/")
// 		x = strings.TrimSuffix(x, ",")

// 		if strings.HasPrefix(x, "http") && !un[x] {
// 			un[x] = true
// 			links = append(links, x)
// 		}
// 	}

// 	return
// }

// func getAllLinksManualInMarkdownMap(v string) (links map[string]bool) {
// 	links = make(map[string]bool)
// 	for _, x := range getAllLinksManualInMarkdown(v) {
// 		links[x] = true
// 	}
// 	return
// }

// no apple.com
// func getAllLinksMap(v string) (links map[string]bool) {
// 	links = make(map[string]bool)
// 	xurlsStrict, _ := xurls.StrictMatchingScheme("http")
// 	for _, link := range xurlsStrict.FindAllString(v, -1) {
// 		link = strings.TrimSpace(link)
// 		if links[link] {
// 			continue
// 		}
// 		links[link] = true
// 	}
// 	return
// }
