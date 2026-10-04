package icon

import "testing"

func TestExtensionFor(t *testing.T) {
	png := []byte{0x89, 'P', 'N', 'G', 0x0d, 0x0a, 0x1a, 0x0a, 0, 0, 0, 0}
	gif := []byte("GIF89a\x01\x00\x01\x00")
	svg := []byte(`<?xml version="1.0"?><svg xmlns="http://www.w3.org/2000/svg"></svg>`)
	ico := []byte{0x00, 0x00, 0x01, 0x00, 0x01, 0x00}
	html := []byte("<html><body>not an image</body></html>")

	cases := []struct {
		name        string
		contentType string
		data        []byte
		want        string
		wantErr     bool
	}{
		{name: "header wins for png", contentType: "image/png", data: png, want: ".png"},
		{name: "jpeg aliases", contentType: "image/jpg", data: png, want: ".jpg"},
		{name: "jpeg with charset", contentType: "image/jpeg; charset=binary", data: png, want: ".jpg"},
		{name: "svg", contentType: "image/svg+xml", data: svg, want: ".svg"},
		{name: "ico via content type", contentType: "image/x-icon", data: ico, want: ".ico"},
		{name: "webp", contentType: "image/webp", data: png, want: ".webp"},
		{name: "sniffs a generic content type", contentType: "application/octet-stream", data: png, want: ".png"},
		{name: "sniffs gif", contentType: "", data: gif, want: ".gif"},
		{name: "sniffs svg", contentType: "text/plain", data: svg, want: ".svg"},
		{name: "sniffs ico magic bytes", contentType: "application/octet-stream", data: ico, want: ".ico"},
		{name: "rejects html", contentType: "text/html", data: html, wantErr: true},
		{name: "rejects an unknown binary payload", contentType: "application/octet-stream", data: []byte{0x01, 0x02, 0x03, 0x04}, wantErr: true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := extensionFor(c.contentType, c.data)
			if c.wantErr {
				if err == nil {
					t.Fatalf("extensionFor(%q) = %q, want an error", c.contentType, got)
				}
				return
			}
			if err != nil {
				t.Fatalf("extensionFor(%q) returned an unexpected error: %v", c.contentType, err)
			}
			if got != c.want {
				t.Fatalf("extensionFor(%q) = %q, want %q", c.contentType, got, c.want)
			}
		})
	}
}

func TestCacheKeyForIsContentAddressed(t *testing.T) {
	first := cacheKeyFor([]byte("icon-a"))
	if first != cacheKeyFor([]byte("icon-a")) {
		t.Fatal("the same payload must produce the same cache key")
	}
	if first == cacheKeyFor([]byte("icon-b")) {
		t.Fatal("different payloads must produce different cache keys")
	}
	if len(first) != 16 {
		t.Fatalf("cache key %q has %d characters, want 16", first, len(first))
	}
}

func TestFilePathRejectsUnsafeKeys(t *testing.T) {
	service := New()
	for _, key := range []string{"", "../../etc/passwd", "..%2fpasswd", "abcdef.png", "/etc/passwd", "0123456789abcdef.sh"} {
		if path := service.FilePath(nil, key); path != "" {
			t.Errorf("FilePath(%q) = %q, want an empty string", key, path)
		}
	}
}

func TestSplitUrl(t *testing.T) {
	cases := []struct {
		raw        string
		wantScheme string
		wantHost   string
		wantOk     bool
	}{
		{raw: "https://github.com/docs", wantScheme: "https", wantHost: "github.com", wantOk: true},
		{raw: "HTTP://Example.com", wantScheme: "http", wantHost: "Example.com", wantOk: true},
		{raw: "example.com", wantOk: false},
		{raw: "", wantOk: false},
	}
	for _, c := range cases {
		scheme, host, ok := splitUrl(c.raw)
		if ok != c.wantOk || scheme != c.wantScheme || host != c.wantHost {
			t.Errorf("splitUrl(%q) = (%q, %q, %v), want (%q, %q, %v)",
				c.raw, scheme, host, ok, c.wantScheme, c.wantHost, c.wantOk)
		}
	}
}
