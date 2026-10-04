package nav

import (
	"net/http"
	"strings"

	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"

	"xyhelper-nav/internal/middleware"
	"xyhelper-nav/internal/service"
)

// iconCacheControl tells browsers that a cached icon may be reused for a week.
// The cache key is part of the address, so an updated icon always has a new
// URL and no stale image can be served.
const iconCacheControl = "public, max-age=604800"

// IconHandler serves the locally cached icon of a site and downloads it on
// first use. It is bound outside the JSON API group because it returns binary
// content instead of the usual response envelope.
func IconHandler(r *ghttp.Request) {
	ctx := r.Context()
	siteId := r.Get("id").Int()
	if siteId <= 0 {
		r.Response.WriteStatus(http.StatusNotFound)
		return
	}
	site, err := service.Site().Get(ctx, siteId)
	if err != nil {
		r.Response.WriteStatus(http.StatusNotFound)
		return
	}
	// A login-protected site keeps its icon private as well.
	visible, err := service.Site().IsVisible(ctx, site, middleware.IsAdmin(ctx))
	if err != nil || !visible {
		r.Response.WriteStatus(http.StatusNotFound)
		return
	}
	key, err := service.Icon().Ensure(ctx, site)
	if err != nil || key == "" {
		writeIconMiss(r)
		return
	}
	path := service.Icon().FilePath(ctx, key)
	if path == "" || !gfile.Exists(path) {
		writeIconMiss(r)
		return
	}
	r.Response.Header().Set("Cache-Control", iconCacheControl)
	r.Response.Header().Set("X-Content-Type-Options", "nosniff")
	if strings.HasSuffix(key, ".svg") {
		// Icons are downloaded from remote hosts, never let a script inside an
		// SVG run in the context of this origin.
		r.Response.Header().Set("Content-Security-Policy", "default-src 'none'; style-src 'unsafe-inline'; sandbox")
	}
	r.Response.ServeFile(path)
}

// writeIconMiss answers with a non-cacheable 404 so the browser retries later
// and the front-end falls back to its generated letter avatar.
func writeIconMiss(r *ghttp.Request) {
	r.Response.Header().Set("Cache-Control", "no-store")
	r.Response.WriteStatus(http.StatusNotFound)
}
