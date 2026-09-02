package email

import (
	"bytes"
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math/rand"
	"net/mail"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/jhillyerd/enmime/v2"

	"github.com/monopolly/errors"
	"github.com/monopolly/file"
	"github.com/monopolly/filetype"
)

func New() (a *Constructor) {
	a = new(Constructor)
	a.m = new(Email)
	t1 := time.Now()
	a.m.Date = &t1
	return
}

type Email struct {
	ID          string          `json:",omitempty"`
	From        mail.Address    `json:",omitempty"`
	To          mail.Address    `json:",omitempty"`
	Date        *time.Time      `json:",omitempty"`
	Subj        string          `json:",omitempty"`
	BCC         []*mail.Address `json:",omitempty"`
	CC          []*mail.Address `json:",omitempty"`
	Forward     []*mail.Address `json:",omitempty"`
	Text        string          `json:",omitempty"`
	Markdown    string          `json:",omitempty"`
	HTML        string          `json:",omitempty"`
	Signature   []string        `json:",omitempty"`
	Inline      []*File         `json:",omitempty"`
	Files       []*File         `json:",omitempty"`
	Images      []*File         `json:",omitempty"`
	IP          string          `json:",omitempty"`
	Unsubscribe *url.URL        `json:",omitempty"` //link
	ParserHTML  bool            `json:",omitempty"` //html or text
	Options     *Options        `json:",omitempty"` //html or text
	EML         []byte          `json:",omitempty"` //html or text
}

func Unpack(b []byte) (a Email, err error) {
	err = json.Unmarshal(b, &a)
	return
}

func (a *Email) DomainFrom() (domain string) {
	return domainonly(a.From.Address)
}

func (a *Email) DomainTo() (domain string) {
	return domainonly(a.To.Address)
}

func (a *Email) CCList() (res []string) {
	for _, x := range a.CC {
		res = append(res, x.Address)
	}
	return
}
func (a *Email) BCCList() (res []string) {
	for _, x := range a.BCC {
		res = append(res, x.Address)
	}
	return
}
func (a *Email) ForwardList() (res []string) {
	for _, x := range a.Forward {
		res = append(res, x.Address)
	}
	return
}

func (a *Email) Pack() (b []byte) {
	b, _ = json.Marshal(a)
	return
}

func (a *Email) Valid() (err errors.E) {
	/* validation */
	switch {
	case a.Subj == "":
		err = errors.Empty("subj")
		return
	case a.From.Address == "":
		err = errors.Empty("from")
		return
	case a.To.Address == "":
		err = errors.Empty("to")
		return
	default:
		return
	}
}

func (a *Email) AddFile(filename, contentType string, body []byte, size ...int) {
	id := sid(body)
	for _, x := range a.Files {
		if x.ID == id {
			return
		}
	}
	s := len(body)
	if len(size) > 0 {
		s = size[0]
	}

	a.Files = append(a.Files, &File{ID: id, Name: filename, Body: body, Type: contentType, Size: s})
}
func (a *Email) AddFileModel(v *File) {
	v.ID = sid(v.Body)
	for _, x := range a.Files {
		if x.ID == v.ID {
			return
		}
	}
	v.Size = len(v.Body)
	a.Files = append(a.Files, v)
}
func (a *Email) AddImageModel(v *File) {
	v.ID = sid(v.Body)
	for _, x := range a.Images {
		if x.ID == v.ID {
			return
		}
	}
	v.Size = len(v.Body)
	a.Images = append(a.Images, v)
}

func (a *Email) AddInlineModel(v *File) {
	v.ID = sid(v.Body)
	for _, x := range a.Inline {
		if x.ID == v.ID {
			return
		}
	}
	v.Size = len(v.Body)
	a.Inline = append(a.Inline, v)
}

func (a *Email) AddImage(filename, contentType string, body []byte, size ...int) {
	id := sid(body)
	for _, x := range a.Images {
		if x.ID == id {
			return
		}
	}
	s := len(body)
	if len(size) > 0 {
		s = size[0]
	}

	a.Images = append(a.Images, &File{ID: id, Name: filename, Body: body, Type: contentType, Size: s})
}

