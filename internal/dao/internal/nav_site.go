// ==========================================================================
// Code generated and maintained by GoFrame CLI tool. DO NOT EDIT.
// ==========================================================================

package internal

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
)

// NavSiteDao is the data access object for the table nav_site.
type NavSiteDao struct {
	table    string             // table is the underlying table name of the DAO.
	group    string             // group is the database configuration group name of the current DAO.
	columns  NavSiteColumns     // columns contains all the column names of Table for convenient usage.
	handlers []gdb.ModelHandler // handlers for customized model modification.
}

// NavSiteColumns defines and stores column names for the table nav_site.
type NavSiteColumns struct {
	Id            string //
	CategoryId    string //
	Name          string //
	Url           string //
	Icon          string //
	Description   string //
	Sort          string //
	Status        string //
	Hits          string //
	IconKey       string //
	IconFetchedAt string //
	CreatedAt     string //
	UpdatedAt     string //
}

// navSiteColumns holds the columns for the table nav_site.
var navSiteColumns = NavSiteColumns{
	Id:            "id",
	CategoryId:    "category_id",
	Name:          "name",
	Url:           "url",
	Icon:          "icon",
	Description:   "description",
	Sort:          "sort",
	Status:        "status",
	Hits:          "hits",
	IconKey:       "icon_key",
	IconFetchedAt: "icon_fetched_at",
	CreatedAt:     "created_at",
	UpdatedAt:     "updated_at",
}

// NewNavSiteDao creates and returns a new DAO object for table data access.
func NewNavSiteDao(handlers ...gdb.ModelHandler) *NavSiteDao {
	return &NavSiteDao{
		group:    "default",
		table:    "nav_site",
		columns:  navSiteColumns,
		handlers: handlers,
	}
}

// DB retrieves and returns the underlying raw database management object of the current DAO.
func (dao *NavSiteDao) DB() gdb.DB {
	return g.DB(dao.group)
}

// Table returns the table name of the current DAO.
func (dao *NavSiteDao) Table() string {
	return dao.table
}

// Columns returns all column names of the current DAO.
func (dao *NavSiteDao) Columns() NavSiteColumns {
	return dao.columns
}

// Group returns the database configuration group name of the current DAO.
func (dao *NavSiteDao) Group() string {
	return dao.group
}

// Ctx creates and returns a Model for the current DAO. It automatically sets the context for the current operation.
func (dao *NavSiteDao) Ctx(ctx context.Context) *gdb.Model {
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
func (dao *NavSiteDao) Transaction(ctx context.Context, f func(ctx context.Context, tx gdb.TX) error) (err error) {
	return dao.Ctx(ctx).Transaction(ctx, f)
}
