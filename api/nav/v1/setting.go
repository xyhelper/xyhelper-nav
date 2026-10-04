package v1

import (
	"github.com/gogf/gf/v2/frame/g"

	"xyhelper-nav/internal/model"
)

// SettingListReq returns the whole site configuration.
type SettingListReq struct {
	g.Meta `path:"/admin/settings" method:"get" tags:"管理端-配置" summary:"获取站点配置"`
}

// SettingListRes is the whole site configuration.
type SettingListRes struct {
	Settings map[string]string `json:"settings" description:"站点配置"`
}

// SettingUpdateReq writes one or more configuration entries.
type SettingUpdateReq struct {
	g.Meta   `path:"/admin/settings" method:"put" tags:"管理端-配置" summary:"保存站点配置"`
	Settings map[string]string `json:"settings" v:"required#配置内容不能为空" description:"站点配置"`
}

// SettingUpdateRes is the response of a configuration update.
type SettingUpdateRes struct{}

// BackupExportReq exports the whole navigation content.
type BackupExportReq struct {
	g.Meta `path:"/admin/backup/export" method:"get" tags:"管理端-数据" summary:"导出数据"`
}

// BackupExportRes is a portable snapshot of the navigation content.
type BackupExportRes struct {
	*model.BackupData
}

// BackupImportReq replaces the navigation content with a snapshot.
type BackupImportReq struct {
	g.Meta `path:"/admin/backup/import" method:"post" tags:"管理端-数据" summary:"导入数据"`
	Data   *model.BackupData `json:"data" v:"required#备份数据不能为空" description:"备份数据"`
}

// BackupImportRes is the response of an import.
type BackupImportRes struct{}
