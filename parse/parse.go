package parse

import (
	"strings"

	"github.com/monopolly/email"
	"github.com/monopolly/email/tools"
	"github.com/monopolly/email/utf"
)

const maxLen = 4000 //telegram

// Parse не меняет переданные options: работает с их копией,
// поэтому один *Options можно использовать из разных горутин.
func Parse(r *email.Email, opts *Options) (err error) {

	options := new(Options)
	if opts != nil {
		*options = *opts
	}

	// любая кодировка -> utf8
	r.Subj = utf.Fix(r.Subj)
	r.Text = tools.NormalizeNewlines(utf.Fix(r.Text))
	r.HTML = utf.Fix(r.HTML)
	r.From.Name = utf.Fix(r.From.Name)
	r.To.Name = utf.Fix(r.To.Name)

	//telegram md escape
	if options.EscapeAddress {
		r.From.Address = tools.TelegramMarkdownEscape(r.From.Address)
		r.To.Address = tools.TelegramMarkdownEscape(r.To.Address)
		r.Subj = tools.TelegramMarkdownEscape(r.Subj)
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
		if tools.IsHTML(source) {
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
		r.HTML = tools.FixHTMLCharset(r.HTML)
		options.Save("htmlfix.html", r.HTML)
	}

	//все картинки
	r.Images = append(r.Images, images...)
	if options.Images && len(r.HTML) > 0 {
		// не качаем повторно то, что уже скачано из ссылок в тексте
		has := map[string]bool{}
		for _, x := range r.Images {
			if x.Link != "" {
				has[x.Link] = true
			}
		}
		var links []string
		for _, x := range tools.ImageLinks(r.HTML) {
			x = strings.TrimSpace(x)
			if has[x] || !strings.HasPrefix(x, "http") {
				continue
			}
			links = append(links, x)
		}
		for _, x := range downloadImages(links, options) {
			if x.Body == nil {
				continue
			}
			addImage(r, x)
		}
	}

	//сокращать только если есть оригинал
	if options.Short && r.HTML != "" {
		if options.MaxWords == 0 {
			options.MaxWords = 100
		}
		r.Markdown = tools.ShortText(r.Markdown, options.MaxWords)
		options.Save("short.md", r.Markdown)
	}

	r.Markdown = tools.TelegramMarkdownEscape(r.Markdown)
	options.Log("total images parse")
	return
}

// addImage добавляет картинку, если картинки с таким ID ещё нет
func addImage(r *email.Email, f *email.File) {
	for _, x := range r.Images {
		if x.ID == f.ID {
			return
		}
	}
	r.Images = append(r.Images, f)
}
