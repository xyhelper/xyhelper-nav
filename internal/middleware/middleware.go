// Package middleware contains the HTTP middlewares of the application.
package middleware

import (
	"context"
	"net/http"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"

	"xyhelper-nav/internal/consts"
)

// CurrentAdminId returns the id of the administrator behind the current
// request, or 0 when the request is anonymous. It is the single place where
// "is the visitor signed in?" is answered.
func CurrentAdminId(ctx context.Context) int {
	request := g.RequestFromCtx(ctx)
	if request == nil {
		return 0
	}
	// Set by the Auth middleware for routes it already protected.
	if adminId := request.GetCtxVar(consts.CtxAdminId).Int(); adminId > 0 {
		return adminId
	}
	adminId, err := request.Session.Get(consts.SessionAdminId)
	if err != nil || adminId == nil {
		return 0
	}
	return adminId.Int()
}

// IsAdmin reports whether the current request carries a valid administrator
// session. It is used by the public API to decide whether login-protected
// content may be disclosed.
func IsAdmin(ctx context.Context) bool {
	return CurrentAdminId(ctx) > 0
}

// Auth guards the administration API. It rejects every request that does not
// carry a valid administrator session.
func Auth(r *ghttp.Request) {
	adminId := CurrentAdminId(r.Context())
	if adminId <= 0 {
		r.Response.Status = http.StatusUnauthorized
		r.Response.WriteJson(g.Map{
			"code":    http.StatusUnauthorized,
			"message": "登录状态已失效，请重新登录",
		})
		r.ExitAll()
		return
	}
	r.SetCtxVar(consts.CtxAdminId, adminId)
	r.Middleware.Next()
}

// CORS allows the API to be consumed from other origins, which is handy when
// the navigation page is embedded into another site.
func CORS(r *ghttp.Request) {
	r.Response.CORSDefault()
	r.Middleware.Next()
}
