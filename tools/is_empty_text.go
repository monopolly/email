package email

import (
	"unicode"
)

func emptyText(t string) (res bool) {
	for _, x := range t {
		if unicode.IsGraphic(x) {
			return
		}
	}
	return true
}
