package admin

import (
	"context"
	"strings"

	v1 "xyhelper-nav/api/nav/v1"
	"xyhelper-nav/internal/model"
	"xyhelper-nav/internal/service"
)

// SiteList returns the navigation entries matching the given filters. The
// administration panel always sees login-protected categories as well.
func (c *Controller) SiteList(ctx context.Context, req *v1.SiteListReq) (res *v1.SiteListRes, err error) {
	records, err := service.Site().List(ctx, model.SiteQuery{
		CategoryId:        req.CategoryId,
		IncludeRestricted: true,
	})
	if err != nil {
		return nil, err
	}
	list := model.NewSiteItems(records)
	if keyword := strings.ToLower(strings.TrimSpace(req.Keyword)); keyword != "" {
		filtered := make([]*model.SiteItem, 0, len(list))
		for _, item := range list {
			if strings.Contains(strings.ToLower(item.Name), keyword) ||
				strings.Contains(strings.ToLower(item.Url), keyword) {
				filtered = append(filtered, item)
			}
		}
		list = filtered
	}
	return &v1.SiteListRes{List: list, Total: len(list)}, nil
}

// SiteCreate stores a new navigation entry.
func (c *Controller) SiteCreate(ctx context.Context, req *v1.SiteCreateReq) (res *v1.SiteCreateRes, err error) {
	id, err := service.Site().Create(ctx, &model.SiteSaveInput{
		CategoryId:  req.CategoryId,
		Name:        req.Name,
		Url:         req.Url,
		Icon:        req.Icon,
		Description: req.Description,
		Sort:        req.Sort,
		Status:      req.Status,
	})
	if err != nil {
		return nil, err
	}
	return &v1.SiteCreateRes{Id: id}, nil
}

// SiteUpdate modifies an existing navigation entry.
func (c *Controller) SiteUpdate(ctx context.Context, req *v1.SiteUpdateReq) (res *v1.SiteUpdateRes, err error) {
	if err = service.Site().Update(ctx, req.Id, &model.SiteSaveInput{
		CategoryId:  req.CategoryId,
		Name:        req.Name,
		Url:         req.Url,
		Icon:        req.Icon,
		Description: req.Description,
		Sort:        req.Sort,
		Status:      req.Status,
	}); err != nil {
		return nil, err
	}
	return &v1.SiteUpdateRes{}, nil
}

// SiteDelete removes a navigation entry.
func (c *Controller) SiteDelete(ctx context.Context, req *v1.SiteDeleteReq) (res *v1.SiteDeleteRes, err error) {
	if err = service.Site().Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.SiteDeleteRes{}, nil
}

// IconRefresh re-downloads the icons of one site or of every site.
func (c *Controller) IconRefresh(ctx context.Context, req *v1.IconRefreshReq) (res *v1.IconRefreshRes, err error) {
	result, err := service.Icon().Refresh(ctx, req.Id)
	if err != nil {
		return nil, err
	}
	return &v1.IconRefreshRes{
		Total:     result.Total,
		Refreshed: result.Refreshed,
		Failed:    result.Failed,
	}, nil
}
