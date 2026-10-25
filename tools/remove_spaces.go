package tools

import (
	"regexp"
)

func removeAllSpaces(text string) (res string) {
	space := regexp.MustCompile(`\s+`)
	return space.ReplaceAllString(text, " ")

}
