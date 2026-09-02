package email

import (
	"bytes"
	"fmt"
	"net/mail"
	"strings"
	"unicode/utf8"

	"github.com/jhillyerd/enmime/v2"
	"github.com/monopolly/filetype"
	"golang.org/x/net/html/charset"
	"golang.org/x/text/encoding/charmap"
)

func ParseRawEmail(raw []byte) (r *Email, err error) {

	var namecount int
	pp := enmime.NewParser()
	en, er := pp.ReadEnvelope(bytes.NewBuffer(raw))
	if er != nil {
		return
	}

	// subj charset
	subj := en.GetHeader("Subject")
	if !utf8.ValidString(subj) {
		subj = utf8convert(subj)
	}

	r = &Email{
		Subj:        subj,
		Text:        en.Text,
		HTML:        en.HTML,
		Unsubscribe: ParseUnsubscribeLink(en.GetHeader("List-Unsubscribe")),
	}

	froms, _ := en.AddressList("From")
	if len(froms) > 0 {
		r.From = *froms[0]
		r.From.Address = strings.ToLower(r.From.Address)
	}

	tos, _ := en.AddressList("To")
	if len(tos) > 0 {
		r.To = *tos[0]
		r.To.Address = strings.ToLower(r.To.Address)
	}

	hasAddress := map[string]bool{}
	// X-Forwarded-To: lava.mobi@gmail.com
	// X-Forwarded-For: vostorgbt.krl@gmail.com lava.mobi@gmail.com
	// Delivered-To: lava.mobi@gmail.com
	for _, tag := range []string{
		"X-Forwarded-To", "X-Forwarded-For", "Delivered-To",
		"In-Reply-To", "X-Original-To", "X-Loop", "References", "Resent-To", "Resent-Cc", "Resent-Bcc"} {

		// mine
		t := en.GetHeader(tag)
		if t != "" {
			t = strings.ReplaceAll(t, tag, "")
			t = strings.ReplaceAll(t, ",", " ")
			t = strings.ReplaceAll(t, ";", " ")
			t = strings.ToLower(t)
			list := strings.Fields(t)
			for _, x := range list {
				if !strings.Contains(x, "@") {
					continue
				}
				x = strings.ReplaceAll(x, "<", "")
				x = strings.ReplaceAll(x, ">", "")
				if hasAddress[x] {
					continue
				}
				r.Forward = append(r.Forward, &mail.Address{Address: x})
				hasAddress[x] = true
			}

		}

	}

	for _, x := range en.Inlines {

		f := new(File)
		f.Type = x.ContentType
		f.Name = x.FileName
		f.Body = x.Content
		f.Size = len(f.Body)

		if f.Name == "" {
			tp, ext := filetype.Detect(f.Body)
			f.Type = tp
			namecount++
			f.Name = fmt.Sprintf("%d.%s", namecount, ext)
		}

		f.Image = strings.Contains(f.Type, "image")

		f.ID = fmt.Sprintf("%d%s%s", len(f.Body), f.Name, f.Type)
		r.Inline = append(r.Inline, f)
		if f.Image {
			r.Images = append(r.Images, f)
		}
	}

	for _, x := range en.Attachments {

		//log.Println("attach", x)

		f := new(File)
		f.Type = x.ContentType
		f.Name = x.FileName
		f.Body = x.Content
		f.Size = len(f.Body)

		if f.Name == "" {
			tp, ext := filetype.Detect(f.Body)
			f.Type = tp
			namecount++
			f.Name = fmt.Sprintf("%d.%s", namecount, ext)
		}

		f.Image = strings.Contains(f.Type, "image")
		f.ID = fmt.Sprintf("%d%s%s", len(f.Body), f.Name, f.Type)
		r.Files = append(r.Files, f)
		if f.Image {
			r.Images = append(r.Images, f)
		}
	}

	return
}

func utf8convert(text string) (res string) {
	_, nn, _ := charset.DetermineEncoding([]byte(text), "")

	switch nn {
	case "windows-1252":
		d := charmap.Windows1251.NewDecoder()
		res, _ = d.String(text)
	default:
		res = text
	}
	return
}
