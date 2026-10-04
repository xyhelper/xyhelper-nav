// Package icon caches remote site icons inside the project.
//
// The browser only ever requests icons from this server: the first request for
// an icon makes the server download it, store it under the configured icon
// directory and serve it from there. Every later request is served from disk,
// which removes the cross-site requests and the latency of the navigation page.
package icon

import (
	"context"
	"net/url"
	"regexp"
	"strings"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"

	"xyhelper-nav/internal/dao"
	"xyhelper-nav/internal/model"
	"xyhelper-nav/internal/model/entity"
	"xyhelper-nav/internal/service"
)

// Configuration keys and their defaults.
const (
	configPath            = "icon.path"
	configTimeout         = "icon.timeout"
	configMaxSize         = "icon.maxSize"
	configRetryInterval   = "icon.retryInterval"
	configRefreshInterval = "icon.refreshInterval"

	defaultPath            = "resource/data/icons"
	defaultTimeout         = "8s"
	defaultMaxSize         = 524288 // 512 KiB
	defaultRetryInterval   = "10m"
	defaultRefreshInterval = "720h" // 30 days
)

// cacheKeyPattern guards against path traversal: only file names produced by
// this package are accepted.
var cacheKeyPattern = regexp.MustCompile(`^[0-9a-f]{16}\.(png|jpg|gif|webp|svg|ico|avif)$`)

type sIcon struct{}

func init() {
	service.RegisterIcon(New())
}

// New creates and returns the icon cache service implementation.
func New() service.IIcon {
	return &sIcon{}
}

var siteColumns = dao.NavSite.Columns()

// Ensure makes sure the icon of a site is cached locally and returns its key.
// An empty key means no icon is available, and the front-end falls back to the
// generated letter avatar.
func (s *sIcon) Ensure(ctx context.Context, site *entity.NavSite) (string, error) {
	if site == nil || s.sourceUrl(site) == "" {
		return "", nil
	}
	dir := s.directory(ctx)
	cached := site.IconKey != "" && gfile.Exists(gfile.Join(dir, site.IconKey))

	if cached && !elapsed(site.IconFetchedAt, s.configDuration(ctx, configRefreshInterval, defaultRefreshInterval)) {
		return site.IconKey, nil
	}
	if !cached && !elapsed(site.IconFetchedAt, s.configDuration(ctx, configRetryInterval, defaultRetryInterval)) {
		// A previous attempt failed recently: back off instead of contacting
		// the remote host on every single page view.
		return "", nil
	}

	key, err := s.fetch(ctx, site)
	if err != nil {
		g.Log().Warningf(ctx, "icon of site %d could not be cached: %v", site.Id, err)
		if cached {
			// Keep serving the stale icon rather than showing nothing.
			return site.IconKey, nil
		}
		return "", nil
	}
	return key, nil
}

// Refresh re-downloads icons: a single site when siteId > 0, every site when
// siteId is 0.
func (s *sIcon) Refresh(ctx context.Context, siteId int) (*model.IconRefreshResult, error) {
	query := dao.NavSite.Ctx(ctx)
	if siteId > 0 {
		query = query.Where(siteColumns.Id, siteId)
	}
	var sites []*entity.NavSite
	if err := query.Scan(&sites); err != nil {
		return nil, err
	}
	result := &model.IconRefreshResult{}
	for _, site := range sites {
		if s.sourceUrl(site) == "" {
			continue
		}
		result.Total++
		key, err := s.fetch(ctx, site)
		if err != nil || key == "" {
			g.Log().Warningf(ctx, "icon refresh failed for site %d: %v", site.Id, err)
			result.Failed++
			continue
		}
		result.Refreshed++
	}
	return result, nil
}

// FilePath resolves a cache key to an absolute file path.
func (s *sIcon) FilePath(ctx context.Context, key string) string {
	if !cacheKeyPattern.MatchString(key) {
		return ""
	}
	return gfile.Join(s.directory(ctx), key)
}

