package email

import (
	"strings"
)

func TelegramMarkdownEscape(text string) (res string) {
	res = text
	res = strings.ReplaceAll(res, "_", "\\_")
	res = strings.ReplaceAll(res, "*", "\\*")
	return
}

func TelegramMarkdownRestore(text string) (res string) {
	res = text
	res = strings.ReplaceAll(res, "\\_", "_")
	res = strings.ReplaceAll(res, "\\*", "*")
	return
}
