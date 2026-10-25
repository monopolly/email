package tools

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
	"unicode"
)

func FixMarkdownLinks(text string) (res string) {

	var lines []string
	for _, i := range strings.Split(text, "\n") {

		//md links
		links := parseMarkdownLinks(i)

		//no links
		if links == nil {
			lines = append(lines, i)
			continue
		}

		//parse line
		for _, x := range links {

			//no link
			if x.link == "" {
				lines = append(lines, strings.TrimSpace(strings.Join(x.title, " ")))
				continue
			}

			//has link only
			if x.title == nil {
				title := parentHost(x.link)
				lines = append(lines, fmt.Sprintf("[%s](%s)", title, x.link))
				continue
			}

			//link and title
			//long title
			lens := len(x.title)

			if lens > 2 {
				if len(x.title[lens-1]) < 5 {
					pre := strings.Join(x.title[:lens-2], " ")
					title := strings.Join(x.title[lens-2:], " ")
					lines = append(lines, fmt.Sprintf("%s [%s](%s)", pre, title, x.link))
				} else {
					pre := strings.Join(x.title[:lens-1], " ")
					title := x.title[lens-1]
					lines = append(lines, fmt.Sprintf("%s [%s](%s)", pre, title, x.link))
				}

				continue
			}
			lines = append(lines, fmt.Sprintf("[%s](%s)", strings.Join(x.title, " "), x.link))
		}
	}

	return strings.Join(lines, "\n")
}

func clearDoubleEmptyLines(text string) (res string) {
	var lines []string
	var space bool
	for _, x := range strings.Split(text, "\n") {
		x = strings.TrimSpace(x)
		if x == "" && space {
			continue
		}
		if x != "" {
			space = false
			lines = append(lines, x)
		}
	}
	return strings.Join(lines, "\n")
}

func parseLine(line string) *mdline {

	line = strings.TrimSpace(line)
	line = cleanTextHard(line)
	if emptyText(line) {
		return nil
	}

	for _, x := range mdreplacers {
		line = strings.ReplaceAll(line, x[0], x[1])
	}
	line = strings.TrimSpace(line)

	if emptyText(line) {
		return nil
	}

	newline := new(mdline)
	newline.body.WriteString(line)
	newline.link = strings.Contains(line, "http")
	newline.list = strings.HasPrefix(line, "-") || strings.HasPrefix(line, "*")
	newline.comma = strings.HasSuffix(line, ",")
	newline.dot = strings.HasSuffix(line, ".")
	newline.twodots = strings.HasSuffix(line, ":")
	newline.letter = unicode.IsLetter(rune(line[0]))

	if unicode.IsDigit(rune(line[0])) || unicode.IsLetter(rune(line[0])) {
		newline.title = strings.ToLower(string(line[0])) != string(line[0])
	}
	newline.lastSmallletter = unicode.IsLetter(rune(line[len(line)-1])) || unicode.IsNumber(rune(line[len(line)-1]))

	_, err := strconv.Atoi(strings.Split(line, ".")[0])
	newline.enum = err == nil

	return newline
}

var mdreplacers = [][]string{
	[]string{"[[", "["},
	[]string{"]]", "]"},
	[]string{"((", "("},
	[]string{"))", ")"},
	[]string{"[ [", "["},
	[]string{"] ]", "]"},
	[]string{"( (", "("},
	[]string{") )", ")"},
	[]string{"[(", "("},
	[]string{"])", "]"},
	[]string{"[ (", "["},
	[]string{"] )", "]"},
	[]string{"([", "["},
	[]string{")]", ")"},
	[]string{"<i>", "_"},
	[]string{"</i>", "_"},
	[]string{")[", ")\n["}, //link,link,link
	[]string{"( [", "["},
	[]string{`( “[`, "["},
	[]string{`)” )`, ")"},
	[]string{"[]", ""},
	[]string{"[!", ""},
	[]string{"]()", "]"},
	[]string{"()", ""},
	// []string{".\n", ". "},
	// []string{"\n- ", "- "},
	// []string{" \n", " "},
	// []string{"\n", " "},
	// []string{") ", ") \n"},
}

type mdline struct {
	link            bool
	list            bool
	enum            bool
	comma           bool
	dot             bool
	letter          bool //start with letter
	lastSmallletter bool //end with letter
	title           bool //first letter
	twodots         bool // :
	body            bytes.Buffer
}

type mdlink struct {
	title []string
	link  string
}

func parseMarkdownLinks(v string) (list []*mdlink) {
	rr := LinksMap(v)
	if rr == nil {
		return
	}

	v = strings.ReplaceAll(v, "![]", " ")
	v = strings.ReplaceAll(v, "![", " ")
	v = strings.ReplaceAll(v, "[", " ")
	v = strings.ReplaceAll(v, "]", " ")
	v = strings.ReplaceAll(v, "(", " ")
	v = strings.ReplaceAll(v, ")", " ")
	v = strings.ReplaceAll(v, "<", " ")
	v = strings.ReplaceAll(v, ">", " ")
	v = strings.ReplaceAll(v, `"`, " ")
	//v = strings.ReplaceAll(v, "http", " http")

	var s []string
	for _, x := range strings.Fields(v) {
		if rr[x] {
			//link
			list = append(list, &mdlink{s, x})
			s = nil
			continue
		}
		s = append(s, x)
	}
	list = append(list, &mdlink{title: s})

	return
}
