package admin

import (
	"context"

	v1 "xyhelper-nav/api/nav/v1"
	"xyhelper-nav/internal/service"
)

// SettingList returns the whole site configuration.
func (c *Controller) SettingList(ctx context.Context, req *v1.SettingListReq) (res *v1.SettingListRes, err error) {
	settings, err := service.Setting().All(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.SettingListRes{Settings: settings}, nil
}

// SettingUpdate writes one or more configuration entries.
func (c *Controller) SettingUpdate(ctx context.Context, req *v1.SettingUpdateReq) (res *v1.SettingUpdateRes, err error) {
	if err = service.Setting().SetMap(ctx, req.Settings); err != nil {
		return nil, err
	}
	return &v1.SettingUpdateRes{}, nil
}

// BackupExport returns the whole navigation content as a portable document.
func (c *Controller) BackupExport(ctx context.Context, req *v1.BackupExportReq) (res *v1.BackupExportRes, err error) {
	data, err := service.Backup().Export(ctx)
	if err != nil {
		return nil, err
	}
	return &v1.BackupExportRes{BackupData: data}, nil
}

// BackupImport replaces the navigation content with the submitted document.
func (c *Controller) BackupImport(ctx context.Context, req *v1.BackupImportReq) (res *v1.BackupImportRes, err error) {
	if err = service.Backup().Import(ctx, req.Data); err != nil {
		return nil, err
	}
	return &v1.BackupImportRes{}, nil
}
