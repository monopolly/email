package utf

import (
	"bytes"
	"regexp"
	"strings"
	"unicode"
	"unicode/utf8"

	"github.com/gogs/chardet"
	"golang.org/x/text/encoding"
	"golang.org/x/text/encoding/htmlindex"
	"golang.org/x/text/encoding/ianaindex"
	xunicode "golang.org/x/text/encoding/unicode"
	"golang.org/x/text/unicode/norm"
)

// ToUTF8 приводит текст в любой кодировке к корректному UTF-8 (NFC).
//
// charset — объявленная кодировка (из заголовка письма или <meta>), может быть
// пустой или неверной: кодировка всё равно проверяется по содержимому.
// Уже валидный UTF-8 проверяется на «кракозябры» (UTF-8, прочитанный как
// windows-1251, windows-1251, прочитанная как windows-1252/KOI8-R, и т.п.)
// и чинится, если перекодировка даёт явно более осмысленный текст.
func Fix(s string, charset ...string) string {
	if s == "" {
		return s
	}
	if utf8.ValidString(s) {
		s = strings.TrimPrefix(s, "\uFEFF")
		return norm.NFC.String(fixMojibake(s))
	}
	return norm.NFC.String(decodeBest([]byte(s), charset...))
}

// кодировки, которые пробуем всегда: кириллица и западноевропейская
var fallbackCharsets = []string{"windows-1251", "koi8-r", "koi8-u", "ibm866", "iso-8859-5", "x-mac-cyrillic", "windows-1252"}

// кодировки, в которых чаще всего ошибочно читают текст (источник кракозябр)
var mojibakeCharsets = []string{"windows-1252", "windows-1251", "koi8-r", "ibm866", "iso-8859-5"}

var metaCharset = regexp.MustCompile(`(?i)<meta[^>]+charset\s*=\s*["']?\s*([a-z0-9_:.\-]+)`)

// decodeBest перебирает возможные кодировки и выбирает ту,
// в которой текст выглядит наиболее осмысленно.
func decodeBest(b []byte, hints ...string) string {
	// BOM однозначно определяет кодировку
	switch {
	case bytes.HasPrefix(b, []byte{0xEF, 0xBB, 0xBF}):
		b = b[3:]
	case bytes.HasPrefix(b, []byte{0xFF, 0xFE}):
		if s, err := xunicode.UTF16(xunicode.LittleEndian, xunicode.ExpectBOM).NewDecoder().Bytes(b); err == nil {
			return string(s)
		}
	case bytes.HasPrefix(b, []byte{0xFE, 0xFF}):
		if s, err := xunicode.UTF16(xunicode.BigEndian, xunicode.ExpectBOM).NewDecoder().Bytes(b); err == nil {
			return string(s)
		}
	}

	var labels []string
	if utf8.Valid(b) {
		labels = append(labels, "utf-8")
	}
	labels = append(labels, hints...)

	// <meta charset> в html
	head := b
	if len(head) > 4096 {
		head = head[:4096]
	}
	if m := metaCharset.FindSubmatch(head); m != nil {
		labels = append(labels, string(m[1]))
	}

	if res, err := chardet.NewTextDetector().DetectAll(b); err == nil {
		for _, x := range res {
			if x.Confidence >= 10 {
				labels = append(labels, x.Charset)
			}
		}
	}
	labels = append(labels, fallbackCharsets...)

	best, bestScore := "", 0
	found := false
	seen := map[encoding.Encoding]bool{}
	for _, label := range labels {
		var text string
		if strings.EqualFold(label, "utf-8") || strings.EqualFold(label, "utf8") {
			if !utf8.Valid(b) {
				continue
			}
			text = string(b)
		} else {
			enc := lookupCharset(label)
			if enc == nil || seen[enc] {
				continue
			}
			seen[enc] = true
			out, err := enc.NewDecoder().Bytes(b)
			if err != nil {
				continue
			}
			text = string(out)
		}
		// при равенстве побеждает кодировка, стоящая раньше (объявленная, chardet)
		if score := textScore(text); !found || score > bestScore {
			best, bestScore, found = text, score, true
		}
	}
	if !found {
		return strings.ToValidUTF8(string(b), "�")
	}
	return best
}

