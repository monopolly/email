package email

import (
	"net/url"
	"strings"
)

// просто заменяет длинные ссылки на короткие, без md
func shortTextLinks(s string, links []string, options *Options) (res string) {
	res = s
	var shorts map[string]string

	//try short
	if options.ShortLinks == nil {
		options.Log("no short links handler")
		return
	}

	//make short
	shorts = options.ShortLinks(links...)
	if shorts == nil {
		options.Log("short links connection server error")
		return
	}

	options.Log("short links count", len(shorts))

	for long, short := range shorts {
		_, err := url.ParseRequestURI(long)
		if err != nil {
			options.Log("link error", err)
			continue
		}

		// if ul.Port() != "" {
		// 	shorts[long] = long
		// 	continue
		// }

		shortURL, _ := url.Parse(short)
		switch options.ShortLinksMustScheme {
		case true:
			shortURL.Scheme = "https"
		default:
			shortURL.Scheme = ""
		}
		options.Log("replace", long, ">>>", shortURL)
		res = strings.ReplaceAll(res, long, shortURL.String())
	}

	return
}
