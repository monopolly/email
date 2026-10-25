package tools

import (
	"fmt"
	"runtime"
	"strings"
	"testing"
)

func TestParseHTML(ggggg *testing.T) {
	function, _, _, _ := runtime.Caller(0)
	fn := runtime.FuncForPC(function).Name()
	fn = fn[strings.LastIndex(fn, ".Test")+5:]
	fn = strings.Join(strings.Split(fn, "_"), ": ")
	fmt.Printf("\033[1;32m%s\033[0m\n", fn)

	ggggg.Skip("ручной тест")
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

	fmt.Println(Retext(text))

	// return
	// r, _ := ParseRawEmail(file.OpenE("mail/love.eml"))
	// p, _ := goquery.NewDocumentFromReader(bytes.NewBuffer([]byte(r.HTML)))
	// e := p.Find("body")
	// file.Save("mail/love.html.txt", []byte(e.Text()))

	select {}
}

func TestMarkdownLinksText(ggggg *testing.T) {
	ggggg.Skip("ручной тест")
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

	fmt.Println("fix:", FixMarkdownLinks(line))

	//fmt.Println(string(deleteTag([]byte(line), "(", ")")))
}
