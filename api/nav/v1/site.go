package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"xyhelper-nav/internal/model"
)

// SiteListReq returns the navigation entries of the administration panel.
type SiteListReq struct {
	g.Meta     `path:"/admin/sites" method:"get" tags:"管理端-站点" summary:"站点列表"`
	CategoryId int    `json:"categoryId" in:"query" description:"分类ID, 0表示全部"`
	Keyword    string `json:"keyword"    in:"query" description:"名称关键字"`
}

// SiteListRes is the site list.
type SiteListRes struct {
	List  []*model.SiteItem `json:"list"  description:"站点列表"`
	Total int               `json:"total" description:"站点数量"`
}

// SiteCreateReq creates a navigation entry.
type SiteCreateReq struct {
	g.Meta      `path:"/admin/site" method:"post" tags:"管理端-站点" summary:"新增站点"`
	CategoryId  int    `json:"categoryId"  v:"required#请选择所属分类" description:"所属分类ID"`
	Name        string `json:"name"        v:"required#站点名称不能为空" description:"站点名称"`
	Url         string `json:"url"         v:"required#站点地址不能为空" description:"站点地址"`
	Icon        string `json:"icon"        description:"站点图标"`
	Description string `json:"description" description:"站点描述"`
	Sort        int    `json:"sort"        description:"排序值"`
	Status      *int   `json:"status"      description:"状态 1显示 0隐藏"`
}

// SiteCreateRes is the id of the created site.
type SiteCreateRes struct {
	Id int `json:"id" description:"站点ID"`
}

// SiteUpdateReq updates a navigation entry.
type SiteUpdateReq struct {
	g.Meta      `path:"/admin/site" method:"put" tags:"管理端-站点" summary:"更新站点"`
	Id          int    `json:"id"          v:"required#站点ID不能为空" description:"站点ID"`
	CategoryId  int    `json:"categoryId"  v:"required#请选择所属分类" description:"所属分类ID"`
	Name        string `json:"name"        v:"required#站点名称不能为空" description:"站点名称"`
	Url         string `json:"url"         v:"required#站点地址不能为空" description:"站点地址"`
	Icon        string `json:"icon"        description:"站点图标"`
	Description string `json:"description" description:"站点描述"`
	Sort        int    `json:"sort"        description:"排序值"`
	Status      *int   `json:"status"      description:"状态 1显示 0隐藏"`
}

// SiteUpdateRes is the response of a site update.
type SiteUpdateRes struct{}

// SiteDeleteReq deletes a navigation entry.
type SiteDeleteReq struct {
	g.Meta `path:"/admin/site" method:"delete" tags:"管理端-站点" summary:"删除站点"`
	Id     int `json:"id" in:"query" v:"required#站点ID不能为空" description:"站点ID"`
}

// SiteDeleteRes is the response of a site deletion.
type SiteDeleteRes struct{}

// IconRefreshReq re-downloads cached site icons.
type IconRefreshReq struct {
	g.Meta `path:"/admin/icons/refresh" method:"post" tags:"管理端-站点" summary:"刷新站点图标缓存"`
	Id     int `json:"id" description:"站点ID, 0 表示全部站点"`
}

// IconRefreshRes reports how many icons were refreshed.
type IconRefreshRes struct {
	Total     int `json:"total"     description:"处理的站点数量"`
	Refreshed int `json:"refreshed" description:"刷新成功的数量"`
	Failed    int `json:"failed"    description:"刷新失败的数量"`
}
