package parse

import (
	"net/url"
	"sort"
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

	// сначала самые длинные: иначе ссылка-префикс испортит более длинную
	longs := make([]string, 0, len(shorts))
	for long := range shorts {
		longs = append(longs, long)
	}
	sort.Slice(longs, func(i, j int) bool { return len(longs[i]) > len(longs[j]) })

	for _, long := range longs {
		_, err := url.ParseRequestURI(long)
		if err != nil {
			options.Log("link error", err)
			continue
		}

		raw := shorts[long]
		if !strings.Contains(raw, "://") {
			raw = "https://" + strings.TrimPrefix(raw, "//")
		}
		shortURL, err := url.Parse(raw)
		if err != nil || shortURL.Host == "" {
			options.Log("short link error", shorts[long], err)
			continue
		}

		var short string
		switch options.ShortLinksMustScheme {
		case true:
			shortURL.Scheme = "https"
			short = shortURL.String()
		default:
			// без схемы url.String() даёт "//host/path"
			shortURL.Scheme = ""
			short = strings.TrimPrefix(shortURL.String(), "//")
		}
		options.Log("replace", long, ">>>", short)
		res = strings.ReplaceAll(res, long, short)
	}

	return
}
