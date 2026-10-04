// Package backup implements import and export of the navigation content.
package backup

import (
	"context"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"xyhelper-nav/internal/dao"
	"xyhelper-nav/internal/model"
	"xyhelper-nav/internal/model/do"
	"xyhelper-nav/internal/model/entity"
	"xyhelper-nav/internal/service"
)

type sBackup struct{}

func init() {
	service.RegisterBackup(New())
}

// New creates and returns the backup service implementation.
func New() service.IBackup {
	return &sBackup{}
}

// Export returns the current content as a backup document.
func (s *sBackup) Export(ctx context.Context) (*model.BackupData, error) {
	settings, err := service.Setting().All(ctx)
	if err != nil {
		return nil, err
	}
	var categories []*entity.NavCategory
	if err = dao.NavCategory.Ctx(ctx).
		OrderAsc(dao.NavCategory.Columns().Sort).
		OrderAsc(dao.NavCategory.Columns().Id).
		Scan(&categories); err != nil {
		return nil, err
	}
	var sites []*entity.NavSite
	if err = dao.NavSite.Ctx(ctx).
		OrderAsc(dao.NavSite.Columns().Sort).
		OrderAsc(dao.NavSite.Columns().Id).
		Scan(&sites); err != nil {
		return nil, err
	}
	grouped := make(map[int][]*model.BackupSite, len(categories))
	for _, record := range sites {
		grouped[record.CategoryId] = append(grouped[record.CategoryId], &model.BackupSite{
			Name:        record.Name,
			Url:         record.Url,
			Icon:        record.Icon,
			Description: record.Description,
			Sort:        record.Sort,
			Status:      record.Status,
			Hits:        record.Hits,
		})
	}
	data := &model.BackupData{
		Version:    model.BackupFormatVersion,
		ExportedAt: gtime.Now(),
		Settings:   settings,
		Categories: make([]*model.BackupCategory, 0, len(categories)),
	}
	for _, record := range categories {
		item := &model.BackupCategory{
			Name:        record.Name,
			Icon:        record.Icon,
			Description: record.Description,
			Sort:        record.Sort,
			Status:      record.Status,
			Visibility:  record.Visibility,
			Sites:       grouped[record.Id],
		}
		if item.Sites == nil {
			item.Sites = make([]*model.BackupSite, 0)
		}
		data.Categories = append(data.Categories, item)
	}
	return data, nil
}

// Import replaces the current content with the given document.
func (s *sBackup) Import(ctx context.Context, data *model.BackupData) error {
	if data == nil || len(data.Categories) == 0 {
		return gerror.New("备份数据为空或格式不正确")
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := dao.NavSite.Ctx(ctx).Where("1=1").Delete(); err != nil {
			return err
		}
		if _, err := dao.NavCategory.Ctx(ctx).Where("1=1").Delete(); err != nil {
			return err
		}
		now := gtime.Now()
		for _, category := range data.Categories {
			categoryId, err := dao.NavCategory.Ctx(ctx).Data(do.NavCategory{
				Name:        category.Name,
				Icon:        category.Icon,
				Description: category.Description,
				Sort:        category.Sort,
				Status:      category.Status,
				Visibility:  category.Visibility,
				CreatedAt:   now,
				UpdatedAt:   now,
			}).InsertAndGetId()
			if err != nil {
				return err
			}
			for _, site := range category.Sites {
				if _, err = dao.NavSite.Ctx(ctx).Data(do.NavSite{
					CategoryId:  int(categoryId),
					Name:        site.Name,
					Url:         site.Url,
					Icon:        site.Icon,
					Description: site.Description,
					Sort:        site.Sort,
					Status:      site.Status,
					Hits:        site.Hits,
					CreatedAt:   now,
					UpdatedAt:   now,
				}).Insert(); err != nil {
					return err
				}
			}
		}
		if len(data.Settings) > 0 {
			for key, value := range data.Settings {
				if err := service.Setting().Set(ctx, key, value); err != nil {
					return err
				}
			}
		}
		return nil
	})
}
