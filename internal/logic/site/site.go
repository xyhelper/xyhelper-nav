// Package site implements the navigation entry service.
package site

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

type sSite struct{}

func init() {
	service.RegisterSite(New())
}

// New creates and returns the site service implementation.
func New() service.ISite {
	return &sSite{}
}

var columns = dao.NavSite.Columns()

// excludeRestricted removes the sites of every login-protected category from
// the query. It is a no-op when the caller may see restricted content.
func (s *sSite) excludeRestricted(ctx context.Context, m *gdb.Model, includeRestricted bool) (*gdb.Model, error) {
	if includeRestricted {
		return m, nil
	}
	restricted, err := service.Category().RestrictedIds(ctx)
	if err != nil {
		return nil, err
	}
	if len(restricted) == 0 {
		return m, nil
	}
	return m.WhereNotIn(columns.CategoryId, restricted), nil
}

// List returns the sites matching the query.
func (s *sSite) List(ctx context.Context, query model.SiteQuery) ([]*entity.NavSite, error) {
	m := dao.NavSite.Ctx(ctx)
	if query.CategoryId > 0 {
		m = m.Where(columns.CategoryId, query.CategoryId)
	}
	if query.OnlyEnabled {
		m = m.Where(columns.Status, consts.StatusOn)
	}
	m, err := s.excludeRestricted(ctx, m, query.IncludeRestricted)
	if err != nil {
		return nil, err
	}
	var sites []*entity.NavSite
	if err = m.OrderAsc(columns.Sort).OrderAsc(columns.Id).Scan(&sites); err != nil {
		return nil, err
	}
	return sites, nil
}

// Get returns one site by id.
func (s *sSite) Get(ctx context.Context, id int) (*entity.NavSite, error) {
	var record *entity.NavSite
	if err := dao.NavSite.Ctx(ctx).Where(columns.Id, id).Scan(&record); err != nil {
		return nil, err
	}
	if record == nil {
		return nil, gerror.Newf("站点 %d 不存在", id)
	}
	return record, nil
}

// Create stores a new site and returns its id.
func (s *sSite) Create(ctx context.Context, input *model.SiteSaveInput) (int, error) {
	if err := s.checkCategory(ctx, input.CategoryId); err != nil {
		return 0, err
	}
	now := gtime.Now()
	id, err := dao.NavSite.Ctx(ctx).Data(do.NavSite{
		CategoryId:  input.CategoryId,
		Name:        strings.TrimSpace(input.Name),
		Url:         NormalizeUrl(input.Url),
		Icon:        strings.TrimSpace(input.Icon),
		Description: strings.TrimSpace(input.Description),
		Sort:        input.Sort,
		Status:      model.StatusValue(input.Status),
		Hits:        0,
		CreatedAt:   now,
		UpdatedAt:   now,
	}).InsertAndGetId()
	if err != nil {
		return 0, err
	}
	return int(id), nil
}

// Update modifies an existing site.
func (s *sSite) Update(ctx context.Context, id int, input *model.SiteSaveInput) error {
	existing, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if err = s.checkCategory(ctx, input.CategoryId); err != nil {
		return err
	}
	address := NormalizeUrl(input.Url)
	icon := strings.TrimSpace(input.Icon)
	data := g.Map{
		columns.CategoryId:  input.CategoryId,
		columns.Name:        strings.TrimSpace(input.Name),
		columns.Url:         address,
		columns.Icon:        icon,
		columns.Description: strings.TrimSpace(input.Description),
		columns.Sort:        input.Sort,
		columns.UpdatedAt:   gtime.Now(),
	}
	if input.Status != nil {
		data[columns.Status] = model.StatusValue(input.Status)
	}
	// The cached icon belongs to the previous address: drop it so that the next
	// page view downloads the icon of the new one.
	var staleKey string
	if address != existing.Url || icon != existing.Icon {
		staleKey = existing.IconKey
		data[columns.IconKey] = ""
		data[columns.IconFetchedAt] = nil
	}
	if _, err = dao.NavSite.Ctx(ctx).Where(columns.Id, id).Data(data).Update(); err != nil {
		return err
	}
	if staleKey != "" {
		return service.Icon().Release(ctx, staleKey)
	}
	return nil
}

