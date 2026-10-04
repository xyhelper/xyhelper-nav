package site

import "testing"

func TestNormalizeUrl(t *testing.T) {
	cases := []struct {
		name string
		raw  string
		want string
	}{
		{name: "keeps https", raw: "https://github.com", want: "https://github.com"},
		{name: "keeps http", raw: "http://example.com", want: "http://example.com"},
		{name: "adds scheme to a bare domain", raw: "example.com", want: "https://example.com"},
		{name: "adds scheme to a path", raw: "example.com/docs", want: "https://example.com/docs"},
		{name: "keeps protocol relative", raw: "//example.com", want: "//example.com"},
		{name: "keeps mailto", raw: "mailto:hi@example.com", want: "mailto:hi@example.com"},
		{name: "trims whitespace", raw: "   example.com  ", want: "https://example.com"},
		{name: "empty stays empty", raw: "", want: ""},
		{name: "scheme check is case insensitive", raw: "HTTPS://example.com", want: "HTTPS://example.com"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := NormalizeUrl(c.raw); got != c.want {
				t.Fatalf("NormalizeUrl(%q) = %q, want %q", c.raw, got, c.want)
			}
		})
	}
}
