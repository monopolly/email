package email

import (
	"bytes"
	"fmt"
	"net/mail"
	"strings"

	"github.com/jhillyerd/enmime/v2"
	"github.com/monopolly/filetype"

	"github.com/monopolly/email/tools"
	"github.com/monopolly/email/utf"
)

func Parse(raw []byte) (r *Email, err error) {

	var namecount int
	pp := enmime.NewParser()
	en, err := pp.ReadEnvelope(bytes.NewBuffer(raw))
	if err != nil {
		return nil, err
	}

	r = &Email{
		ID:          en.GetHeader("Message-Id"),
		Subj:        utf.Fix(en.GetHeader("Subject")),
		Text:        utf.Fix(en.Text, tools.PartCharset(en, "text/plain")),
		HTML:        utf.Fix(en.HTML, tools.PartCharset(en, "text/html")),
		Unsubscribe: tools.ParseUnsubscribeLink(en.GetHeader("List-Unsubscribe")),
	}

	if date, err := en.Date(); err == nil {
		r.Date = &date
	}

	froms, _ := en.AddressList("From")
	if len(froms) > 0 {
		r.From = *froms[0]
		r.From.Name = utf.Fix(r.From.Name)
		r.From.Address = strings.ToLower(r.From.Address)
	}

	tos, _ := en.AddressList("To")
	if len(tos) > 0 {
		r.To = *tos[0]
		r.To.Name = utf.Fix(r.To.Name)
		r.To.Address = strings.ToLower(r.To.Address)
	}

	r.CC = tools.AddressList(en, "Cc")
	r.BCC = tools.AddressList(en, "Bcc")

	hasAddress := map[string]bool{}
	// X-Forwarded-To: lava.mobi@gmail.com
	// X-Forwarded-For: vostorgbt.krl@gmail.com lava.mobi@gmail.com
	// Delivered-To: lava.mobi@gmail.com
	// In-Reply-To и References не адреса, а Message-Id писем (<id@host>)
	for _, tag := range []string{
		"X-Forwarded-To", "X-Forwarded-For", "Delivered-To",
		"X-Original-To", "X-Loop", "Resent-To", "Resent-Cc", "Resent-Bcc"} {

		// mine
		t := en.GetHeader(tag)
		if t != "" {
			t = strings.ReplaceAll(t, ",", " ")
			t = strings.ReplaceAll(t, ";", " ")
			t = strings.ToLower(t)
			list := strings.FieldsSeq(t)
			for x := range list {
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
		f.Name = utf.Fix(x.FileName)
		f.CID = x.ContentID
		f.Inline = true
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
		f.Name = utf.Fix(x.FileName)
		f.CID = x.ContentID
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
