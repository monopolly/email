package tools

import (
	"fmt"
	"net"
	"net/url"
	"strings"

	"golang.org/x/net/publicsuffix"
)

// api.apple.com = apple.com, news.bbc.co.uk = bbc.co.uk
func parentHost(host string) (newhost string) {
	u, err := url.Parse(host)
	if err != nil {
		return host
	}
	hostname := u.Hostname()
	if net.ParseIP(hostname) != nil {
		return hostname
	}
	if domain, err := publicsuffix.EffectiveTLDPlusOne(hostname); err == nil {
		return domain
	}
	list := strings.Split(hostname, ".")
	if len(list) < 2 {
		return host
	}
	return fmt.Sprintf("%s.%s", list[len(list)-2], list[len(list)-1])
}
