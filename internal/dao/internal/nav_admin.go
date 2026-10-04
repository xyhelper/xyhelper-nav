// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// NavAdminDao is the data access object for the table nav_admin.
type NavAdminDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  NavAdminColumns    // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// NavAdminColumns defines and stores column names for the table nav_admin.
type NavAdminColumns struct {
	Id          string //
	Username    string //
	Password    string //
	Nickname    string //
	LastLoginAt string //
	LastLoginIp string //
	CreatedAt   string //
	UpdatedAt   string //
}

// navAdminColumns holds the columns for the table nav_admin.
var navAdminColumns = NavAdminColumns{
	Id:          "id",
	Username:    "username",
	Password:    "password",
	Nickname:    "nickname",
	LastLoginAt: "last_login_at",
	LastLoginIp: "last_login_ip",
	CreatedAt:   "created_at",
	UpdatedAt:   "updated_at",
}

// NewNavAdminDao creates and returns a new DAO object for table data access.
func NewNavAdminDao(handlers ...gdb.ModelHandler) *NavAdminDao {
	return &NavAdminDao{
		group:    "default",
		table:    "nav_admin",
		columns:  navAdminColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *NavAdminDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *NavAdminDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *NavAdminDao) Columns() NavAdminColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *NavAdminDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *NavAdminDao) Ctx(ctx context.Context) *gdb.Model {
	model := dao.DB().Model(dao.table)
	for _, handler := range dao.handlers {
		model = handler(model)
	}
	return model.Safe().Ctx(ctx)
}

// Transaction wraps the transaction logic using function f.
// It rolls back the transaction and returns the error if function f returns a non-nil error.
// It commits the transaction and returns nil if function f returns nil.
//
// Note: Do not commit or roll back the transaction in function f,
// as it is automatically handled by this function.
func (dao *NavAdminDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
