// Package category implements the navigation category service.
package category

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"xyhelper-nav/internal/consts"
	"xyhelper-nav/internal/dao"
	"xyhelper-nav/internal/model"
	"xyhelper-nav/internal/model/do"
	"xyhelper-nav/internal/model/entity"
	"xyhelper-nav/internal/service"
)

type sCategory struct{}

func init() {
	service.RegisterCategory(New())
}

// New creates and returns the category service implementation.
func New() service.ICategory {
	return &sCategory{}
}

var (
	columns     = dao.NavCategory.Columns()
	siteColumns = dao.NavSite.Columns()
)

// ListEntities returns the raw category records the caller may see.
func (s *sCategory) ListEntities(ctx context.Context, query model.CategoryQuery) ([]*entity.NavCategory, error) {
	m := dao.NavCategory.Ctx(ctx)
	if query.OnlyEnabled {
		m = m.Where(columns.Status, consts.StatusOn)
	}
	if !query.IncludeRestricted {
		m = m.Where(columns.Visibility, consts.VisibilityPublic)
	}
	var categories []*entity.NavCategory
	if err := m.OrderAsc(columns.Sort).OrderAsc(columns.Id).Scan(&categories); err != nil {
		return nil, err
	}
	return categories, nil
}

// RestrictedIds returns the ids of the categories that require a login.
func (s *sCategory) RestrictedIds(ctx context.Context) ([]int, error) {
	values, err := dao.NavCategory.Ctx(ctx).
		Where(columns.Visibility, consts.VisibilityLogin).
		Array(columns.Id)
	if err != nil {
		return nil, err
	}
	return values.Ints(), nil
}

// List returns every category the caller may see, with its sites attached.
func (s *sCategory) List(ctx context.Context, query model.CategoryQuery) ([]*model.CategoryItem, error) {
	categories, err := s.ListEntities(ctx, query)
	if err != nil {
		return nil, err
	}
	sites, err := service.Site().List(ctx, model.SiteQuery{
		OnlyEnabled:       query.OnlyEnabled,
		IncludeRestricted: query.IncludeRestricted,
	})
	if err != nil {
		return nil, err
	}
	grouped := make(map[int][]*model.SiteItem, len(categories))
	for _, record := range sites {
		grouped[record.CategoryId] = append(grouped[record.CategoryId], model.NewSiteItem(record))
	}
	items := make([]*model.CategoryItem, 0, len(categories))
	for _, record := range categories {
		item := model.NewCategoryItem(record)
		if list, ok := grouped[record.Id]; ok {
			item.Sites = list
		}
		item.SiteCount = len(item.Sites)
		items = append(items, item)
	}
	return items, nil
}

// Get returns one category by id.
func (s *sCategory) Get(ctx context.Context, id int) (*entity.NavCategory, error) {
	var record *entity.NavCategory
	if err := dao.NavCategory.Ctx(ctx).Where(columns.Id, id).Scan(&record); err != nil {
		return nil, err
	}
	if record == nil {
		return nil, gerror.Newf("分类 %d 不存在", id)
	}
	return record, nil
}

// Create stores a new category and returns its id.
func (s *sCategory) Create(ctx context.Context, input *model.CategorySaveInput) (int, error) {
	now := gtime.Now()
	id, err := dao.NavCategory.Ctx(ctx).Data(do.NavCategory{
		Name:        strings.TrimSpace(input.Name),
		Icon:        strings.TrimSpace(input.Icon),
		Description: strings.TrimSpace(input.Description),
		Sort:        input.Sort,
		Status:      model.StatusValue(input.Status),
		Visibility:  model.VisibilityValue(input.Visibility),
		CreatedAt:   now,
		UpdatedAt:   now,
	}).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// Update modifies an existing category.
func (s *sCategory) Update(ctx context.Context, id int, input *model.CategorySaveInput) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	data := g.Map{
		columns.Name:        strings.TrimSpace(input.Name),
		columns.Icon:        strings.TrimSpace(input.Icon),
		columns.Description: strings.TrimSpace(input.Description),
		columns.Sort:        input.Sort,
		columns.UpdatedAt:   gtime.Now(),
	}
	if input.Status != nil {
		data[columns.Status] = model.StatusValue(input.Status)
	}
	if input.Visibility != nil {
		data[columns.Visibility] = model.VisibilityValue(input.Visibility)
	}
	_, err := dao.NavCategory.Ctx(ctx).Where(columns.Id, id).Data(data).Update()
	return err
}

// Delete removes a category and every site that belongs to it.
func (s *sCategory) Delete(ctx context.Context, id int) error {
	if _, err := s.Get(ctx, id); err != nil {
		return err
	}
	return g.DB().Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
		if _, err := dao.NavSite.Ctx(ctx).Where(siteColumns.CategoryId, id).Delete(); err != nil {
			return err
		}
		_, err := dao.NavCategory.Ctx(ctx).Where(columns.Id, id).Delete()
		return err
	})
}
