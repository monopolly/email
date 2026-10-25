package email

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/mail"
	"net/url"
	"path/filepath"
	"strings"
	"time"

	"github.com/jhillyerd/enmime/v2"

	"github.com/monopolly/errors"
	"github.com/monopolly/file"
	"github.com/monopolly/filetype"

	"github.com/monopolly/email/tools"
)

func New() (a *Constructor) {
	a = new(Constructor)
	a.m = new(Email)
	t1 := time.Now()
	a.m.Date = &t1
	return
}

type Email struct {
	ID          string          `json:"id,omitempty"`
	From        mail.Address    `json:"from,omitempty"`
	To          mail.Address    `json:"to,omitempty"`
	Date        *time.Time      `json:"date,omitempty"`
	Subj        string          `json:"subj,omitempty"`
	BCC         []*mail.Address `json:"bcc,omitempty"`
	CC          []*mail.Address `json:"cc,omitempty"`
	Forward     []*mail.Address `json:"forward,omitempty"`
	Text        string          `json:"text,omitempty"`
	Markdown    string          `json:"markdown,omitempty"`
	HTML        string          `json:"html,omitempty"`
	Signature   []string        `json:"signature,omitempty"`
	Inline      []*File         `json:"inline,omitempty"`
	Files       []*File         `json:"files,omitempty"`
	Images      []*File         `json:"images,omitempty"`
	IP          string          `json:"ip,omitempty"`
	Unsubscribe *url.URL        `json:"unsubscribe,omitempty"` //link
	ParserHTML  bool            `json:"parserhtml,omitempty"`  //html or text
	EML         []byte          `json:"eml,omitempty"`         //html or text
}

func Unpack(b []byte) (a Email, err error) {
	err = json.Unmarshal(b, &a)
	return
}

func (a *Email) DomainFrom() (domain string) {
	return tools.Domain(a.From.Address)
}

func (a *Email) DomainTo() (domain string) {
	return tools.Domain(a.To.Address)
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
	s := len(body)
	if len(size) > 0 {
		s = size[0]
	}
	addFile(&a.Files, &File{ID: tools.SHA1(body), Name: filename, Body: body, Type: contentType, Size: s})
}
func (a *Email) AddFileModel(v *File) {
	v.ID = tools.SHA1(v.Body)
	v.Size = len(v.Body)
	addFile(&a.Files, v)
}
func (a *Email) AddImageModel(v *File) {
	v.ID = tools.SHA1(v.Body)
	v.Size = len(v.Body)
	addFile(&a.Images, v)
}

func (a *Email) AddInlineModel(v *File) {
	v.ID = tools.SHA1(v.Body)
	v.Size = len(v.Body)
	addFile(&a.Inline, v)
}

func (a *Email) AddImage(filename, contentType string, body []byte, size ...int) {
	s := len(body)
	if len(size) > 0 {
		s = size[0]
	}
	addFile(&a.Images, &File{ID: tools.SHA1(body), Name: filename, Body: body, Type: contentType, Size: s})
}

func (a *Email) AddInline(filename, contentType string, body []byte) {
	addFile(&a.Inline, &File{ID: tools.SHA1(body), Name: filename, Body: body, Type: contentType})
}

func (a *Email) AddFiles(filename ...string) {
	for _, x := range filename {
		f, err := openFile(x)
		if err != nil {
			fmt.Println("File error", err)
			continue
		}
		addFile(&a.Files, f)
	}
}

func (a *Email) AddImages(filename ...string) {
	for _, x := range filename {
		f, err := openFile(x)
		if err != nil {
			fmt.Println("File error", err)
			continue
		}
		addFile(&a.Images, f)
	}
}

func (a *Email) AddInlines(filename ...string) {
	for _, x := range filename {
		f, err := openFile(x)
		if err != nil {
			continue
		}
		f.Size = 0
		addFile(&a.Inline, f)
	}
}

// addFile добавляет файл, если файла с таким же содержимым (ID) в списке ещё нет
func addFile(list *[]*File, f *File) {
	for _, x := range *list {
		if x.ID == f.ID {
			return
		}
	}
	*list = append(*list, f)
}

func openFile(path string) (*File, error) {
	b, err := file.Open(path)
	if err != nil {
		return nil, err
	}
	types, _ := filetype.Detect(b)
	return &File{ID: tools.SHA1(b), Name: filepath.Base(path), Body: b, Type: types, Size: len(b)}, nil
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
	CID    string `json:"cid,omitempty"`  //Content-ID inline файла; в html: src="cid:..."
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

func (a *Constructor) Subj(v string, args ...any) *Constructor {
	if len(args) > 0 {
		v = fmt.Sprintf(v, args...)
	}
	a.m.Subj = v
	return a
}

func (a *Constructor) Signature(lines ...string) *Constructor {
	a.m.Signature = lines
	return a
}

func (a *Constructor) UnsubscribeLink(link *url.URL) *Constructor {
	a.m.Unsubscribe = link
	return a
}

// Deprecated: используйте UnsubscribeLink.
func (a *Constructor) UnsibscribeLink(link *url.URL) *Constructor {
	return a.UnsubscribeLink(link)
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

func (a *Constructor) AddFile(filename, contentType string, body []byte, size ...int) *Constructor {
	a.m.AddFile(filename, contentType, body, size...)
	return a
}

func (a *Constructor) AddInline(filename, contentType string, body []byte) *Constructor {
	a.m.AddInline(filename, contentType, body)
	return a
}

func (a *Constructor) AddFiles(filename ...string) *Constructor {
	a.m.AddFiles(filename...)
	return a
}

func (a *Constructor) AddInlines(filename ...string) *Constructor {
	a.m.AddInlines(filename...)
	return a
}

func (a *Constructor) Done() *Email {
	return a.m
}

func (a *Email) Render() (res bytes.Buffer) {

	domain := tools.Domain(a.From.Address)

	date := time.Now()
	if a.Date != nil {
		date = *a.Date
	}

	/* id */
	if a.ID == "" {
		a.ID = fmt.Sprintf("<%d.%s@%s>", date.Unix(), tools.RandomID(), domain)
	}

	/* mail mimes */
	b := enmime.Builder().
		Date(date).
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

	/* inline: в html на файл ссылаются как src="cid:<CID или ID>" */
	for _, x := range a.Inline {
		cid := x.CID
		if cid == "" {
			cid = x.ID
		}
		if cid == "" {
			cid = tools.SHA1(x.Body)
		}
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
