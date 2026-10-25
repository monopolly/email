package parse

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/monopolly/email"
)

func TestShortTextLinks(t *testing.T) {
	long := "https://example.com/very/long/link"
	longer := long + "/and/more"
	text := long + " " + longer
	shorts := map[string]string{long: "https://go.senthy.com/1", longer: "go.senthy.com/2", "https://bad": "::"}
	opt := &Options{ShortLinks: func(...string) map[string]string { return shorts }}

	assert.Equal(t, "go.senthy.com/1 go.senthy.com/2", shortTextLinks(text, nil, opt))

	opt.ShortLinksMustScheme = true
	assert.Equal(t, "https://go.senthy.com/1 https://go.senthy.com/2", shortTextLinks(text, nil, opt))
}

func TestParseOptions(t *testing.T) {
	// nil options не паникуют
	r := &email.Email{Text: "hello"}
	require.NoError(t, Parse(r, nil))

	// переданные options не меняются
	opt := &Options{Short: true}
	r = &email.Email{HTML: "<p>hello world</p>"}
	require.NoError(t, Parse(r, opt))
	assert.Zero(t, opt.MaxTextLen)
	assert.Zero(t, opt.MaxWords)
}
