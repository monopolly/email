package tools

import (
	"strings"
)

// only http* links
// no apple.com
func DeleteRelativeLinks(text string) (res string) {

	//fmt.Println("delete relative links")
	res = text
	for _, v := range strings.Split(res, "\n") {
		v = strings.ReplaceAll(v, "![]", "")
		v = strings.ReplaceAll(v, "[", " ")
		v = strings.ReplaceAll(v, "]", " ")
		v = strings.ReplaceAll(v, "(", " ")
		v = strings.ReplaceAll(v, ")", " ")
		v = strings.ReplaceAll(v, "<", " ")
		v = strings.ReplaceAll(v, ">", " ")
		for _, f := range strings.Fields(v) {
			if strings.Contains(f, "/") && strings.Contains(f, ".") {
				if strings.Contains(f, "http") {
					continue
				}
				//fmt.Println("match", f)
				res = strings.ReplaceAll(res, f, "")
			}
		}
	}
	return
}
