package parse

//testing

import (
	"encoding/json"
	"fmt"
	"log"
	"runtime"
	"strings"
	"testing"

	"github.com/monopolly/file"
	"github.com/monopolly/images"
	"github.com/stretchr/testify/assert"
	"github.com/valyala/fasthttp"

	"github.com/monopolly/email"
)

func TestDownload(ggggg *testing.T) {
	function, _, _, _ := runtime.Caller(0)
	fn := runtime.FuncForPC(function).Name()
	fn = fn[strings.LastIndex(fn, ".Test")+5:]
	fn = strings.Join(strings.Split(fn, "_"), ": ")
	fmt.Printf("\033[1;32m%s\033[0m\n", fn)

	ggggg.Skip("ручной тест")
	files := []string{
		"https://www.e-xecutive.ru/assets/mail/social/linkedin-d80e9a9ff2085061efd4efc4c1a511c35034373a1995369fe38e38e28882e70e.png",
		"https://www.e-xecutive.ru/assets/mail/social/telegram-10242950a80856c7c5b3ed8244be0a1e7b478bab40a86be2e9c3e775f8da84ac.png",
		"https://www.e-xecutive.ru/assets/mail/social/twitter-0f764750af7b71666777c6bab349c8de603b17dda5630dbfe72b19b8c8ba90e5.png",
		"https://www.e-xecutive.ru/assets/mail/social/vkontakte-d94c3fe46ffb7d8f25d10cb97b52598328569a12e7488393639043e8ca6e74b6.png",
		"https://www.e-xecutive.ru/uploads/article/image/1996648/mini_business.jpg",
		"https://www.e-xecutive.ru/uploads/article/image/1996651/mini_tehnopolis-moskva.jpg",
		"https://www.e-xecutive.ru/uploads/article/image/1996653/mini_time.jpg",
		"https://www.e-xecutive.ru/uploads/news/image/2097293/mini_molchanov.jpeg",
	}

	_, b, _ := fasthttp.Get(nil, files[0])
	if b == nil {
		panic("nil file")
	}

	p, err := images.New(b)
	if err != nil {
		panic("image")
	}

	fmt.Println(p.Ext)

	select {}
}
func TestParseText(ggggg *testing.T) {
	function, _, _, _ := runtime.Caller(0)
	fn := runtime.FuncForPC(function).Name()
	fn = fn[strings.LastIndex(fn, ".Test")+5:]
	fn = strings.Join(strings.Split(fn, "_"), ": ")
	fmt.Printf("\033[1;32m%s\033[0m\n", fn)

	ggggg.Skip("ручной тест")
	a := assert.New(ggggg)
	_ = a

	file1 := []string{"../mail/imagelinks.eml"}
	_ = file1

	file2 := file.FileList("../mail", "eml")
	_ = file2

	for _, x := range file1 { //  {
		x = strings.TrimPrefix(x, "../mail/")
		name := strings.TrimSuffix(x, ".eml")
		//fmt.Println(x)

		options := &Options{
			ShortLinks: func(s ...string) (res map[string]string) {
				res = make(map[string]string)
				for _, x := range s {
					res[x] = "testlink"
				}
				return
			},
			ShortLink:        func(s string) string { return "testlink" },
			ShortLinksDomain: "go.senthy.com",
			//Proxy:                "http://144.217.197.151:38611",
			Debug:                true,
			MinImageSizeBytes:    10_000,
			ShortLinksMustScheme: true,
			Images:               true,
			Short:                false,
			DebugSave:            true,
			HTML:                 true,
		}

		options.Clean()

		r, err := email.Parse(file.OpenE(fmt.Sprintf("../mail/%s.eml", name)))
		// fmt.Println(r)
		// fmt.Println("text", r.Text, len(r.Text))
		// fmt.Println("html", r.HTML, len(r.HTML))

		Parse(r, options)

		// if !r.ParserHTML {
		// 	continue
		// }

		if err != nil {
			fmt.Println(err)
			return
		}

		fmt.Println("markdown", len(r.Text))
		fmt.Println("html", len(r.HTML))
		fmt.Println("files", len(r.Files))
		fmt.Println("images", len(r.Images))

		file.Save(fmt.Sprintf("../mail/%s.html", name), []byte(r.HTML))
		file.Save(fmt.Sprintf("../mail/%s.txt", name), []byte(r.Text))
		file.Save(fmt.Sprintf("../mail/%s.md", name), []byte(r.Markdown))

		for _, fn := range r.Images {
			if fn.Size > options.MinImageSizeBytes {
				ff := fmt.Sprintf("../mail/%s_%s", name, fn.Name)
				options.Log("save image", fn.Name)
				file.Save(ff, fn.Body)
			} else {
				options.Log("image size too small")
			}
		}

		for _, fn := range r.Files {
			if fn.Size > options.MinImageSizeBytes {
				ff := fmt.Sprintf("../mail/%s_%s", name, fn.Name)
				options.Log("save file", fn.Name)
				file.Save(ff, fn.Body)
			} else {
				options.Log("image size too small")
			}
		}
	}

}

func TestUnsibscribe(ggggg *testing.T) {
	function, _, _, _ := runtime.Caller(0)
	fn := runtime.FuncForPC(function).Name()
	fn = fn[strings.LastIndex(fn, ".Test")+5:]
	fn = strings.Join(strings.Split(fn, "_"), ": ")
	fmt.Printf("\033[1;32m%s\033[0m\n", fn)

	ggggg.Skip("ручной тест")
	a := assert.New(ggggg)
	_ = a

	file1 := "../mail/unsubscribeLink.eml"
	// file1 := "../mail/unsubscribeEmail.eml"

	r, _ := email.Parse(file.OpenE(file1))
	fmt.Println(r.Unsubscribe)
}

