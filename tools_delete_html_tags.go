package email

import "strings"

// clean html tags in text, a, h1 ,etc
func deleteHTMLTags(s string) (res string) {
	res = s

	//links
	res = strings.ReplaceAll(res, "<a>", "")
	res = strings.ReplaceAll(res, "<a", "")
	res = strings.ReplaceAll(res, "</a>", "")
	res = strings.ReplaceAll(res, `href="`, "")
	res = strings.ReplaceAll(res, `target="_blank"`, "")

	//img
	res = strings.ReplaceAll(res, "<img>", "")
	res = strings.ReplaceAll(res, "</img>", "")
	res = strings.ReplaceAll(res, `src="`, "")

	//h1
	res = strings.ReplaceAll(res, "<h1>", "*")
	res = strings.ReplaceAll(res, "</h1>", "*")
	res = strings.ReplaceAll(res, "<h2>", "*")
	res = strings.ReplaceAll(res, "</h2>", "*")
	res = strings.ReplaceAll(res, "<h3>", "*")
	res = strings.ReplaceAll(res, "</h3>", "*")
	res = strings.ReplaceAll(res, "<h4>", "*")
	res = strings.ReplaceAll(res, "</h4>", "*")
	res = strings.ReplaceAll(res, "<h5>", "*")
	res = strings.ReplaceAll(res, "</h5>", "*")
	res = strings.ReplaceAll(res, "<h6>", "*")
	res = strings.ReplaceAll(res, "</h6>", "*")

	//span
	res = strings.ReplaceAll(res, "<body>", "")
	res = strings.ReplaceAll(res, "</body>", "")

	//head, style
	res = string(deleteTag([]byte(res), "<head>", "</head>"))
	res = string(deleteTag([]byte(res), "<style>", "</style>"))

	//span
	res = strings.ReplaceAll(res, "<span>", "")
	res = strings.ReplaceAll(res, "</span>", "")

	//p
	res = strings.ReplaceAll(res, "<p>", "")
	res = strings.ReplaceAll(res, "</p>", "")

	//code
	res = strings.ReplaceAll(res, "<code>", "```")
	res = strings.ReplaceAll(res, "</code>", "```")

	//pre
	res = strings.ReplaceAll(res, "<pre>", "```")
	res = strings.ReplaceAll(res, "</pre>", "```")

	//divs
	res = strings.ReplaceAll(res, "<div>", "")
	res = strings.ReplaceAll(res, "</div>", "")

	//br, hr
	res = strings.ReplaceAll(res, "<br>", "")
	res = strings.ReplaceAll(res, "<hr>", "")

	res = strings.ReplaceAll(res, "<b>", "*")
	res = strings.ReplaceAll(res, "</b>", "*")

	res = strings.ReplaceAll(res, "<blockquote>", "_")
	res = strings.ReplaceAll(res, "</blockquote>", "_")
	res = strings.ReplaceAll(res, "<i>", "_")
	res = strings.ReplaceAll(res, "</i>", "_")

	return
}
