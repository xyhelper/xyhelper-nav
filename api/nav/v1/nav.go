// Package v1 defines the API contract of the navigation service.
// Every request struct carries a g.Meta tag so GoFrame can bind it to a route.
package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"xyhelper-nav/internal/model"
)

// IndexReq is the request of the public navigation page payload.
type IndexReq struct {
	g.Meta `path:"/nav/index" method:"get" tags:"前台导航" summary:"获取导航页全部数据"`
}

// IndexRes is the response of the public navigation page payload.
type IndexRes struct {
	Settings   map[string]string     `json:"settings"   description:"站点配置"`
	Categories []*model.CategoryItem `json:"categories" description:"分类与站点"`
	Total      int                   `json:"total"      description:"站点总数"`
}

// HitReq records one visit of a navigation entry.
type HitReq struct {
	g.Meta `path:"/nav/hit" method:"post" tags:"前台导航" summary:"记录站点点击"`
	Id     int `json:"id" v:"required#站点ID不能为空" description:"站点ID"`
}

// HitRes is the response of the hit API.
type HitRes struct {
	Id   int    `json:"id"   description:"站点ID"`
	Url  string `json:"url"  description:"站点地址"`
	Hits int    `json:"hits" description:"最新点击次数"`
}

// SearchReq searches navigation entries.
type SearchReq struct {
	g.Meta  `path:"/nav/search" method:"get" tags:"前台导航" summary:"搜索站点"`
	Keyword string `json:"keyword" in:"query" description:"搜索关键字"`
	Limit   int    `json:"limit"   in:"query" description:"返回数量上限"`
}

// SearchRes is the response of the search API.
type SearchRes struct {
	Keyword string            `json:"keyword" description:"搜索关键字"`
	Total   int               `json:"total"   description:"命中数量"`
	Sites   []*model.SiteItem `json:"sites"   description:"命中的站点"`
}
