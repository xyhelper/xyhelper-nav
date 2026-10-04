package model

import "github.com/gogf/gf/v2/os/gtime"

// BackupFormatVersion is the version of the exported JSON document.
const BackupFormatVersion = "1"

// BackupData is the full, portable snapshot of the navigation content.
type BackupData struct {
	Version    string            `json:"version"    description:"备份格式版本"`
	ExportedAt *gtime.Time       `json:"exportedAt" description:"导出时间"`
	Settings   map[string]string `json:"settings"   description:"站点配置"`
	Categories []*BackupCategory `json:"categories" description:"分类与站点"`
}

// BackupCategory is a category inside a backup document.
type BackupCategory struct {
	Name        string        `json:"name"        description:"分类名称"`
	Icon        string        `json:"icon"        description:"分类图标"`
	Description string        `json:"description" description:"分类描述"`
	Sort        int           `json:"sort"        description:"排序值"`
	Status      int           `json:"status"      description:"状态 1显示 0隐藏"`
	Visibility  int           `json:"visibility"  description:"可见性 0公开 1登录后可见"`
	Sites       []*BackupSite `json:"sites"       description:"站点列表"`
}

// BackupSite is a navigation entry inside a backup document.
type BackupSite struct {
	Name        string `json:"name"        description:"站点名称"`
	Url         string `json:"url"         description:"站点地址"`
	Icon        string `json:"icon"        description:"站点图标"`
	Description string `json:"description" description:"站点描述"`
	Sort        int    `json:"sort"        description:"排序值"`
	Status      int    `json:"status"      description:"状态 1显示 0隐藏"`
	Hits        int    `json:"hits"        description:"点击次数"`
}
