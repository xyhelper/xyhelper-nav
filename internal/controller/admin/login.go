package admin

import (
	"context"

	"github.com/gogf/gf/v2/frame/g"

	v1 "xyhelper-nav/api/nav/v1"
	"xyhelper-nav/internal/consts"
	"xyhelper-nav/internal/service"
)

// LoginController serves the public login endpoint. It lives in its own
// controller so that it can be bound outside of the authentication middleware.
type LoginController struct{}

// NewLogin creates and returns the login controller.
func NewLogin() *LoginController {
	return &LoginController{}
}

// Login verifies the credentials and opens an administration session.
func (c *LoginController) Login(ctx context.Context, req *v1.LoginReq) (res *v1.LoginRes, err error) {
	request := g.RequestFromCtx(ctx)
	account, err := service.Admin().Login(ctx, req.Username, req.Password, request.GetClientIp())
	if err != nil {
		return nil, err
	}
	if err = request.Session.Set(consts.SessionAdminId, account.Id); err != nil {
		return nil, err
	}
	if err = request.Session.Set(consts.SessionAdminName, account.Username); err != nil {
		return nil, err
	}
	return &v1.LoginRes{
		Id:       account.Id,
		Username: account.Username,
		Nickname: account.Nickname,
	}, nil
}
