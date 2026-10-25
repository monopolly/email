package tools

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
)

func LinksHTML(htmlsource string) (list []string) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer([]byte(htmlsource)))
	if err != nil {
		return
	}
	doc.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		link, _ := item.Attr("href")
		list = append(list, link)
	})
	return
}

func ImageLinks(htmlsource string) (list []string) {
	doc, err := goquery.NewDocumentFromReader(bytes.NewBuffer([]byte(htmlsource)))
	if err != nil {
		return
	}

	un := map[string]bool{}
	doc.Find("img[src]").Each(func(index int, item *goquery.Selection) {
		link, _ := item.Attr("src")
		if un[link] {
			return
		}
		un[link] = true
		list = append(list, link)
	})

	for _, link := range LinksHTML(htmlsource) {
		if IsImageLink(link) {
			if un[link] {
				continue
			}
			un[link] = true
			list = append(list, link)
		}
	}

	return
}
