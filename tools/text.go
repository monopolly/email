package tools

import "strings"

// Domain — домен из email-адреса в нижнем регистре: a@B.com -> b.com
func Domain(email string) (domain string) {
	i := strings.Index(email, "@")
	if i == -1 {
		return email
	}
	domain = email[i+1:]
	domain = strings.TrimSpace(domain)
	domain = strings.ToLower(domain)
	return
}

// NormalizeNewlines: \r\n и \r -> \n
func NormalizeNewlines(s string) string {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	return strings.ReplaceAll(s, "\r", "\n")
}
