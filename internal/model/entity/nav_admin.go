// =================================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// =================================================================================

package entity

import (
	"time"
)

// NavAdmin is the golang structure for table nav_admin.
type NavAdmin struct {
	Id          int       `json:"id"          orm:"id"            description:""` //
	Username    string    `json:"username"    orm:"username"      description:""` //
	Password    string    `json:"password"    orm:"password"      description:""` //
	Nickname    string    `json:"nickname"    orm:"nickname"      description:""` //
	LastLoginAt time.Time `json:"lastLoginAt" orm:"last_login_at" description:""` //
	LastLoginIp string    `json:"lastLoginIp" orm:"last_login_ip" description:""` //
	CreatedAt   time.Time `json:"createdAt"   orm:"created_at"    description:""` //
	UpdatedAt   time.Time `json:"updatedAt"   orm:"updated_at"    description:""` //
}
