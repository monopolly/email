package parse

import (
	"flag"
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
	"unicode/utf8"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/monopolly/email"
)

// go test -run TestMail -update — перезаписать mail/golden/*.golden
var update = flag.Bool("update", false, "перезаписать golden-файлы")

// опции без сети: картинки не качаются, короткие ссылки детерминированные
func mailTestOptions() *Options {
	return &Options{
		ShortLinks: func(links ...string) map[string]string {
			res := make(map[string]string)
			for _, x := range links {
				h := fnv.New32a()
				h.Write([]byte(x))
				res[x] = fmt.Sprintf("https://go.senthy.com/%08x", h.Sum32())
			}
			return res
		},
	}
}

func TestMail(t *testing.T) {
	list, err := filepath.Glob("../mail/*.eml")
	require.NoError(t, err)
	require.NotEmpty(t, list)

	for _, path := range list {
		name := strings.TrimSuffix(filepath.Base(path), ".eml")
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			raw, err := os.ReadFile(path)
			require.NoError(t, err)

			r, err := email.Parse(raw)
			require.NoError(t, err)
			require.NoError(t, Parse(r, mailTestOptions()))

			checkMailInvariants(t, r)

			got := dumpMail(r)
			golden := filepath.Join("..", "mail", "golden", name+".golden")
			if *update {
				require.NoError(t, os.MkdirAll(filepath.Dir(golden), 0o755))
				require.NoError(t, os.WriteFile(golden, []byte(got), 0o644))
				return
			}
			want, err := os.ReadFile(golden)
			require.NoError(t, err, "нет golden-файла, запустите: go test -run TestMail -update")
			assert.Equal(t, string(want), got)
		})
	}
}

var (
	htmlTagRe   = regexp.MustCompile(`(?i)</?(div|span|p|br|table|td|tr|img|a|style|head|body|html|font|center)\b[^>]*>`)
	messageIDRe = regexp.MustCompile(`^\d{8,}\.`)
)

// общие проверки, которые должны выполняться для любого письма
func checkMailInvariants(t *testing.T, r *email.Email) {
	t.Helper()

	for field, v := range map[string]string{
		"Subj": r.Subj, "From.Name": r.From.Name, "To.Name": r.To.Name,
		"Text": r.Text, "HTML": r.HTML, "Markdown": r.Markdown,
	} {
		assert.True(t, utf8.ValidString(v), "%s: невалидный UTF-8", field)
		assert.NotContains(t, v, "�", "%s: символ замены (битая кодировка)", field)
	}

	assert.Equal(t, strings.ToLower(r.From.Address), r.From.Address)
	assert.Equal(t, strings.ToLower(r.To.Address), r.To.Address)

	for _, x := range r.Forward {
		assert.Contains(t, x.Address, "@")
		assert.False(t, messageIDRe.MatchString(x.Address), "Forward похож на Message-Id: %s", x.Address)
		assert.NotEqual(t, r.From.Address, "", "")
	}

	for _, list := range [][]*email.File{r.Files, r.Inline} {
		for _, f := range list {
			assert.NotEmpty(t, f.Name, "файл без имени")
			assert.True(t, utf8.ValidString(f.Name), "имя файла: невалидный UTF-8: %q", f.Name)
			assert.Equal(t, len(f.Body), f.Size, "размер файла %s", f.Name)
		}
	}

	md := r.Markdown
	assert.False(t, htmlTagRe.MatchString(md), "html-теги в markdown: %q", htmlTagRe.FindString(md))
	assert.NotContains(t, md, "](//", "короткая ссылка с // в начале")
	assert.NotContains(t, md, "\n\n\n", "больше одной пустой строки подряд")
	assert.Equal(t, strings.TrimSpace(md), md, "пробелы по краям markdown")
	if strings.TrimSpace(r.Text) != "" {
		assert.NotEmpty(t, strings.TrimSpace(md), "есть тело письма, но пустой markdown")
	}
}

// текстовый снимок результата разбора для golden-файла
func dumpMail(r *email.Email) string {
	var b strings.Builder
	addr := func(name, address string) string {
		if name == "" {
			return address
		}
		return fmt.Sprintf("%s <%s>", name, address)
	}

	fmt.Fprintf(&b, "Subj: %s\n", r.Subj)
	fmt.Fprintf(&b, "From: %s\n", addr(r.From.Name, r.From.Address))
	fmt.Fprintf(&b, "To: %s\n", addr(r.To.Name, r.To.Address))
	for _, x := range r.CC {
		fmt.Fprintf(&b, "CC: %s\n", addr(x.Name, x.Address))
	}
	for _, x := range r.BCC {
		fmt.Fprintf(&b, "BCC: %s\n", addr(x.Name, x.Address))
	}
	for _, x := range r.Forward {
		fmt.Fprintf(&b, "Forward: %s\n", x.Address)
	}
	if r.Date != nil {
		fmt.Fprintf(&b, "Date: %s\n", r.Date.UTC().Format("2006-01-02 15:04:05"))
	}
	fmt.Fprintf(&b, "ID: %s\n", r.ID)
	if r.Unsubscribe != nil {
		fmt.Fprintf(&b, "Unsubscribe: %s\n", r.Unsubscribe)
	}
	fmt.Fprintf(&b, "ParserHTML: %v\n", r.ParserHTML)
	fmt.Fprintf(&b, "Text: %d bytes, HTML: %d bytes\n", len(r.Text), len(r.HTML))

	files := func(kind string, list []*email.File) {
		for _, f := range list {
			fmt.Fprintf(&b, "%s: %q %s %d bytes", kind, f.Name, f.Type, f.Size)
			if f.CID != "" {
				fmt.Fprintf(&b, " cid=%s", f.CID)
			}
			b.WriteString("\n")
		}
	}
	files("File", r.Files)
	files("Inline", r.Inline)
	files("Image", r.Images)

	b.WriteString("\n--- markdown ---\n")
	b.WriteString(r.Markdown)
	b.WriteString("\n")
	return b.String()
}
