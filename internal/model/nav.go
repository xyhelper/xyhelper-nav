// Package model contains the domain models shared by the service layer,
// the controllers and the API definitions.
package model

import (
	"fmt"
	"time"

	"xyhelper-nav/internal/consts"
	"xyhelper-nav/internal/model/entity"
)

// SiteItem is a single navigation entry. It is used both by the public page and
// by the administration panel.
type SiteItem struct {
	Id          int       `json:"id"          description:"站点ID"`
	CategoryId  int       `json:"categoryId"  description:"所属分类ID"`
	Name        string    `json:"name"        description:"站点名称"`
	Url         string    `json:"url"         description:"站点地址"`
	Icon        string    `json:"icon"        description:"图标来源地址"`
	IconUrl     string    `json:"iconUrl"     description:"本地缓存图标地址"`
	Description string    `json:"description" description:"站点描述"`
	Sort        int       `json:"sort"        description:"排序值"`
	Status      int       `json:"status"      description:"状态 1显示 0隐藏"`
	Hits        int       `json:"hits"        description:"点击次数"`
	CreatedAt   time.Time `json:"createdAt"   description:"创建时间"`
	UpdatedAt   time.Time `json:"updatedAt"   description:"更新时间"`
}

// IconRefreshResult reports the outcome of an icon cache refresh.
type IconRefreshResult struct {
	Total     int `json:"total"     description:"处理的站点数量"`
	Refreshed int `json:"refreshed" description:"刷新成功的数量"`
	Failed    int `json:"failed"    description:"刷新失败的数量"`
}

// CategoryItem is a category together with its navigation entries.
type CategoryItem struct {
	Id          int         `json:"id"          description:"分类ID"`
	Name        string      `json:"name"        description:"分类名称"`
	Icon        string      `json:"icon"        description:"分类图标"`
	Description string      `json:"description" description:"分类描述"`
	Sort        int         `json:"sort"        description:"排序值"`
	Status      int         `json:"status"      description:"状态 1显示 0隐藏"`
	Visibility  int         `json:"visibility"  description:"可见性 0公开 1登录后可见"`
	SiteCount   int         `json:"siteCount"   description:"站点数量"`
	Sites       []*SiteItem `json:"sites"       description:"站点列表"`
	CreatedAt   time.Time   `json:"createdAt"   description:"创建时间"`
	UpdatedAt   time.Time   `json:"updatedAt"   description:"更新时间"`
}

// CategoryQuery filters a category listing.
type CategoryQuery struct {
	// OnlyEnabled drops categories whose status is "hidden".
	OnlyEnabled bool
	// IncludeRestricted keeps categories that require a login. Administrators
	// pass true, anonymous visitors pass false.
	IncludeRestricted bool
}

// SiteQuery filters a site listing or a site search.
type SiteQuery struct {
	// CategoryId limits the result to a single category, 0 means every category.
	CategoryId int
	// Keyword is the search term, it is only used by search.
	Keyword string
	// Limit caps the number of returned records, it is only used by search.
	Limit int
	// OnlyEnabled drops sites whose status is "hidden".
	OnlyEnabled bool
	// IncludeRestricted keeps the sites of categories that require a login.
	IncludeRestricted bool
}

// NavIndex is the whole payload required to render the public navigation page.
type NavIndex struct {
	Settings   map[string]string `json:"settings"   description:"站点配置"`
	Categories []*CategoryItem   `json:"categories" description:"分类与站点"`
	Total      int               `json:"total"      description:"站点总数"`
}

// SearchResult is the payload of the site search API.
type SearchResult struct {
	Keyword string      `json:"keyword" description:"搜索关键字"`
	Total   int         `json:"total"   description:"命中数量"`
	Sites   []*SiteItem `json:"sites"   description:"命中的站点"`
}

// CategorySaveInput is the writable payload of a category.
type CategorySaveInput struct {
	Name        string `json:"name"        v:"required#分类名称不能为空" description:"分类名称"`
	Icon        string `json:"icon"        description:"分类图标"`
	Description string `json:"description" description:"分类描述"`
	Sort        int    `json:"sort"        description:"排序值"`
	Status      *int   `json:"status"      description:"状态 1显示 0隐藏"`
	Visibility  *int   `json:"visibility"  description:"可见性 0公开 1登录后可见"`
}

// SiteSaveInput is the writable payload of a navigation entry.
type SiteSaveInput struct {
	CategoryId  int    `json:"categoryId"  v:"required#请选择所属分类" description:"所属分类ID"`
	Name        string `json:"name"        v:"required#站点名称不能为空" description:"站点名称"`
	Url         string `json:"url"         v:"required#站点地址不能为空" description:"站点地址"`
	Icon        string `json:"icon"        description:"站点图标"`
	Description string `json:"description" description:"站点描述"`
	Sort        int    `json:"sort"        description:"排序值"`
	Status      *int   `json:"status"      description:"状态 1显示 0隐藏"`
}

// NewSiteItem converts a database entity into the shared site model.
func NewSiteItem(site *entity.NavSite) *SiteItem {
	if site == nil {
		return nil
	}
	return &SiteItem{
		Id:          site.Id,
		CategoryId:  site.CategoryId,
		Name:        site.Name,
		Url:         site.Url,
		Icon:        site.Icon,
		IconUrl:     IconUrl(site),
		Description: site.Description,
		Sort:        site.Sort,
		Status:      site.Status,
		Hits:        site.Hits,
		CreatedAt:   site.CreatedAt,
		UpdatedAt:   site.UpdatedAt,
	}
}

// IconUrl returns the local address that serves the cached icon of a site.
// The browser therefore never contacts the remote site directly. The cache key
// is appended as a version so that the response can be cached aggressively and
// still changes as soon as the icon is refreshed.
func IconUrl(site *entity.NavSite) string {
	if site == nil || (site.Icon == "" && site.Url == "") {
		return ""
	}
	address := fmt.Sprintf("%s/%d", consts.IconEndpointPrefix, site.Id)
	if site.IconKey != "" {
		address += "?v=" + site.IconKey
	}
	return address
}

// NewSiteItems converts a slice of database entities into site models.
func NewSiteItems(sites []*entity.NavSite) []*SiteItem {
	items := make([]*SiteItem, 0, len(sites))
	for _, site := range sites {
		items = append(items, NewSiteItem(site))
	}
	return items
}

// NewCategoryItem converts a database entity into the shared category model.
func NewCategoryItem(category *entity.NavCategory) *CategoryItem {
	if category == nil {
		return nil
	}
	return &CategoryItem{
		Id:          category.Id,
		Name:        category.Name,
		Icon:        category.Icon,
		Description: category.Description,
		Sort:        category.Sort,
		Status:      category.Status,
		Visibility:  category.Visibility,
		Sites:       make([]*SiteItem, 0),
		CreatedAt:   category.CreatedAt,
		UpdatedAt:   category.UpdatedAt,
	}
}

// StatusValue resolves an optional status into the value stored in the
// database. A missing status means "visible".
func StatusValue(status *int) int {
	if status != nil && *status == consts.StatusOff {
		return consts.StatusOff
	}
	return consts.StatusOn
}

// VisibilityValue resolves an optional visibility into the value stored in the
// database. A missing visibility means "public".
func VisibilityValue(visibility *int) int {
	if visibility != nil && *visibility == consts.VisibilityLogin {
		return consts.VisibilityLogin
	}
	return consts.VisibilityPublic
}
