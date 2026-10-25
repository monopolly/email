package parse

import (
	"strings"

	"github.com/monopolly/email"
	"github.com/monopolly/email/tools"
)

func parseText(source string, options *Options) (marked string, imagelist []*email.File) {

	marked = source
	options.Log("parse text")

	//remove html tags
	marked = tools.DeleteHTMLTags(marked)

	//get links
	links := tools.Links(marked)
	if links == nil && len(marked) < options.MaxTextLen {
		// короткий текст без ссылок не переписываем, только убираем лишние пустые строки
		marked = strings.TrimSpace(tools.RemoveRepeatedEmptyLines(marked, 1))
		return
	}

	options.Log("links found", len(links))

	//delete relative links
	marked = tools.DeleteRelativeLinks(marked)
	if options.Debug {
		options.Save("removerelative.md", marked)
	}

	//clean
	options.Log("clean")
	marked = tools.CleanText(marked)
	if options.Debug {
		options.Save("clean_text.md", marked)
	}

	// remove repeat lines
	marked = tools.RemoveRepeatedEmptyLines(marked, 1)
	if options.Debug {
		options.Save("remove_repeated_lines.md", marked)
	}

	images := make(map[string]*email.File)

	//delete image links
	//if options.html {
	if options.Images {
		options.Log("download all images")
		images = downloadImages(links, options)
		options.Log("images count", len(images))
		for link, img := range images {
			marked = tools.DeleteLink(marked, link)
			imagelist = append(imagelist, img)
		}
		if options.Debug {
			options.Save("removeImageLinks.md", marked)
		}
	}

	// fix md
	options.Log("fix links")
	marked = tools.FixMarkdownLinks(marked)
	if options.Debug {
		options.Save("fixMDlinks.md", marked)
	}

	//weblinks only
	options.Log("get weblinks")
	options.Log("links found", len(links))
	var weblinks []string
	for _, x := range links {
		options.Log("link", x)
		if images[x] != nil {
			options.Log("image link")
			continue
		}
		//если ссылка и так короткая
		if len(x) > 40 {
			options.Log("weblink")
			weblinks = append(weblinks, x)
		}

	}

	options.Log("weblinks count", len(weblinks))

	//short links

	options.Log("make short web links")
	marked = shortTextLinks(marked, weblinks, options)
	if options.Debug {
		options.Save("shortTextLinks.md", marked)
	}

	options.Log(marked)
	marked = strings.TrimSpace(marked)
	for {
		if !strings.Contains(marked, "\n\n\n") {
			break
		}
		marked = strings.ReplaceAll(marked, "\n\n\n", "\n\n")

	}
	if options.Debug {
		options.Save("trimspace.md", marked)
	}

	marked = tools.Retext(marked)
	if options.Debug {
		options.Save("retext.md", marked)
	}

	options.Log("total all images final", len(imagelist))
	return
}