func lookupCharset(label string) encoding.Encoding {
	label = strings.Trim(strings.TrimSpace(label), `"'`)
	if label == "" {
		return nil
	}
	if e, err := htmlindex.Get(label); err == nil && e != nil {
		return e
	}
	for _, idx := range []*ianaindex.Index{ianaindex.MIME, ianaindex.IANA} {
		if e, err := idx.Encoding(label); err == nil && e != nil {
			return e
		}
	}
	return nil
}

// fixMojibake чинит валидный UTF-8, который когда-то был прочитан не в той кодировке.
func fixMojibake(s string) string {
	if !looksLikeMojibake(s) {
		return s
	}
	best, bestScore := s, textScore(s)
	for _, label := range mojibakeCharsets {
		raw, err := lookupCharset(label).NewEncoder().String(s)
		if err != nil {
			// в тексте есть символы вне этой кодировки — значит, читали не в ней
			continue
		}
		text := decodeBest([]byte(raw))
		// небольшой запас, чтобы не трогать нормальный текст
		if score := textScore(text); score > bestScore+2 {
			best, bestScore = text, score
		}
	}
	return best
}

// looksLikeMojibake — быстрая проверка на типичные признаки кракозябр.
func looksLikeMojibake(s string) bool {
	var prev rune
	for _, r := range s {
		switch {
		case r >= 0x80 && r < 0xA0, r >= 0x2500 && r <= 0x259F:
			// C1-управляющие и псевдографика
			return true
		case isLatin1Letter(prev) && isLatin1Letter(r):
			// "Ïðèâåò" — windows-1251, прочитанная как windows-1252
			return true
		case isUTF8Lead(prev, r):
			// "CafÃ©" — UTF-8, прочитанный как latin1
			return true
		case unicode.Is(unicode.Cyrillic, prev) && unicode.Is(unicode.Cyrillic, r) &&
			unicode.IsLower(prev) && unicode.IsUpper(r):
			// "РџСЂРёРІРµС‚", "оПХБЕР" — UTF-8/KOI8-R, прочитанные как windows-1251
			return true
		}
		prev = r
	}
	return false
}

func isLatin1Letter(r rune) bool {
	return r >= 0xC0 && r <= 0xFF && r != 0xD7 && r != 0xF7
}

// isUTF8Lead — пара "Ã©": первый байт двухбайтового UTF-8 символа
// и знак из U+00A0..U+00BF, т.е. байты UTF-8, прочитанные как latin1
func isUTF8Lead(prev, r rune) bool {
	return prev >= 0xC2 && prev <= 0xDF && r >= 0xA0 && r <= 0xBF
}

// textScore — насколько текст похож на нормальный: соседние буквы одной
// письменности и без скачков регистра внутри слова дают плюс,
// мусорные символы и смесь алфавитов — минус.
func textScore(s string) (score int) {
	var prev rune
	for _, r := range s {
		switch {
		case r == utf8.RuneError:
			score -= 10
		case r >= 0x80 && r < 0xA0:
			score -= 5
		case unicode.IsControl(r) && r != '\n' && r != '\r' && r != '\t':
			score -= 5
		case r >= 0x2500 && r <= 0x259F, isUTF8Lead(prev, r):
			score -= 3
		}
		if unicode.IsLetter(prev) && unicode.IsLetter(r) {
			score += letterPairScore(prev, r)
		}
		prev = r
	}
	return
}

func letterPairScore(a, b rune) int {
	switch {
	case letterScript(a) != letterScript(b):
		return -3
	case unicode.IsLower(a) && unicode.IsUpper(b):
		return -3
	case isLatin1Letter(a) && isLatin1Letter(b):
		return -2
	}
	return 1
}

func letterScript(r rune) int {
	switch {
	case r < 0x80, unicode.Is(unicode.Latin, r):
		return 1
	case unicode.Is(unicode.Cyrillic, r):
		return 2
	case unicode.Is(unicode.Greek, r):
		return 3
	}
	return 4
}