// Delete removes a single site.
func (s *sSite) Delete(ctx context.Context, id int) error {
	existing, err := s.Get(ctx, id)
	if err != nil {
		return err
	}
	if _, err = dao.NavSite.Ctx(ctx).Where(columns.Id, id).Delete(); err != nil {
		return err
	}
	if existing.IconKey != "" {
		return service.Icon().Release(ctx, existing.IconKey)
	}
	return nil
}

// Hit records one visit and returns the visited site.
func (s *sSite) Hit(ctx context.Context, id int, includeRestricted bool) (*entity.NavSite, error) {
	record, err := s.Get(ctx, id)
	if err != nil {
		return nil, err
	}
	visible, err := s.IsVisible(ctx, record, includeRestricted)
	if err != nil {
		return nil, err
	}
	if !visible {
		// Do not reveal the address of a login-protected site.
		return nil, gerror.Newf("站点 %d 不存在", id)
	}
	if _, err = dao.NavSite.Ctx(ctx).Where(columns.Id, id).Increment(columns.Hits, 1); err != nil {
		return nil, err
	}
	record.Hits++
	return record, nil
}

// IsVisible reports whether a site may be disclosed to the caller.
func (s *sSite) IsVisible(ctx context.Context, site *entity.NavSite, includeRestricted bool) (bool, error) {
	if site == nil {
		return false, nil
	}
	if includeRestricted {
		return true, nil
	}
	restricted, err := service.Category().RestrictedIds(ctx)
	if err != nil {
		return false, err
	}
	for _, categoryId := range restricted {
		if categoryId == site.CategoryId {
			return false, nil
		}
	}
	return true, nil
}

// Search looks for sites matching the query keyword.
func (s *sSite) Search(ctx context.Context, query model.SiteQuery) (*model.SearchResult, error) {
	keyword := strings.TrimSpace(query.Keyword)
	result := &model.SearchResult{Keyword: keyword, Sites: make([]*model.SiteItem, 0)}
	if keyword == "" {
		return result, nil
	}
	limit := query.Limit
	if limit <= 0 || limit > 200 {
		limit = 50
	}
	m := dao.NavSite.Ctx(ctx)
	if query.OnlyEnabled {
		m = m.Where(columns.Status, consts.StatusOn)
	}
	like := "%" + keyword + "%"
	m = m.Where(
		"("+columns.Name+" LIKE ? OR "+columns.Description+" LIKE ? OR "+columns.Url+" LIKE ?)",
		like, like, like,
	)
	m, err := s.excludeRestricted(ctx, m, query.IncludeRestricted)
	if err != nil {
		return nil, err
	}
	var sites []*entity.NavSite
	if err = m.OrderDesc(columns.Hits).OrderAsc(columns.Sort).Limit(limit).Scan(&sites); err != nil {
		return nil, err
	}
	result.Sites = model.NewSiteItems(sites)
	result.Total = len(result.Sites)
	return result, nil
}

// checkCategory makes sure the referenced category exists.
func (s *sSite) checkCategory(ctx context.Context, categoryId int) error {
	if categoryId <= 0 {
		return gerror.New("请选择所属分类")
	}
	if _, err := service.Category().Get(ctx, categoryId); err != nil {
		return err
	}
	return nil
}

// NormalizeUrl trims an address and adds the https scheme when it is missing.
func NormalizeUrl(raw string) string {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ""
	}
	lower := strings.ToLower(raw)
	for _, scheme := range []string{"http://", "https://", "//", "mailto:", "tel:"} {
		if strings.HasPrefix(lower, scheme) {
			return raw
		}
	}
	return "https://" + raw
}