func TestForward(ggggg *testing.T) {
	function, _, _, _ := runtime.Caller(0)
	fn := runtime.FuncForPC(function).Name()
	fn = fn[strings.LastIndex(fn, ".Test")+5:]
	fn = strings.Join(strings.Split(fn, "_"), ": ")
	fmt.Printf("\033[1;32m%s\033[0m\n", fn)

	// return
	a := assert.New(ggggg)
	_ = a

	ggggg.Skip("ручной тест")

	file1 := "../mail/redirect.eml"
	// file1 := "../mail/unsubscribeEmail.eml"

	r, err := email.Parse(file.OpenE(file1))
	if err != nil {
		fmt.Println(err)
		return
	}
	for _, x := range r.Forward {
		fmt.Println(x.Address)
	}

}
func TestNoRead(ggggg *testing.T) {
	function, _, _, _ := runtime.Caller(0)
	fn := runtime.FuncForPC(function).Name()
	fn = fn[strings.LastIndex(fn, ".Test")+5:]
	fn = strings.Join(strings.Split(fn, "_"), ": ")
	fmt.Printf("\033[1;32m%s\033[0m\n", fn)

	ggggg.Skip("ручной тест")

	// return
	a := assert.New(ggggg)
	_ = a

	file1 := "../mail/noread.eml"
	// file1 := "../mail/unsubscribeEmail.eml"

	r, err := email.Parse(file.OpenE(file1))
	if err != nil {
		fmt.Println(err)
		return
	}

	options := &Options{
		ShortLinks: func(s ...string) (res map[string]string) {
			res = make(map[string]string)
			for _, x := range s {
				res[x] = "testlink"
			}
			return
		},
		ShortLink:        func(s string) string { return "testlink" },
		ShortLinksDomain: "go.senthy.com",
		//Proxy:                "http://144.217.197.151:38611",
		Debug:                true,
		MinImageSizeBytes:    10_000,
		ShortLinksMustScheme: true,
		Images:               true,
		Short:                false,
		DebugSave:            true,
		HTML:                 true,
	}

	options.Clean()

	Parse(r, options)

	fmt.Println(r.Subj)

}

func TestRaw(ggggg *testing.T) {
	function, _, _, _ := runtime.Caller(0)
	fn := runtime.FuncForPC(function).Name()
	fn = fn[strings.LastIndex(fn, ".Test")+5:]
	fn = strings.Join(strings.Split(fn, "_"), ": ")
	fmt.Printf("\033[1;32m%s\033[0m\n", fn)

	// return
	a := assert.New(ggggg)
	_ = a

	// file1 := "../mail/rawmail_error.eml"
	file1 := "../mail/raw_error2.eml"
	// file1 := "../mail/unsubscribeEmail.eml"

	r, err := email.Parse(file.OpenE(file1))
	if err != nil {
		fmt.Println(err)
		return
	}

	options := &Options{
		ShortLinks: func(s ...string) (res map[string]string) {
			res = make(map[string]string)
			for _, x := range s {
				res[x] = "testlink"
			}
			return
		},
		ShortLink:        func(s string) string { return "testlink" },
		ShortLinksDomain: "go.senthy.com",
		//Proxy:                "http://144.217.197.151:38611",
		Debug:                true,
		MinImageSizeBytes:    10_000,
		ShortLinksMustScheme: true,
		Images:               true,
		Short:                false,
		DebugSave:            true,
		HTML:                 true,
	}

	options.Clean()

	err = Parse(r, options)
	if err != nil {
		log.Fatal(err)
	}

}

func TestForwardsGmail(ggggg *testing.T) {
	function, _, _, _ := runtime.Caller(0)
	fn := runtime.FuncForPC(function).Name()
	fn = fn[strings.LastIndex(fn, ".Test")+5:]
	fn = strings.Join(strings.Split(fn, "_"), ": ")
	fmt.Printf("\033[1;32m%s\033[0m\n", fn)

	ggggg.Skip("ручной тест")

	// return
	a := assert.New(ggggg)
	_ = a

	file1 := "../mail/gmail.eml"

	r, err := email.Parse(file.OpenE(file1))
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("to", r.To.Address)
	fmt.Println("forward", r.Forward)

	options := &Options{
		ShortLinks: func(s ...string) (res map[string]string) {
			res = make(map[string]string)
			for _, x := range s {
				res[x] = "testlink"
			}
			return
		},
		ShortLink:        func(s string) string { return "testlink" },
		ShortLinksDomain: "go.senthy.com",
		//Proxy:                "http://144.217.197.151:38611",
		Debug:                true,
		MinImageSizeBytes:    10_000,
		ShortLinksMustScheme: true,
		Images:               true,
		Short:                false,
		DebugSave:            true,
		HTML:                 true,
	}

	options.Clean()

	err = Parse(r, options)
	if err != nil {
		log.Fatal(err)
	}

}

//badhtml

// senthy links client
const host = "https://go.senthy.com"
const token = "a7f8s6f8s6"

func SenthyShortLinks(links ...string) (res map[string]string) {

	var shortlinkhost = fmt.Sprintf("%s/%s", host, token)

	res = make(map[string]string)
	b, _ := json.Marshal(links)

	r := file.Post(shortlinkhost, b)

	json.Unmarshal(r.Body(), &res)
	return
}

// client
func SenthyShortLink(link string) string {
	var shortlinkhost = fmt.Sprintf("%s/%s/%s", host, token, link)
	r, _ := file.Get(shortlinkhost)
	return string(r)
}
