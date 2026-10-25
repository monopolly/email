package tools

import (
	"bytes"
	"unicode"
)

func cleanTextHard(text string) (res string) {
	var b bytes.Buffer
	var space bool
	for _, x := range text {
		switch {
		case
			unicode.IsDigit(x),
			unicode.IsLetter(x),
			unicode.IsSymbol(x),
			unicode.IsPunct(x):
			b.WriteRune(x)
			space = false
		default:
			if unicode.IsSpace(x) && !space {
				b.WriteString(" ")
				space = true
			}

		}
	}
	return b.String()
}
