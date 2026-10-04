package icon

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
)

const (
	// userAgent identifies this application to the remote host.
	userAgent = "xyhelper-nav/1.0 (+icon-cache)"
	// maxRedirects limits how many redirects are followed.
	maxRedirects = 5
	// sniffLimit is how much of the payload is inspected for format detection.
	sniffLimit = 1024
)

// download fetches an icon and returns its bytes together with the file
// extension derived from the payload.
func (s *sIcon) download(ctx context.Context, source string) ([]byte, string, error) {
	maxSize := s.configSize(ctx, configMaxSize, defaultMaxSize)

	client := &http.Client{
		Timeout: s.configDuration(ctx, configTimeout, defaultTimeout),
		CheckRedirect: func(request *http.Request, via []*http.Request) error {
			if len(via) >= maxRedirects {
				return gerror.Newf("stopped after %d redirects", maxRedirects)
			}
			return nil
		},
	}
	request, err := http.NewRequestWithContext(ctx, http.MethodGet, source, nil)
	if err != nil {
		return nil, "", err
	}
	request.Header.Set("User-Agent", userAgent)
	request.Header.Set("Accept", "image/*,*/*;q=0.8")

	response, err := client.Do(request)
	if err != nil {
		return nil, "", err
	}
	defer func() {
		if closeErr := response.Body.Close(); closeErr != nil {
			g.Log().Debugf(ctx, "closing icon response failed: %v", closeErr)
		}
	}()

	if response.StatusCode != http.StatusOK {
		return nil, "", gerror.Newf("unexpected status %d", response.StatusCode)
	}
	// Read one byte more than allowed so that oversized payloads are detected.
	data, err := io.ReadAll(io.LimitReader(response.Body, maxSize+1))
	if err != nil {
		return nil, "", err
	}
	if int64(len(data)) > maxSize {
		return nil, "", gerror.Newf("icon is larger than the %d byte limit", maxSize)
	}
	if len(data) == 0 {
		return nil, "", gerror.New("icon response is empty")
	}
	extension, err := extensionFor(response.Header.Get("Content-Type"), data)
	if err != nil {
		return nil, "", err
	}
	return data, extension, nil
}

// extensionFor resolves the file extension of an icon, preferring the declared
// content type and falling back to inspecting the payload.
func extensionFor(contentType string, data []byte) (string, error) {
	switch mediaType(contentType) {
	case "image/png":
		return ".png", nil
	case "image/jpeg", "image/jpg", "image/pjpeg":
		return ".jpg", nil
	case "image/gif":
		return ".gif", nil
	case "image/webp":
		return ".webp", nil
	case "image/svg+xml":
		return ".svg", nil
	case "image/avif":
		return ".avif", nil
	case "image/x-icon", "image/vnd.microsoft.icon", "image/ico", "image/icon":
		return ".ico", nil
	}
	return sniffExtension(data)
}

// sniffExtension detects the icon format from the payload, which is needed
// because many servers answer with a generic content type.
func sniffExtension(data []byte) (string, error) {
	switch http.DetectContentType(data) {
	case "image/png":
		return ".png", nil
	case "image/jpeg":
		return ".jpg", nil
	case "image/gif":
		return ".gif", nil
	case "image/webp":
		return ".webp", nil
	}
	if bytes.Contains(bytes.ToLower(trimTo(data, sniffLimit)), []byte("<svg")) {
		return ".svg", nil
	}
	// ICO files start with a reserved zero word followed by the type word 1.
	if len(data) >= 4 && data[0] == 0x00 && data[1] == 0x00 && data[2] == 0x01 && data[3] == 0x00 {
		return ".ico", nil
	}
	return "", gerror.Newf("unsupported icon content type %q", http.DetectContentType(data))
}

// mediaType strips the parameters of a content type header.
func mediaType(contentType string) string {
	return strings.ToLower(strings.TrimSpace(strings.Split(contentType, ";")[0]))
}

// trimTo caps a byte slice to the first limit bytes.
func trimTo(data []byte, limit int) []byte {
	if len(data) > limit {
		return data[:limit]
	}
	return data
}

// cacheKeyFor derives a content addressed file name: identical icons share a
// single file and the address can be cached forever by the browser.
func cacheKeyFor(data []byte) string {
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:8])
}

// store writes an icon into the storage directory. The write goes through a
// temporary file that is renamed into place, so a partially written icon can
// never be served and concurrent downloads of the same icon are harmless.
func store(ctx context.Context, dir, key string, data []byte) error {
	if err := gfile.Mkdir(dir); err != nil {
		return err
	}
	target := gfile.Join(dir, key)
	if gfile.Exists(target) {
		return nil
	}
	temp, err := os.CreateTemp(dir, ".icon-*")
	if err != nil {
		return err
	}
	name := temp.Name()
	cleanup := func() {
		if removeErr := os.Remove(name); removeErr != nil && !os.IsNotExist(removeErr) {
			g.Log().Debugf(ctx, "removing temporary icon %s failed: %v", filepath.Base(name), removeErr)
		}
	}
	if _, err = temp.Write(data); err != nil {
		if closeErr := temp.Close(); closeErr != nil {
			g.Log().Debugf(ctx, "closing temporary icon failed: %v", closeErr)
		}
		cleanup()
		return err
	}
	if err = temp.Chmod(0o644); err != nil {
		if closeErr := temp.Close(); closeErr != nil {
			g.Log().Debugf(ctx, "closing temporary icon failed: %v", closeErr)
		}
		cleanup()
		return err
	}
	if err = temp.Close(); err != nil {
		cleanup()
		return err
	}
	if err = os.Rename(name, target); err != nil {
		cleanup()
		return err
	}
	return nil
}
