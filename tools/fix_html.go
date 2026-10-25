package email

import (
	"strings"
)

func htmlfix(html string, options *Options) (res string) {

	res = html
	res = strings.TrimSpace(res)
	if len(res) < 20 {
		options.Log("html too short")
		return
	}

	if strings.Contains(strings.ToLower(res), `utf-8`) {
		options.Log("html has utf8")
		return
	}

	//html
	if !strings.Contains(res, "<html") {
		res = `<!DOCTYPE html><html><head><meta charset="utf-8"></head><body>` + res + "</body></html>"
		return
	}

	i := strings.Index(res, "<head>")
	switch i > 0 {
	case true:
		i = i + 6
		res = res[:i] + `<meta charset="utf-8">` + res[i:]
	default:
		i = strings.Index(res, "<html")
		if i < 0 {
			break
		}
		end := strings.Index(res[i:], ">")
		if end == -1 {
			return `<html><head><meta charset="utf-8"></head><body>` + res + "</body></html>" //invalid
		}
		x := i + end
		res = res[:x] + `<head><meta charset="utf-8"></head>` + res[x:]

	}

	return
}
