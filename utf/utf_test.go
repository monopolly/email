package utf

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
	xunicode "golang.org/x/text/encoding/unicode"
)

func encodeTo(t *testing.T, e encoding.Encoding, s string) string {
	t.Helper()
	b, err := e.NewEncoder().String(s)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestToUTF8(t *testing.T) {
	const ru = "Привет, как дела? Съешь же ещё этих мягких французских булок"
	const short = "Тема: счёт"

	cases := []struct {
		name    string
		in      string
		charset []string
		want    string
	}{
		{"utf8 as is", ru, nil, ru},
		{"ascii as is", "Hello, World", nil, "Hello, World"},
		{"german stays", "Größe und Maße für Österreich", nil, "Größe und Maße für Österreich"},
		{"french stays", "Élève à l'école, déjà vu", nil, "Élève à l'école, déjà vu"},
		{"brand case", "Мы вКонтакте и в iPhone", nil, "Мы вКонтакте и в iPhone"},
		{"windows-1251", encodeTo(t, charmap.Windows1251, ru), nil, ru},
		{"windows-1251 short", encodeTo(t, charmap.Windows1251, short), nil, short},
		{"koi8-r", encodeTo(t, charmap.KOI8R, ru), nil, ru},
		{"koi8-r short", encodeTo(t, charmap.KOI8R, short), nil, short},
		{"ibm866", encodeTo(t, charmap.CodePage866, ru), nil, ru},
		{"iso-8859-5", encodeTo(t, charmap.ISO8859_5, ru), nil, ru},
		{"declared charset", encodeTo(t, charmap.KOI8R, short), []string{"koi8-r"}, short},
		{"wrong declared charset", encodeTo(t, charmap.Windows1251, ru), []string{"koi8-r"}, ru},
		{"latin1", encodeTo(t, charmap.Windows1252, "Größe und Maße"), nil, "Größe und Maße"},
		{"utf-16 bom", encodeTo(t, xunicode.UTF16(xunicode.LittleEndian, xunicode.UseBOM), ru), nil, ru},
		{"utf-8 bom", "\xEF\xBB\xBF" + ru, nil, ru},
		// валидный UTF-8, но с кракозябрами
		{"1251 read as 1252", "Ïðèâåò, êàê äåëà?", nil, "Привет, как дела?"},
		{"utf8 read as 1251", "РџСЂРёРІРµС‚, РєР°Рє РґРµР»Р°?", nil, "Привет, как дела?"},
		{"utf8 read as latin1", "CafÃ© crÃ¨me", nil, "Café crème"},
		// NFD -> NFC: «й» из «и» + U+0306
		{"nfc", "йод", nil, "йод"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			assert.Equal(t, c.want, Fix(c.in, c.charset...))
		})
	}
}

func TestToUTF8HTMLMeta(t *testing.T) {
	html := `<html><head><meta http-equiv="Content-Type" content="text/html; charset=koi8-r"></head><body>` +
		encodeTo(t, charmap.KOI8R, "Добрый день") + `</body></html>`
	assert.Contains(t, Fix(html), "Добрый день")
}
