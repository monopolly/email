package email

import (
	"fmt"
	"log"

	"github.com/monopolly/file"
)

type Options struct {
	ID                   string                            `json:"id,omitempty"`
	Images               bool                              `json:"images,omitempty"`   //download images
	Short                bool                              `json:"short,omitempty"`    //one paragraph
	MaxWords             int                               `json:"maxwords,omitempty"` //for short
	ShortLinks           func(...string) map[string]string `json:"-"`                  //short links
	ShortLink            func(string) string               `json:"-"`                  //short links
	ShortLinksDomain     string                            `json:"-"`                  //go.senthy.com
	ShortLinksMustScheme bool                              `json:"-"`
	MaxTextLen           int                               `json:"-"` //telegram max len
	Proxy                string                            `json:"-"`
	Debug                bool                              `json:"-"`
	DebugSave            bool                              `json:"-"`
	HTML                 bool                              `json:"-"` //prioritize
	MinImageSizeBytes    int                               `json:"-"` //5000
	NoGIF                bool                              `json:"-"`
	EscapeAddress        bool                              `json:"-"` //telegram escape strings.ReplaceAll(r.From.Address, "_", "\\_")

	html       bool //parsed from html
	debugcount int
}

func (a *Options) Log(arg ...interface{}) {
	if a.Debug {
		arg = append([]interface{}{a.ID}, arg...)
		log.Println(arg...)
	}
}

func (a *Options) Save(title string, body string) {
	if !a.Debug && !a.DebugSave {
		return
	}
	a.debugcount++
	file.Save(fmt.Sprintf("x__%d_%s", a.debugcount, title), []byte(body))
}

func (a *Options) Clean() {
	a.debugcount = 0
	file.DeleteMask("x__*")
}
