package service

import (
	"context"

	"xyhelper-nav/internal/model/entity"
)

// IAdmin manages administrator accounts and authentication.
type IAdmin interface {
	// Login verifies the credentials, records the login and returns the account.
	Login(ctx context.Context, username, password, ip string) (*entity.NavAdmin, error)
	// Get returns one administrator by id.
	Get(ctx context.Context, id int) (*entity.NavAdmin, error)
	// ChangePassword replaces the password after verifying the current one.
	ChangePassword(ctx context.Context, id int, oldPassword, newPassword string) error
	// UpdateProfile updates the mutable profile fields.
	UpdateProfile(ctx context.Context, id int, nickname string) error
}

var localAdmin IAdmin

// Admin returns the administrator service implementation.
func Admin() IAdmin {
	if localAdmin == nil {
		panic("implementation of IAdmin is not registered, did you forget to import internal/logic?")
	}
	return localAdmin
}

// RegisterAdmin registers the administrator service implementation.
func RegisterAdmin(impl IAdmin) {
	localAdmin = impl
}
