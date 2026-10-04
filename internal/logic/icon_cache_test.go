package logic_test

import (
	"context"
	"encoding/base64"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gogf/gf/v2/os/gfile"

	"xyhelper-nav/internal/consts"
	"xyhelper-nav/internal/model"
	"xyhelper-nav/internal/service"
)

// testIcon is a valid 1x1 PNG used as the payload of the fake icon host.
var testIcon = mustDecodeBase64(
	"iVBORw0KGgoAAAANSUhEUgAAAAEAAAABCAYAAAAfFcSJAAAADUlEQVR42mP8z8BQDwAEhQGAhKmMIQAAAABJRU5ErkJggg==",
)

// mustDecodeBase64 decodes a constant and panics on a programming error.
func mustDecodeBase64(value string) []byte {
	data, err := base64.StdEncoding.DecodeString(value)
	if err != nil {
		panic(err)
	}
	return data
}

// TestIconCacheDownloadsOnce proves that the remote host is contacted a single
// time per icon: later page views are served from the local cache, which is the
// whole point of caching icons inside the project.
func TestIconCacheDownloadsOnce(t *testing.T) {
	var hits int
	remote := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		w.Header().Set("Content-Type", "image/png")
		if _, err := w.Write(testIcon); err != nil {
			t.Errorf("writing the fake icon failed: %v", err)
		}
	}))
	defer remote.Close()

	var (
		ctx    = context.Background()
		status = consts.StatusOn
	)
	categoryId, err := service.Category().Create(ctx, &model.CategorySaveInput{
		Name:   "图标缓存测试",
		Status: &status,
	})
	if err != nil {
		t.Fatalf("create category: %v", err)
	}
	siteId, err := service.Site().Create(ctx, &model.SiteSaveInput{
		CategoryId: categoryId,
		Name:       "Icon Fixture",
		Url:        remote.URL,
		Icon:       remote.URL + "/favicon.ico",
		Status:     &status,
	})
	if err != nil {
		t.Fatalf("create site: %v", err)
	}
	defer func() {
		if err = service.Site().Delete(ctx, siteId); err != nil {
			t.Errorf("cleaning up the fixture failed: %v", err)
		}
	}()

	// First view: the icon is downloaded and stored locally.
	site, err := service.Site().Get(ctx, siteId)
	if err != nil {
		t.Fatalf("get site: %v", err)
	}
	key, err := service.Icon().Ensure(ctx, site)
	if err != nil {
		t.Fatalf("ensure icon: %v", err)
	}
	if key == "" {
		t.Fatal("icon was not cached")
	}
	if hits != 1 {
		t.Fatalf("remote host was contacted %d times, want 1", hits)
	}
	path := service.Icon().FilePath(ctx, key)
	if !gfile.Exists(path) {
		t.Fatalf("cached icon file %q does not exist", path)
	}

	// Later views reuse the cached file and never touch the remote host again.
	site, err = service.Site().Get(ctx, siteId)
	if err != nil {
		t.Fatalf("reload site: %v", err)
	}
	if key, err = service.Icon().Ensure(ctx, site); err != nil {
		t.Fatalf("ensure icon again: %v", err)
	}
	if key == "" {
		t.Fatal("cached icon disappeared on the second view")
	}
	if hits != 1 {
		t.Fatalf("remote host was contacted %d times after a second view, want 1", hits)
	}
}

// TestIconCacheBacksOffAfterFailure proves that an unreachable or broken icon
// host is not retried on every page view.
func TestIconCacheBacksOffAfterFailure(t *testing.T) {
	var hits int
	remote := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		hits++
		http.Error(w, "no icon here", http.StatusNotFound)
	}))
	defer remote.Close()

	var (
		ctx    = context.Background()
		status = consts.StatusOn
	)
	categoryId, err := service.Category().Create(ctx, &model.CategorySaveInput{
		Name:   "图标失败测试",
		Status: &status,
	})
	if err != nil {
		t.Fatalf("create category: %v", err)
	}
	siteId, err := service.Site().Create(ctx, &model.SiteSaveInput{
		CategoryId: categoryId,
		Name:       "Broken Icon Fixture",
		Url:        remote.URL,
		Status:     &status,
	})
	if err != nil {
		t.Fatalf("create site: %v", err)
	}
	defer func() {
		if err = service.Site().Delete(ctx, siteId); err != nil {
			t.Errorf("cleaning up the fixture failed: %v", err)
		}
	}()

	for attempt := 1; attempt <= 3; attempt++ {
		site, err := service.Site().Get(ctx, siteId)
		if err != nil {
			t.Fatalf("get site: %v", err)
		}
		key, err := service.Icon().Ensure(ctx, site)
		if err != nil {
			t.Fatalf("ensure icon: %v", err)
		}
		if key != "" {
			t.Fatalf("attempt %d cached an icon that does not exist", attempt)
		}
	}
	if hits != 1 {
		t.Fatalf("broken remote host was contacted %d times, want 1 (back-off)", hits)
	}
}