// Release deletes a cached file once no site references it anymore.
func (s *sIcon) Release(ctx context.Context, key string) error {
	if !cacheKeyPattern.MatchString(key) {
		return nil
	}
	count, err := dao.NavSite.Ctx(ctx).Where(siteColumns.IconKey, key).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	return gfile.Remove(gfile.Join(s.directory(ctx), key))
}

// fetch downloads the icon of a site, stores it and records the cache key.
func (s *sIcon) fetch(ctx context.Context, site *entity.NavSite) (string, error) {
	source := s.sourceUrl(site)
	if source == "" {
		return "", nil
	}
	// Record the attempt first so that a failure also triggers the back-off.
	if err := s.markFetched(ctx, site.Id); err != nil {
		return "", err
	}
	data, extension, err := s.download(ctx, source)
	if err != nil {
		return "", err
	}
	key := cacheKeyFor(data) + extension
	if err = store(ctx, s.directory(ctx), key, data); err != nil {
		return "", err
	}
	if err = s.markCached(ctx, site.Id, key); err != nil {
		return "", err
	}
	if previous := site.IconKey; previous != "" && previous != key {
		// The icon changed: drop the old file when nothing references it.
		if err = s.Release(ctx, previous); err != nil {
			return "", err
		}
	}
	return key, nil
}

// markFetched records a download attempt.
func (s *sIcon) markFetched(ctx context.Context, siteId int) error {
	_, err := dao.NavSite.Ctx(ctx).
		Where(siteColumns.Id, siteId).
		Data(g.Map{siteColumns.IconFetchedAt: gtime.Now()}).
		Update()
	return err
}

// markCached stores the cache key of a successfully downloaded icon.
func (s *sIcon) markCached(ctx context.Context, siteId int, key string) error {
	_, err := dao.NavSite.Ctx(ctx).
		Where(siteColumns.Id, siteId).
		Data(g.Map{
			siteColumns.IconKey:       key,
			siteColumns.IconFetchedAt: gtime.Now(),
		}).
		Update()
	return err
}

// sourceUrl returns the address an icon is downloaded from. It uses the icon
// configured by the administrator and falls back to the favicon of the site.
func (s *sIcon) sourceUrl(site *entity.NavSite) string {
	if configured := strings.TrimSpace(site.Icon); configured != "" {
		if scheme, _, ok := splitUrl(configured); ok && (scheme == "http" || scheme == "https") {
			return configured
		}
		return ""
	}
	scheme, host, ok := splitUrl(strings.TrimSpace(site.Url))
	if !ok || (scheme != "http" && scheme != "https") || host == "" {
		return ""
	}
	return scheme + "://" + host + "/favicon.ico"
}

// directory returns the icon storage directory.
func (s *sIcon) directory(ctx context.Context) string {
	return s.configString(ctx, configPath, defaultPath)
}

// configString reads a string option from the configuration.
func (s *sIcon) configString(ctx context.Context, key, fallback string) string {
	if value, err := g.Cfg().Get(ctx, key); err == nil {
		if text := strings.TrimSpace(value.String()); text != "" {
			return text
		}
	}
	return fallback
}

// configDuration reads a duration option from the configuration.
func (s *sIcon) configDuration(ctx context.Context, key, fallback string) time.Duration {
	duration, err := gtime.ParseDuration(s.configString(ctx, key, fallback))
	if err != nil {
		duration, _ = gtime.ParseDuration(fallback)
	}
	return duration
}

// configSize reads an integer option from the configuration.
func (s *sIcon) configSize(ctx context.Context, key string, fallback int64) int64 {
	if value, err := g.Cfg().Get(ctx, key); err == nil {
		if size := value.Int64(); size > 0 {
			return size
		}
	}
	return fallback
}

// elapsed reports whether a timestamp is unset or older than the interval.
func elapsed(at time.Time, interval time.Duration) bool {
	if at.IsZero() {
		return true
	}
	return time.Since(at) >= interval
}

// splitUrl returns the lower-cased scheme and the host of an address.
func splitUrl(raw string) (scheme, host string, ok bool) {
	parsed, err := url.Parse(raw)
	if err != nil || parsed.Scheme == "" {
		return "", "", false
	}
	return strings.ToLower(parsed.Scheme), parsed.Host, true
}
