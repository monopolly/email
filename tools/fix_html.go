package tools

import (
	"regexp"
	"strings"
)

var (
	// <meta charset="..."> и <meta http-equiv="Content-Type" content="text/html; charset=...">
	charsetDecl = regexp.MustCompile(`(?i)(<meta[^>]*?charset\s*=\s*["']?\s*)[a-z0-9_:.\-]+`)
	headOpen    = regexp.MustCompile(`(?i)<head(\s[^>]*)?>`)
	htmlOpen    = regexp.MustCompile(`(?i)<html(\s[^>]*)?>`)
)

// FixHTMLCharset помечает html как utf-8: сам текст к этому моменту уже в utf-8 (ToUTF8)
func FixHTMLCharset(html string) (res string) {

	res = html
	res = strings.TrimSpace(res)
	if len(res) < 20 {
		return
	}

	// объявленная кодировка (windows-1251 и т.п.) больше не соответствует тексту
	if charsetDecl.MatchString(res) {
		return charsetDecl.ReplaceAllString(res, "${1}utf-8")
	}

	const meta = `<meta charset="utf-8">`

	if loc := headOpen.FindStringIndex(res); loc != nil {
		return res[:loc[1]] + meta + res[loc[1]:]
	}

	if loc := htmlOpen.FindStringIndex(res); loc != nil {
		return res[:loc[1]] + "<head>" + meta + "</head>" + res[loc[1]:]
	}

	return `<!DOCTYPE html><html><head>` + meta + `</head><body>` + res + "</body></html>"
}
