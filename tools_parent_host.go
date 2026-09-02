package email

import (
	"fmt"
	"net/url"
	"strings"
)

// api.apple.com = apple.com
func parentHost(host string) (newhost string) {
	u, err := url.Parse(host)
	if err != nil {
		return host
	}
	list := strings.Split(u.Hostname(), ".")
	if len(list) < 2 {
		return host
	}
	return fmt.Sprintf("%s.%s", list[len(list)-2], list[len(list)-1])
}
