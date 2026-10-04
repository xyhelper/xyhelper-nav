package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"xyhelper-nav/internal/model"
)

// CategoryListReq returns every category together with its sites.
type CategoryListReq struct {
	g.Meta `path:"/admin/categories" method:"get" tags:"管理端-分类" summary:"分类列表"`
}

// CategoryListRes is the category list.
type CategoryListRes struct {
	List  []*model.CategoryItem `json:"list"  description:"分类列表"`
	Total int                   `json:"total" description:"分类数量"`
}

// CategoryCreateReq creates a category.
type CategoryCreateReq struct {
	g.Meta      `path:"/admin/category" method:"post" tags:"管理端-分类" summary:"新增分类"`
	Name        string `json:"name"        v:"required#分类名称不能为空" description:"分类名称"`
	Icon        string `json:"icon"        description:"分类图标"`
	Description string `json:"description" description:"分类描述"`
	Sort        int    `json:"sort"        description:"排序值"`
	Status      *int   `json:"status"      description:"状态 1显示 0隐藏"`
	Visibility  *int   `json:"visibility"  description:"可见性 0公开 1登录后可见"`
}

// CategoryCreateRes is the id of the created category.
type CategoryCreateRes struct {
	Id int `json:"id" description:"分类ID"`
}

// CategoryUpdateReq updates a category.
type CategoryUpdateReq struct {
	g.Meta      `path:"/admin/category" method:"put" tags:"管理端-分类" summary:"更新分类"`
	Id          int    `json:"id"          v:"required#分类ID不能为空" description:"分类ID"`
	Name        string `json:"name"        v:"required#分类名称不能为空" description:"分类名称"`
	Icon        string `json:"icon"        description:"分类图标"`
	Description string `json:"description" description:"分类描述"`
	Sort        int    `json:"sort"        description:"排序值"`
	Status      *int   `json:"status"      description:"状态 1显示 0隐藏"`
	Visibility  *int   `json:"visibility"  description:"可见性 0公开 1登录后可见"`
}

// CategoryUpdateRes is the response of a category update.
type CategoryUpdateRes struct{}

// CategoryDeleteReq deletes a category and all of its sites.
type CategoryDeleteReq struct {
	g.Meta `path:"/admin/category" method:"delete" tags:"管理端-分类" summary:"删除分类"`
	Id     int `json:"id" in:"query" v:"required#分类ID不能为空" description:"分类ID"`
}

// CategoryDeleteRes is the response of a category deletion.
type CategoryDeleteRes struct{}
