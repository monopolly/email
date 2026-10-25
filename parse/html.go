package email

import (
	"bytes"

	"github.com/PuerkitoBio/goquery"
	"github.com/k3a/html2text"

	//md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/monopolly/markdown"
)

func parseHTML(html string, options *Options) (res string, err error) {

	return markdown.Parse(html)
	// //empty
	// html = strings.TrimSpace(html)
	// if emptyText(html) {
	// 	return
	// }

	// opt := new(md.Options)
	// //opt.StrongDelimiter = "*"
	// // opt.CodeBlockStyle = "fenced"
	// // opt.HorizontalRule = ""
	// // opt.LinkReferenceStyle = "collapsed"
	// // opt.LinkStyle = "referenced"

	// converter := md.NewConverter("", true, opt)

	// converter.AddRules(
	// 	md.Rule{
	// 		Filter: []string{"a"},
	// 		Replacement: func(c string, selection *goquery.Selection, opt *md.Options) *string {
	// 			title := strings.TrimSpace(selection.Text())
	// 			if title == "" || strings.Contains(title, "<img>") {
	// 				return md.String("")
	// 			}

	// 			title = strings.TrimSpace(title)

	// 			//link
	// 			link, has := selection.Attr("href")
	// 			if !has {
	// 				return md.String("")
	// 			}

	// 			//log.Println("link title", title)

	// 			//no web
	// 			if !strings.HasPrefix(link, "http") {
	// 				return md.String("")
	// 			}

	// 			// if options.NoImages {
	// 			// 	if isImageLink(link, options) {
	// 			// 		return md.String("")
	// 			// 	}
	// 			// }

	// 			//absolute
	// 			u, err := url.ParseRequestURI(link)
	// 			if err != nil {
	// 				return md.String("")
	// 			}

	// 			if title == "" {
	// 				title = u.Hostname()
	// 			}

	// 			if options.MarkdownLinkFormater != nil {
	// 				title, link = options.MarkdownLinkFormater(title, link)
	// 			}

	// 			p := fmt.Sprintf("[%s]( %s )", title, link)
	// 			return md.String(p)
	// 		},
	// 	},

	// 	md.Rule{
	// 		Filter: []string{"h1", "h2", "h3", "h4", "h5", "h6", "b"},
	// 		Replacement: func(content string, selec *goquery.Selection, opt *md.Options) *string {
	// 			content = strings.TrimSpace(content)
	// 			return md.String(fmt.Sprintf("*%s*", content))
	// 		},
	// 	},
	// )

	// return converter.ConvertString(html)

}

func parseSimple(r *Email, options *Options) {
	text := html2text.HTML2TextWithOptions(string(r.HTML), html2text.WithListSupport(), html2text.WithUnixLineBreaks())
	r.Markdown = cleanText(text)
}

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
		list = append(list, link)
	})

	for _, link := range LinksHTML(htmlsource) {
		if isImageLink(link) {
			if un[link] {
				continue
			}
			list = append(list, link)
		}
	}

	return
}

// func Markdown(htmlsource string) {

// 	return
// }
