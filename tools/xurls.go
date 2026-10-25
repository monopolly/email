package email

import (
	"regexp"

	"mvdan.cc/xurls/v2"
)

// xurls builds a very large regexp on every call, so compile each one once.
var (
	xurlsStrict = xurls.Strict()
	// The scheme must include the colon: "tel" alone also matches
	// "telnogo-belya" and "hotel.pro/page".
	xurlsTel    = strictScheme("tel:")
	xurlsMailto = strictScheme("mailto:")
)

func strictScheme(scheme string) *regexp.Regexp {
	rx, err := xurls.StrictMatchingScheme(scheme)
	if err != nil {
		panic("email: bad xurls scheme " + scheme + ": " + err.Error())
	}
	return rx
}
