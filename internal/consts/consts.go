// Package consts holds project-wide constants.
package consts

const (
	// SessionAdminId is the session key that stores the logged-in administrator id.
	// SessionAdminName is the session key that stores the logged-in administrator name.
	SessionAdminId   = "xyhelperNavAdminId"
	SessionAdminName = "xyhelperNavAdminName"

	// CtxAdminId is the context key for the current administrator id.
	CtxAdminId = "xyhelperNavCtxAdminId"

	// IconEndpointPrefix is the public route that serves cached site icons.
	IconEndpointPrefix = "/api/v1/nav/icon"

	// DefaultAdminUser is the username created on first boot when no administrator exists.
	DefaultAdminUser = "admin"
	// DefaultAdminPass is the password created on first boot when no administrator exists.
	DefaultAdminPass = "admin123"

	// StatusOn marks a record as visible on the public navigation page.
	StatusOn = 1
	// StatusOff marks a record as hidden from the public navigation page.
	StatusOff = 0
)

// Visibility levels of a category.
const (
	// VisibilityPublic means everybody can see the category, including guests.
	VisibilityPublic = 0
	// VisibilityLogin means the category is only listed for a signed-in
	// administrator.
	VisibilityLogin = 1
)

// Setting keys of the nav_setting table.
const (
	SettingSiteName     = "site_name"
	SettingSiteSubtitle = "site_subtitle"
	SettingDescription  = "site_description"
	SettingKeywords     = "site_keywords"
	SettingFooter       = "site_footer"
	SettingThemeColor   = "theme_color"
	SettingDefaultTheme = "default_theme"
	SettingOpenNewTab   = "open_new_tab"
	SettingShowHits     = "show_hits"
	SettingBackground   = "background_image"
	SettingAnnouncement = "announcement"
	SettingICP          = "icp"
)

// DefaultSettings are written to the database on first boot.
var DefaultSettings = map[string]string{
	SettingSiteName:     "XYHelper 导航",
	SettingSiteSubtitle: "轻量 · 美观 · 自托管导航页",
	SettingDescription:  "XYHelper Nav 是一个基于 GoFrame 构建的轻量美观导航页，支持分类管理、站点收藏与全文搜索。",
	SettingKeywords:     "导航,网址导航,XYHelper,GoFrame",
	SettingFooter:       "Powered by xyhelper-nav · GoFrame",
	SettingThemeColor:   "#6366f1",
	SettingDefaultTheme: "dark",
	SettingOpenNewTab:   "1",
	SettingShowHits:     "1",
	SettingBackground:   "",
	SettingAnnouncement: "",
	SettingICP:          "",
}
