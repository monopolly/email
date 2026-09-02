package email

const maxLen = 4000 //telegram

func Parse(r *Email, options *Options) (err error) {

	r.Options = options

	//telegram md escape
	if options.EscapeAddress {
		r.From.Address = TelegramMarkdownEscape(r.From.Address)
		r.To.Address = TelegramMarkdownEscape(r.To.Address)
		r.Subj = TelegramMarkdownEscape(r.Subj)
	}

	//обрезать текст
	if options.MaxTextLen == 0 {
		options.MaxTextLen = maxLen
	}

	var source string
	textmode := len(r.HTML) == 0

	//body
	switch textmode {
	case true:
		//has text
		source = r.Text
		//если в обычном тексте есть html
		if isHTML(source) {
			options.Log("html in text")
			source, _ = parseHTML(source, options)
		}
		options.Log("as text")
		options.Save("origin.txt", source)
	default:
		//html only
		options.Log("as html")
		r.ParserHTML = true
		options.html = true
		options.Save("origin.html", r.HTML)

		source, err = parseHTML(r.HTML, options)
		if err != nil {
			return
		}

		if options.Debug {
			options.Save("html_converted.html", source)
		}
	}

	//парсим текст
	md, images := parseText(source, options)
	r.Markdown = md

	if r.ParserHTML {
		options.Save("result.md", r.Markdown)
	}

	//original fix utf8
	if len(r.HTML) > 0 {
		//fix utf8 encode
		options.Log("html fix")
		r.HTML = htmlfix(r.HTML, options)
		options.Save("htmlfix.html", r.HTML)
	}

	//все картинки
	r.Images = append(r.Images, images...)
	if options.Images && len(r.HTML) > 0 {
		links := ImageLinks(r.HTML)
		if links != nil {
			for _, x := range downloadImages(links, options) {
				if x.Body == nil {
					continue
				}
				r.Images = append(r.Images, x)
			}
		}
	}

	//сокращать только если есть оригинал
	if options.Short && r.HTML != "" {
		if options.MaxWords == 0 {
			options.MaxWords = 100
		}
		r.Markdown = shortText(r.Markdown, options.MaxWords)
		options.Save("short.md", r.Markdown)
	}

	r.Markdown = TelegramMarkdownEscape(r.Markdown)
	options.Log("total images parse")
	return
}