func (a *Email) AddInline(filename, contentType string, body []byte) {
	id := sid(body)
	for _, x := range a.Inline {
		if x.ID == id {
			return
		}
	}
	a.Inline = append(a.Inline, &File{ID: id, Name: filename, Body: body, Type: contentType})
}

func (a *Email) AddFiles(filename ...string) {
	for _, x := range filename {
		b, err := file.Open(x)
		if err != nil {
			fmt.Println("File error", err)
			continue
		}
		id := sid(b)
		for _, x := range a.Files {
			if x.ID == id {
				continue
			}
		}

		name := filepath.Base(x)
		types, _ := filetype.Detect(b)
		a.Files = append(a.Files, &File{ID: id, Name: name, Body: b, Type: types, Size: len(b)})
	}
}

func (a *Email) AddImages(filename ...string) {
	for _, x := range filename {
		b, err := file.Open(x)
		if err != nil {
			fmt.Println("File error", err)
			continue
		}
		id := sid(b)
		for _, x := range a.Images {
			if x.ID == id {
				continue
			}
		}

		name := filepath.Base(x)
		types, _ := filetype.Detect(b)
		a.Images = append(a.Images, &File{ID: id, Name: name, Body: b, Type: types, Size: len(b)})
	}
}

func (a *Email) AddInlines(filename ...string) {
	for _, x := range filename {
		b, err := file.Open(x)
		if err != nil {
			//fmt.Println("File error", err)
			continue
		}

		id := sid(b)
		for _, x := range a.Inline {
			if x.ID == id {
				continue
			}
		}

		name := filepath.Base(x)
		types, _ := filetype.Detect(b)
		//fmt.Println(name, types)
		a.Inline = append(a.Inline, &File{ID: id, Name: name, Body: b, Type: types})
	}
}

type File struct {
	ID     string `json:"id,omitempty"`
	Name   string `json:"name,omitempty"`
	Type   string `json:"type,omitempty"`
	Body   []byte `json:"body,omitempty"`
	Path   string `json:"path,omitempty"`
	Size   int    `json:"size,omitempty"` //for sending only
	Inline bool   `json:"inline,omitempty"`
	Image  bool   `json:"image,omitempty"`
	Link   string `json:"link,omitempty"` //origin link
}

type Constructor struct {
	m *Email
}

func (a *Constructor) To(email string, name ...string) *Constructor {
	p := mail.Address{Address: email}
	if len(name) > 0 {
		p.Name = name[0]
	}
	a.m.To = p
	return a
}

func (a *Constructor) BCC(email string, name ...string) *Constructor {
	p := &mail.Address{Address: email}
	if len(name) > 0 {
		p.Name = name[0]
	}
	a.m.BCC = append(a.m.BCC, p)
	return a
}

func (a *Constructor) CC(email string, name ...string) *Constructor {
	p := &mail.Address{Address: email}
	if len(name) > 0 {
		p.Name = name[0]
	}
	a.m.CC = append(a.m.CC, p)
	return a
}

func (a *Constructor) Subj(v string, args ...interface{}) *Constructor {
	a.m.Subj = fmt.Sprintf(v, args...)
	return a
}

func (a *Constructor) Signature(lines ...string) *Constructor {
	a.m.Signature = lines
	return a
}

func (a *Constructor) UnsibscribeLink(link *url.URL) *Constructor {
	a.m.Unsubscribe = link
	return a
}

func (a *Constructor) From(email string, name ...string) *Constructor {
	a.m.From = mail.Address{Address: email}
	if len(name) > 0 {
		a.m.From.Name = name[0]
	}
	return a
}

func (a *Constructor) Text(v string) *Constructor {
	a.m.Text = v
	return a
}

func (a *Constructor) HTML(v string) *Constructor {
	a.m.HTML = v
	return a
}

// sid возвращает content-id вложения.
// sha1.Sum, а не общий hash.Hash: письма собираются параллельно,
// общее состояние хэшера давало гонку и неверные cid.
func sid(v []byte) string {
	if v == nil {
		return ""
	}
	sum := sha1.Sum(v)
	return hex.EncodeToString(sum[:])
}

