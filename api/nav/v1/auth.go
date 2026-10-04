package v1

import (
	"github.com/gogf/gf/v2/frame/g"
	"time"
)

// LoginReq authenticates an administrator and opens a session.
type LoginReq struct {
	g.Meta   `path:"/admin/login" method:"post" tags:"管理端-认证" summary:"管理员登录"`
	Username string `json:"username" v:"required#请输入用户名" description:"用户名"`
	Password string `json:"password" v:"required#请输入密码" description:"密码"`
}

// LoginRes is the response of a successful login.
type LoginRes struct {
	Id       int    `json:"id"       description:"管理员ID"`
	Username string `json:"username" description:"用户名"`
	Nickname string `json:"nickname" description:"昵称"`
}

// LogoutReq closes the current session.
type LogoutReq struct {
	g.Meta `path:"/admin/logout" method:"post" tags:"管理端-认证" summary:"退出登录"`
}

// LogoutRes is the response of a logout.
type LogoutRes struct{}

// ProfileReq returns the current administrator.
type ProfileReq struct {
	g.Meta `path:"/admin/profile" method:"get" tags:"管理端-认证" summary:"获取当前管理员"`
}

// ProfileRes is the current administrator profile.
type ProfileRes struct {
	Id          int       `json:"id"          description:"管理员ID"`
	Username    string    `json:"username"    description:"用户名"`
	Nickname    string    `json:"nickname"    description:"昵称"`
	LastLoginAt time.Time `json:"lastLoginAt" description:"最后登录时间"`
	LastLoginIp string    `json:"lastLoginIp" description:"最后登录IP"`
}

// ProfileUpdateReq updates the current administrator profile.
type ProfileUpdateReq struct {
	g.Meta   `path:"/admin/profile" method:"put" tags:"管理端-认证" summary:"更新管理员资料"`
	Nickname string `json:"nickname" v:"required#昵称不能为空" description:"昵称"`
}

// ProfileUpdateRes is the response of a profile update.
type ProfileUpdateRes struct{}

// PasswordReq changes the current administrator password.
type PasswordReq struct {
	g.Meta      `path:"/admin/password" method:"post" tags:"管理端-认证" summary:"修改密码"`
	OldPassword string `json:"oldPassword" v:"required#请输入原密码" description:"原密码"`
	NewPassword string `json:"newPassword" v:"required|length:6,64#请输入新密码|新密码长度需为6-64位" description:"新密码"`
}

// PasswordRes is the response of a password change.
type PasswordRes struct{}
