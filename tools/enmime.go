package tools

import (
	"net/mail"
	"strings"

	"github.com/jhillyerd/enmime/v2"

	"github.com/monopolly/email/utf"
)

// AddressList — адреса из заголовка письма: имена в utf-8, адреса в нижнем регистре
func AddressList(en *enmime.Envelope, key string) (res []*mail.Address) {
	list, _ := en.AddressList(key)
	for _, x := range list {
		x.Name = utf.Fix(x.Name)
		x.Address = strings.ToLower(x.Address)
		res = append(res, x)
	}
	return
}

// PartCharset — объявленная кодировка первой части нужного типа (если enmime не смог перекодировать)
func PartCharset(en *enmime.Envelope, contentType string) string {
	if en.Root == nil {
		return ""
	}
	p := en.Root.BreadthMatchFirst(func(p *enmime.Part) bool {
		return p.ContentType == contentType
	})
	if p == nil {
		return ""
	}
	return p.Charset
}