func (a *Constructor) AddFile(filename, contentType string, body []byte, size ...int) *Constructor {
	id := sid(body)
	for _, x := range a.m.Files {
		if x.ID == id {
			return a
		}
	}
	s := len(body)
	if len(size) > 0 {
		s = size[0]
	}

	a.m.Files = append(a.m.Files, &File{ID: id, Name: filename, Body: body, Type: contentType, Size: s})
	return a
}

func (a *Constructor) AddInline(filename, contentType string, body []byte) *Constructor {
	id := sid(body)
	for _, x := range a.m.Inline {
		if x.ID == id {
			return a
		}
	}
	a.m.Inline = append(a.m.Inline, &File{ID: id, Name: filename, Body: body, Type: contentType})
	return a
}

func (a *Constructor) AddFiles(filename ...string) *Constructor {
	for _, x := range filename {
		b, err := file.Open(x)
		if err != nil {
			fmt.Println("File error", err)
			continue
		}
		id := sid(b)
		for _, x := range a.m.Files {
			if x.ID == id {
				continue
			}
		}

		name := filepath.Base(x)
		types, _ := filetype.Detect(b)
		//fmt.Println(name, types)
		a.m.Files = append(a.m.Files, &File{ID: id, Name: name, Body: b, Type: types, Size: len(b)})
	}
	return a
}

func (a *Constructor) AddInlines(filename ...string) *Constructor {
	for _, x := range filename {
		b, err := file.Open(x)
		if err != nil {
			//fmt.Println("File error", err)
			continue
		}

		id := sid(b)
		for _, x := range a.m.Inline {
			if x.ID == id {
				continue
			}
		}

		name := filepath.Base(x)
		types, _ := filetype.Detect(b)
		//fmt.Println(name, types)
		a.m.Inline = append(a.m.Inline, &File{ID: id, Name: name, Body: b, Type: types})
	}
	return a
}

func (a *Constructor) Done() *Email {
	return a.m
}

func domainonly(email string) (domain string) {
	i := strings.Index(email, "@")
	if i == -1 {
		return email
	}
	domain = email[i+1:]
	domain = strings.TrimSpace(domain)
	domain = strings.ToLower(domain)
	return
}

func (a *Email) Render() (res bytes.Buffer) {

	domain := domainonly(a.From.Address)

	/* id */
	if a.ID == "" {
		a.ID = fmt.Sprintf("<%d@%d.%s>", a.Date.Unix(), rand.Intn(999), domain)
	}

	/* mail mimes */
	b := enmime.Builder().
		Date(*a.Date).
		Header("Message-Id", a.ID).
		From(a.From.Name, a.From.Address).
		To(a.To.Name, a.To.Address).
		Subject(a.Subj)

	if a.Unsubscribe != nil {
		b = b.Header("List-Unsubscribe", a.Unsubscribe.String())
	}

	/* text body */
	if a.Text != "" {
		s := a.Text + "\n"
		s += "\n\n"
		s += strings.Join(a.Signature, "\n")
		b = b.Text([]byte(strings.TrimSpace(s)))
	}

	/* html body */
	if a.HTML != "" {
		s := a.HTML
		s += "\n\n"
		s += strings.Join(a.Signature, "\n")
		b = b.HTML([]byte(strings.TrimSpace(s)))
	}

	/* bcc */
	for _, x := range a.BCC {
		b = b.BCC(x.Name, x.Address)
	}

	/* cc */
	for _, x := range a.CC {
		b = b.CC(x.Name, x.Address)
	}

	/* files */
	for _, x := range a.Files {
		b = b.AddAttachment(x.Body, x.Type, x.Name)
	}

	/* inline */
	for _, x := range a.Inline {
		cid := fmt.Sprintf("<cid.%d@inline.%d.%s>", rand.Int63(), rand.Intn(999), domain)
		b = b.AddInline(x.Body, x.Type, x.Name, cid)
	}

	/* render */
	m, err := b.Build()
	if err != nil {
		fmt.Println(err)
		return
	}

	/* encode */
	m.Encode(&res)
	return
}
