package email

import (
	"strings"
)

/*
some words with lower case
other words with lower case,
or...
*/
func shortText(s string, words int) (res string) {

	count := 0
	var lines []string
	for _, i := range strings.Split(s, "\n") {
		c := len(strings.Fields(i))
		if count+c < words {
			lines = append(lines, i)
			count = count + c
			continue
		}
		if len(lines) > 0 {
			lines[len(lines)-1] = lines[len(lines)-1] + "..."
			break
		}

	}
	return strings.Join(lines, "\n")
}
