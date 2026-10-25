package email

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
)

func TestDownload(ggggg *testing.T) {
	function, _, _, _ := runtime.Caller(0)
	fn := runtime.FuncForPC(function).Name()
	fn = fn[strings.LastIndex(fn, ".Test")+5:]
	fn = strings.Join(strings.Split(fn, "_"), ": ")
	fmt.Printf("\033[1;32m%s\033[0m\n", fn)

	return
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
func TestParseHTML(ggggg *testing.T) {
	function, _, _, _ := runtime.Caller(0)
	fn := runtime.FuncForPC(function).Name()
	fn = fn[strings.LastIndex(fn, ".Test")+5:]
	fn = strings.Join(strings.Split(fn, "_"), ": ")
	fmt.Printf("\033[1;32m%s\033[0m\n", fn)

	return
	// a, _ := utf8.DecodeLastRuneInString("а")
	// fmt.Println(unicode.IsLower(a))
	// fmt.Println(unicode.IsLower('а'))
	// fmt.Println(unicode.IsLower('.'))
	// fmt.Println(unicode.IsLower('к'))
	// fmt.Println(unicode.IsLower('К'))
	// select {}

	text := `Анна, здравствуйте!
	С вами Светлана Белова.
	Сегодня вы узнаете как правильно
	пользоваться нумерологией и
	научитесь
	делать свой первый расчет.
	Отличие нумерологии от других систем
	том, что для того, чтобы получить
	описание,
	вам нужна только дата рождения
	человека.`

	fmt.Println(retext(text))

	// return
	// r, _ := ParseRawEmail(file.OpenE("mail/love.eml"))
	// p, _ := goquery.NewDocumentFromReader(bytes.NewBuffer([]byte(r.HTML)))
	// e := p.Find("body")
	// file.Save("mail/love.html.txt", []byte(e.Text()))

	select {}
}
func TestMarkdownLinksText(ggggg *testing.T) {
	return
	//EXCLUSIVE DEAL ( go.senthy.com/3b8aa233 )
	//( go.senthy.com/3b8aa233 )
	line := "Редактировать ( https://follow.it/c6gav4?action=editFeedSettings ) Отписаться ( https://api.follow.it/unfollow/v2/AmaWpaPd1B7aGwbVRjkGe5TI3rpcT4N1ef3J7xQ2_fQ )"
	line = "EXCLUSIVE DEAL ( https://go.senthy.com/3b8aa233 )"
	line = "🔊︱community-announcement ( https://go.senthy.com/163d404c ) (KaratDAO)"
	line = "Новинки ( http://archiv.fagms.de/c/r?EMID=09E03MVKUHN1J0E40O30U0EODNH01BIRI58&RED=https%3A%2F%2Fwww.lamoda.ru%2Fc%2F4153%2Fdefault-women%2F%3Fis_new%3D1%26is_sale%3D0%26utm_content%3Dmenu_1%26utm_source%3Dnl%26utm_medium%3Dem%26utm_campaign%3D20230719_083425_LAM_NL_RU_SURF_HISTORY%26utm_term%3D1042522732%26adjust_campaign%3D20230719_083425_LAM_NL_RU_SURF_HISTORY%26adjust_t%3Dilg2wn6_a45rsrx ) Premium ( http://archiv.fagms.de/c/r?EMID=09E03MVKUHN1L0E40O30U0EODNH01PIFNID&RED=https%3A%2F%2Fwww.lamoda.ru%2Fhome-premium%2F%3Futm_content%3Dmenu_2%26utm_source%3Dnl%26utm_medium%3Dem%26utm_campaign%3D20230719_083425_LAM_NL_RU_SURF_HISTORY%26utm_term%3D1042522732%26adjust_campaign%3D20230719_083425_LAM_NL_RU_SURF_HISTORY%26adjust_t%3Dilg2wn6_a45rsrx ) Женщинам ( http://archiv.fagms.de/c/r?EMID=09E03MVKUHN1O0E40O30U0EODNH01U33A5J&RED=https%3A%2F%2Fwww.lamoda.ru%2Fwomen-home%2F%3Futm_content%3Dmenu_3%26utm_source%3Dnl%26utm_medium%3Dem%26utm_campaign%3D20230719_083425_LAM_NL_RU_SURF_HISTORY%26utm_term%3D1042522732%26adjust_campaign%3D20230719_083425_LAM_NL_RU_SURF_HISTORY%26adjust_t%3Dilg2wn6_a45rsrx ) Мужчинам ( http://archiv.fagms.de/c/r?EMID=09E03MVKUHN1N0E40O30U0EODNH0074R95E&RED=https%3A%2F%2Fwww.lamoda.ru%2Fmen-home%2F%3Futm_content%3Dmenu_4%26utm_source%3Dnl%26utm_medium%3Dem%26utm_campaign%3D20230719_083425_LAM_NL_RU_SURF_HISTORY%26utm_term%3D1042522732%26adjust_campaign%3D20230719_083425_LAM_NL_RU_SURF_HISTORY%26adjust_t%3Dilg2wn6_a45rsrx ) Детям ( http://archiv.fagms.de/c/r?EMID=09E03MVKUHN1M0E40O30U0EODNH00G2A9I3&RED=https%3A%2F%2Fwww.lamoda.ru%2Fkids-home%2F%3Futm_content%3Dmenu_5%26utm_source%3Dnl%26utm_medium%3Dem%26utm_campaign%3D20230719_083425_LAM_NL_RU_SURF_HISTORY%26utm_term%3D1042522732%26adjust_campaign%3D20230719_083425_LAM_NL_RU_SURF_HISTORY%26adjust_t%3Dilg2wn6_a45rsrx ) Распродажа % ( http://archiv.fagms.de/c/r?EMID=09E03MVKUHN1K0E40O30U0EODNH012IGNEA&RED=https%3A%2F%2Fwww.lamoda.ru%2Fc%2F4153%2Fdefault-women%2F%3Fis_sale%3D1%26display_locations%3Doutlet%26utm_content%3Dmenu_6%26utm_source%3Dnl%26utm_medium%3Dem%26utm_campaign%3D20230719_083425_LAM_NL_RU_SURF_HISTORY%26utm_term%3D1042522732%26adjust_campaign%3D20230719_083425_LAM_NL_RU_SURF_HISTORY%26adjust_t%3Dilg2wn6_a45rsrx )"
	fmt.Println("original:", line)
	fmt.Println()
	for _, x := range parseMarkdownLinks(line) {
		fmt.Println(x.title, x.link)
		fmt.Println()
	}

	fmt.Println("fix:", fixMarkdownLinks(line))

	//fmt.Println(string(deleteTag([]byte(line), "(", ")")))
}

func TestParseText(ggggg *testing.T) {
	function, _, _, _ := runtime.Caller(0)
	fn := runtime.FuncForPC(function).Name()
	fn = fn[strings.LastIndex(fn, ".Test")+5:]
	fn = strings.Join(strings.Split(fn, "_"), ": ")
	fmt.Printf("\033[1;32m%s\033[0m\n", fn)

	return
	a := assert.New(ggggg)
	_ = a

	file1 := []string{"mail/imagelinks.eml"}
	_ = file1

	file2 := file.FileList("mail", "eml")
	_ = file2

	for _, x := range file1 { //  {
		x = strings.TrimPrefix(x, "mail/")
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

		r, err := ParseRawEmail(file.OpenE(fmt.Sprintf("mail/%s.eml", name)))
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

		file.Save(fmt.Sprintf("mail/%s.html", name), []byte(r.HTML))
		file.Save(fmt.Sprintf("mail/%s.txt", name), []byte(r.Text))
		file.Save(fmt.Sprintf("mail/%s.md", name), []byte(r.Markdown))

		for _, fn := range r.Images {
			if fn.Size > options.MinImageSizeBytes {
				ff := fmt.Sprintf("mail/%s_%s", name, fn.Name)
				options.Log("save image", fn.Name)
				file.Save(ff, fn.Body)
			} else {
				options.Log("image size too small")
			}
		}

		for _, fn := range r.Files {
			if fn.Size > options.MinImageSizeBytes {
				ff := fmt.Sprintf("mail/%s_%s", name, fn.Name)
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

	return
	a := assert.New(ggggg)
	_ = a

	file1 := "mail/unsubscribeLink.eml"
	// file1 := "mail/unsubscribeEmail.eml"

	r, _ := ParseRawEmail(file.OpenE(file1))
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

	return

	file1 := "mail/redirect.eml"
	// file1 := "mail/unsubscribeEmail.eml"

	r, err := ParseRawEmail(file.OpenE(file1))
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

	return

	// return
	a := assert.New(ggggg)
	_ = a

	file1 := "mail/noread.eml"
	// file1 := "mail/unsubscribeEmail.eml"

	r, err := ParseRawEmail(file.OpenE(file1))
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

	// file1 := "mail/rawmail_error.eml"
	file1 := "mail/raw_error2.eml"
	// file1 := "mail/unsubscribeEmail.eml"

	r, err := ParseRawEmail(file.OpenE(file1))
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

	return

	// return
	a := assert.New(ggggg)
	_ = a

	file1 := "mail/gmail.eml"

	r, err := ParseRawEmail(file.OpenE(file1))
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
