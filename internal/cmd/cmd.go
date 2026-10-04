// Package cmd contains the entry commands of the application.
package cmd

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gfile"

	"xyhelper-nav/internal/boot"
	"xyhelper-nav/internal/consts"
	"xyhelper-nav/internal/controller/admin"
	"xyhelper-nav/internal/controller/nav"
	"xyhelper-nav/internal/middleware"

	// Register all business implementations.
	_ "xyhelper-nav/internal/logic"
)

// Main is the default command that starts the HTTP server.
var Main = gcmd.Command{
	Name:  "main",
	Usage: "main",
	Brief: "start the xyhelper-nav http server",
	Func: func(ctx context.Context, parser *gcmd.Parser) (err error) {
		if err = boot.Init(ctx); err != nil {
			return err
		}
		s := g.Server()
		serveStatic(ctx, s)
		// Site icons are served from the local cache, so the browser never has
		// to contact the remote sites itself.
		s.BindHandler("GET:"+consts.IconEndpointPrefix+"/{id}", nav.IconHandler)

		s.Group("/api/v1", func(group *ghttp.RouterGroup) {
			group.Middleware(ghttp.MiddlewareHandlerResponse)
			group.Middleware(middleware.CORS)

			// Public endpoints.
			group.Bind(
				nav.New(),
				admin.NewLogin(),
			)

			// Administration endpoints, protected by the session middleware.
			group.Group("/", func(adminGroup *ghttp.RouterGroup) {
				adminGroup.Middleware(middleware.Auth)
				adminGroup.Bind(admin.New())
			})
		})

		s.Run()
		return nil
	},
}

// serveStatic mounts the front-end assets and the admin entry page.
func serveStatic(ctx context.Context, s *ghttp.Server) {
	root := g.Cfg().MustGet(ctx, "server.serverRoot", "resource/public").String()
	if !gfile.Exists(root) {
		// Fall back to the directory of the executable, which is where a
		// packaged deployment usually keeps its resources.
		root = gfile.Join(gfile.SelfDir(), "resource/public")
	}
	if !gfile.Exists(root) {
		g.Log().Warningf(ctx, `static root %q not found, front-end pages are disabled`, root)
		return
	}
	s.SetServerRoot(root)
	// Friendly entries: /admin renders the administration panel.
	s.BindHandler("GET:/admin", func(r *ghttp.Request) {
		r.Response.ServeFile(gfile.Join(root, "admin.html"))
	})
}
