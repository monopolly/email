package tools

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
func Retext(s string) (res string) {
	var lines []string

	for _, x := range strings.Split(s, "\n") {
		x = strings.TrimSpace(removeUnderscores(x))

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
	for i, x := range lines {
		lines[i] = splitEnumeration(x)
	}
	res = strings.Join(lines, "\n")

	res = strings.ReplaceAll(res, ")\n\n[", ")\n[")
	return
}

// "_" (курсив) убираем везде, кроме ссылок и адресов: в них "_" — часть адреса
func removeUnderscores(line string) string {
	if !strings.Contains(line, "_") {
		return line
	}
	words := strings.Split(line, " ")
	for i, w := range words {
		if strings.Contains(w, "://") || strings.Contains(w, "www.") || strings.Contains(w, "@") ||
			(strings.Contains(w, "/") && strings.Contains(w, ".")) {
			continue
		}
		words[i] = strings.ReplaceAll(w, "_", " ")
	}
	return strings.Join(words, " ")
}

// "Шаги: 1. открыть 2. нажать 3. готово" -> каждый пункт с новой строки.
// Разбиваем, только если после "1. " идёт " 2. ": одиночное "версия 1. ..." не трогаем.
func splitEnumeration(line string) string {
	start := -1
	switch {
	case strings.HasPrefix(line, "1. "):
		start = 0
	case strings.Contains(line, " 1. "):
		start = strings.Index(line, " 1. ") + 1
	}
	if start < 0 || !strings.Contains(line[start:], " 2. ") {
		return line
	}

	var b strings.Builder
	if start > 0 {
		b.WriteString(strings.TrimSpace(line[:start]))
		b.WriteString("\n")
	}
	rest := line[start:]
	for n := 2; ; n++ {
		i := strings.Index(rest, fmt.Sprintf(" %d. ", n))
		if i < 0 {
			break
		}
		b.WriteString(rest[:i])
		b.WriteString("\n")
		rest = rest[i+1:]
	}
	b.WriteString(rest)
	return b.String()
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
