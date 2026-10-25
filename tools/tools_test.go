package tools

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCleanTextKeepsBrackets(t *testing.T) {
	got := CleanText("[link](https://a.com) a|b")
	assert.Equal(t, "[link](https://a.com) a|b", got)
	// брайль, смайлы и комбинируемые знаки по-прежнему удаляются
	assert.Equal(t, "ab", CleanText("a⠀b😀"))
	// NFD «й» не превращается в «и»
	assert.Equal(t, "йод", CleanText("йод"))
	// пустые строки схлопываются в одну, а не пропадают после первой
	assert.Equal(t, "a\n\nb\n\nc", CleanText("a\n\n\nb\n\n\nc"))
}

func TestRemoveRepeatedEmptyLines(t *testing.T) {
	in := "a\n\n\n\nb\nc\n\nd"
	assert.Equal(t, "a\nb\nc\nd", RemoveRepeatedEmptyLines(in, 0))
	assert.Equal(t, "a\n\nb\nc\n\nd", RemoveRepeatedEmptyLines(in, 1))
	assert.Equal(t, "a\n\n\nb\nc\n\nd", RemoveRepeatedEmptyLines(in, 2))
}

func TestImageLinksDedupe(t *testing.T) {
	html := `<img src="https://a.com/1.png"><img src="https://a.com/1.png"><a href="https://a.com/1.png">x</a>`
	assert.Equal(t, []string{"https://a.com/1.png"}, ImageLinks(html))
}

func TestRetext(t *testing.T) {
	got := Retext("Ссылка https://a.com/some_path и почта john_doe@mail.ru, _курсив_")
	assert.Contains(t, got, "https://a.com/some_path")
	assert.Contains(t, got, "john_doe@mail.ru")
	assert.NotContains(t, got, "_курсив")

	assert.Equal(t, "Шаги:\n1. открыть\n2. нажать\n3. готово", Retext("Шаги: 1. открыть 2. нажать 3. готово"))
	assert.Equal(t, "Вышла версия 1. Подробности на сайте", Retext("Вышла версия 1. Подробности на сайте"))
}

func TestShortText(t *testing.T) {
	assert.Equal(t, "one two...", ShortText("one two three four\nfive", 2))
	assert.Equal(t, "a b\nc...", ShortText("a b\nc\nd e f g", 4))
	assert.Equal(t, "a b...", ShortText("a b\n\nc d e f g", 4))
}

func TestHasRepeats(t *testing.T) {
	assert.True(t, HasRepeats("ааа"))
	assert.False(t, HasRepeats("аба"))
	assert.True(t, HasRepeats("---"))
}

func TestParentHost(t *testing.T) {
	assert.Equal(t, "apple.com", parentHost("https://api.apple.com/x"))
	assert.Equal(t, "bbc.co.uk", parentHost("https://news.bbc.co.uk/x"))
	assert.Equal(t, "127.0.0.1", parentHost("http://127.0.0.1:8080/x"))
}

func TestHTMLFix(t *testing.T) {
	assert.Equal(t,
		`<html lang="ru"><head><meta charset="utf-8"></head><body>привет, это текст</body></html>`,
		FixHTMLCharset(`<html lang="ru"><body>привет, это текст</body></html>`))
	assert.Equal(t,
		`<html><head><meta http-equiv="Content-Type" content="text/html; charset=utf-8"></head></html>`,
		FixHTMLCharset(`<html><head><meta http-equiv="Content-Type" content="text/html; charset=windows-1251"></head></html>`))
	assert.Equal(t,
		`<html><head id="h"><meta charset="utf-8"><title>t</title></head></html>`,
		FixHTMLCharset(`<html><head id="h"><title>t</title></head></html>`))
}

func TestParseUnsubscribeLink(t *testing.T) {
	assert.Nil(t, ParseUnsubscribeLink(""))
	assert.Equal(t, "https://a.com/u", ParseUnsubscribeLink("<mailto:u@a.com>, <https://a.com/u>").String())
	assert.Equal(t, "mailto:u@a.com", ParseUnsubscribeLink("<mailto:u@a.com>").String())
}
