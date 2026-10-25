package email

import (
	"fmt"
	"strings"
	"unicode"
	"unicode/utf8"
)

/*
some words with lower case
other words with lower case,
or words with comma
*/
func retext(s string) (res string) {
	var lines []string

	s = strings.ReplaceAll(s, "_", " ")

	for _, x := range strings.Split(s, "\n") {
		x = strings.TrimSpace(x)

		switch lines == nil {
		case true:
			lines = append(lines, x)
		default:
			if x == "" {
				lines = append(lines, x)
				continue
			}
			if lines[len(lines)-1] == "" {
				lines = append(lines, x)
				continue
			}

			now, _ := utf8.DecodeLastRuneInString(x)
			first, _ := utf8.DecodeRune([]byte(x))
			last, _ := utf8.DecodeLastRuneInString(lines[len(lines)-1])

			switch {
			case unicode.IsLower(now) && unicode.IsUpper(last):
				//fmt.Println(1)
				lines = addToLastNewDotLineText(lines, x)
			case unicode.IsLower(first) && unicode.IsLower(last):
				lines = addToLastLineText(lines, x)
				//fmt.Println(2)
			case string(last) == "," && unicode.IsLower(now):
				lines = addToLastLineText(lines, x)
				//fmt.Println(3)
			case string(last) == "-" && unicode.IsLower(now):
				lines = addToLastLineText(lines, x)
				//fmt.Println(4)
			case string(last) == "." && unicode.IsUpper(first):
				lines = addToLastLineText(lines, x)
				//fmt.Println(5)
			case unicode.IsNumber(first) && unicode.IsNumber(last):
				//202 3. fix
				lines = addToLastLineText(lines, x)
			default:
				//fmt.Println("❌", x, string(now), ":", string(last))
				//fmt.Println(x)
				lines = append(lines, x)
				continue
			}

			//fmt.Println("✅", x, string(now), ":", string(last))

		}

	}

	//list fix 1. 2. 3.
	res = strings.Join(lines, "\n")
	for i := 1; i < 20; i++ {
		res = strings.ReplaceAll(res, fmt.Sprintf(" %d. ", i), fmt.Sprintf("\n%d. ", i))
	}

	res = strings.ReplaceAll(res, ")\n\n[", ")\n[")
	return
}

func addToLastLineText(lines []string, v string) []string {
	if lines == nil {
		return []string{v}
	}
	if lines[len(lines)-1] == "" {
		return lines
	}

	lines[len(lines)-1] = lines[len(lines)-1] + " " + v
	return lines
}

func addToLastNewDotLineText(lines []string, v string) []string {
	if lines == nil {
		return []string{v}
	}
	if lines[len(lines)-1] == "" {
		return lines
	}

	lines[len(lines)-1] = lines[len(lines)-1] + ". " + v
	return lines
}
