package email

import (
	"net/mail"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/jhillyerd/enmime/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/charmap"
)

func TestParseError(t *testing.T) {
	r, err := Parse(nil)
	if err == nil {
		// enmime разбирает пустой ввод без ошибки — главное, не nil, nil
		assert.NotNil(t, r)
		return
	}
	assert.Nil(t, r)
}

func TestParseHeaders(t *testing.T) {
	raw, err := os.ReadFile("mail/raw_error2.eml")
	require.NoError(t, err)
	r, err := Parse(raw)
	require.NoError(t, err)

	for _, x := range r.Forward {
		// Message-Id из In-Reply-To/References не адреса
		assert.NotContains(t, x.Address, "sip3.comfortel.pro")
		assert.NotContains(t, x.Address, "tickets.comfortel.pro")
	}
	require.NotNil(t, r.Date)
	assert.NotEmpty(t, r.ID)
}

func TestParse8bitHeaders(t *testing.T) {
	subj := encodeTo(t, charmap.Windows1251, "Счёт на оплату")
	body := encodeTo(t, charmap.KOI8R, "Добрый день! Высылаем счёт.")
	raw := "From: a@b.com\r\nTo: c@d.com\r\nSubject: " + subj + "\r\n" +
		"Content-Type: text/plain\r\n\r\n" + body + "\r\n"
	r, err := Parse([]byte(raw))
	require.NoError(t, err)
	assert.Equal(t, "Счёт на оплату", r.Subj)
	assert.Contains(t, r.Text, "Добрый день! Высылаем счёт.")
}

func TestRender(t *testing.T) {
	m := &Email{
		From: mail.Address{Address: "a@b.com"},
		To:   mail.Address{Address: "c@d.com"},
		Subj: "test",
		HTML: `<img src="cid:x">`,
	}
	m.AddInline("a.png", "image/png", []byte("png"))

	// без даты не паникует
	out := m.Render()
	require.NotZero(t, out.Len())

	// cid вложения совпадает с File.ID, на него можно сослаться из html
	env, err := enmime.ReadEnvelope(strings.NewReader(out.String()))
	require.NoError(t, err)
	require.Len(t, env.Inlines, 1)
	assert.Equal(t, m.Inline[0].ID, env.Inlines[0].ContentID)
}

func TestAddFilesDedupe(t *testing.T) {
	dir := t.TempDir()
	a := filepath.Join(dir, "a.txt")
	b := filepath.Join(dir, "b.txt")
	require.NoError(t, os.WriteFile(a, []byte("same"), 0o644))
	require.NoError(t, os.WriteFile(b, []byte("same"), 0o644))

	m := New().AddFiles(a, b).AddInlines(a, b).Done()
	assert.Len(t, m.Files, 1)
	assert.Len(t, m.Inline, 1)

	e := &Email{}
	e.AddImages(a, b)
	assert.Len(t, e.Images, 1)
}

func TestSubjPercent(t *testing.T) {
	assert.Equal(t, "скидка 50%", New().Subj("скидка 50%").Done().Subj)
	assert.Equal(t, "скидка 50%", New().Subj("скидка %d%%", 50).Done().Subj)
}

func encodeTo(t *testing.T, e encoding.Encoding, s string) string {
	t.Helper()
	b, err := e.NewEncoder().String(s)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func TestPackJSONKeys(t *testing.T) {
	m := &Email{ID: "<1@a>", Subj: "тема", HTML: "<p>x</p>", ParserHTML: true, From: mail.Address{Address: "a@b.com"}}
	b := m.Pack()
	for _, key := range []string{`"id"`, `"subj"`, `"html"`, `"parserhtml"`, `"from"`} {
		assert.Contains(t, string(b), key)
	}
	assert.NotContains(t, string(b), `"Subj"`)

	r, err := Unpack(b)
	require.NoError(t, err)
	assert.Equal(t, m.Subj, r.Subj)
	assert.Equal(t, m.HTML, r.HTML)

	// JSON в старом формате (ключи с заглавной) тоже читается
	old, err := Unpack([]byte(`{"ID":"<2@a>","Subj":"старая","ParserHTML":true}`))
	require.NoError(t, err)
	assert.Equal(t, "<2@a>", old.ID)
	assert.Equal(t, "старая", old.Subj)
	assert.True(t, old.ParserHTML)
}
