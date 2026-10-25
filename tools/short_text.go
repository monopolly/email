package tools

import (
	"strings"
)

/*
some words with lower case
other words with lower case,
or...
*/
func ShortText(s string, words int) (res string) {
	if words <= 0 {
		return s
	}

	count := 0
	var lines []string
	for _, i := range strings.Split(s, "\n") {
		c := len(strings.Fields(i))
		if count+c < words {
			lines = append(lines, i)
			count = count + c
			continue
		}

		// первая же строка длиннее лимита — режем её по словам
		if len(strings.Join(lines, "")) == 0 {
			lines = []string{strings.Join(strings.Fields(i)[:words], " ")}
		}

		// "..." ставим к последней непустой строке
		for len(lines) > 0 && strings.TrimSpace(lines[len(lines)-1]) == "" {
			lines = lines[:len(lines)-1]
		}
		lines[len(lines)-1] = lines[len(lines)-1] + "..."
		break
	}
	return strings.Join(lines, "\n")
}
