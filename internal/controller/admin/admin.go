// Package admin implements the administration API controllers.
package admin

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	v1 "xyhelper-nav/api/nav/v1"
	"xyhelper-nav/internal/middleware"
	"xyhelper-nav/internal/service"
)

// Controller serves the protected administration API.
type Controller struct{}

// New creates and returns the administration controller.
func New() *Controller {
	return &Controller{}
}

// Logout closes the current administration session.
func (c *Controller) Logout(ctx context.Context, req *v1.LogoutReq) (res *v1.LogoutRes, err error) {
	if err = g.RequestFromCtx(ctx).Session.RemoveAll(); err != nil {
		return nil, err
	}
	return &v1.LogoutRes{}, nil
}

// Profile returns the current administrator.
func (c *Controller) Profile(ctx context.Context, req *v1.ProfileReq) (res *v1.ProfileRes, err error) {
	account, err := service.Admin().Get(ctx, middleware.CurrentAdminId(ctx))
	if err != nil {
		return nil, err
	}
	return &v1.ProfileRes{
		Id:          account.Id,
		Username:    account.Username,
		Nickname:    account.Nickname,
		LastLoginAt: account.LastLoginAt,
		LastLoginIp: account.LastLoginIp,
	}, nil
}

// ProfileUpdate updates the current administrator profile.
func (c *Controller) ProfileUpdate(ctx context.Context, req *v1.ProfileUpdateReq) (res *v1.ProfileUpdateRes, err error) {
	if err = service.Admin().UpdateProfile(ctx, middleware.CurrentAdminId(ctx), req.Nickname); err != nil {
		return nil, err
	}
	return &v1.ProfileUpdateRes{}, nil
}

// Password changes the current administrator password.
func (c *Controller) Password(ctx context.Context, req *v1.PasswordReq) (res *v1.PasswordRes, err error) {
	if err = service.Admin().ChangePassword(ctx, middleware.CurrentAdminId(ctx), req.OldPassword, req.NewPassword); err != nil {
		return nil, err
	}
	return &v1.PasswordRes{}, nil
}
