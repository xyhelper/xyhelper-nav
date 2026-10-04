// Package admin implements administrator accounts and authentication.
package admin

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"golang.org/x/crypto/bcrypt"

	"xyhelper-nav/internal/dao"
	"xyhelper-nav/internal/model/entity"
	"xyhelper-nav/internal/service"
)

type sAdmin struct{}

func init() {
	service.RegisterAdmin(New())
}

// New creates and returns the administrator service implementation.
func New() service.IAdmin {
	return &sAdmin{}
}

// Login verifies the credentials, records the login and returns the account.
func (s *sAdmin) Login(ctx context.Context, username, password, ip string) (*entity.NavAdmin, error) {
	username = strings.TrimSpace(username)
	if username == "" || password == "" {
		return nil, gerror.New("用户名和密码不能为空")
	}
	var account *entity.NavAdmin
	if err := dao.NavAdmin.Ctx(ctx).
		Where(dao.NavAdmin.Columns().Username, username).
		Scan(&account); err != nil {
		return nil, err
	}
	if account == nil {
		return nil, gerror.New("用户名或密码错误")
	}
	if err := bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(password)); err != nil {
		return nil, gerror.New("用户名或密码错误")
	}
	if _, err := dao.NavAdmin.Ctx(ctx).
		Where(dao.NavAdmin.Columns().Id, account.Id).
		Data(g.Map{
			dao.NavAdmin.Columns().LastLoginAt: gtime.Now(),
			dao.NavAdmin.Columns().LastLoginIp: ip,
			dao.NavAdmin.Columns().UpdatedAt:   gtime.Now(),
		}).Update(); err != nil {
		return nil, err
	}
	return account, nil
}

// Get returns one administrator by id.
func (s *sAdmin) Get(ctx context.Context, id int) (*entity.NavAdmin, error) {
	var account *entity.NavAdmin
	if err := dao.NavAdmin.Ctx(ctx).Where(dao.NavAdmin.Columns().Id, id).Scan(&account); err != nil {
		return nil, err
	}
	if account == nil {
		return nil, gerror.New("管理员不存在")
	}
	return account, nil
}

// ChangePassword replaces the password after verifying the current one.
func (s *sAdmin) ChangePassword(ctx context.Context, id int, oldPassword, newPassword string) error {
	if len(strings.TrimSpace(newPassword)) < 6 {
		return gerror.New("新密码长度不能少于 6 位")
	}
	account, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if err = bcrypt.CompareHashAndPassword([]byte(account.Password), []byte(oldPassword)); err != nil {
		return gerror.New("原密码不正确")
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	_, err = dao.NavAdmin.Ctx(ctx).
		Where(dao.NavAdmin.Columns().Id, id).
		Data(g.Map{
			dao.NavAdmin.Columns().Password:  string(hashed),
			dao.NavAdmin.Columns().UpdatedAt: gtime.Now(),
		}).Update()
	return err
}

// UpdateProfile updates the mutable profile fields.
func (s *sAdmin) UpdateProfile(ctx context.Context, id int, nickname string) error {
	_, err := dao.NavAdmin.Ctx(ctx).
		Where(dao.NavAdmin.Columns().Id, id).
		Data(g.Map{
			dao.NavAdmin.Columns().Nickname:  strings.TrimSpace(nickname),
			dao.NavAdmin.Columns().UpdatedAt: gtime.Now(),
		}).Update()
	return err
}
