// Package nav implements the public (front-end) HTTP controllers.
package nav

import (
	"context"

	v1 "xyhelper-nav/api/nav/v1"
	"xyhelper-nav/internal/middleware"
	"xyhelper-nav/internal/model"
	"xyhelper-nav/internal/service"
)

// Controller serves the public navigation API.
type Controller struct{}

// New creates and returns the public navigation controller.
func New() *Controller {
	return &Controller{}
}

// Index returns the whole payload required to render the navigation page.
// Categories marked as "登录后可见" are only included for a signed-in
// administrator.
func (c *Controller) Index(ctx context.Context, req *v1.IndexReq) (res *v1.IndexRes, err error) {
	settings, err := service.Setting().All(ctx)
	if err != nil {
		return nil, err
	}
	categories, err := service.Category().List(ctx, model.CategoryQuery{
		OnlyEnabled:       true,
		IncludeRestricted: middleware.IsAdmin(ctx),
	})
	if err != nil {
		return nil, err
	}
	total := 0
	for _, category := range categories {
		total += len(category.Sites)
	}
	return &v1.IndexRes{
		Settings:   settings,
		Categories: categories,
		Total:      total,
	}, nil
}

// Hit records one visit of a navigation entry and returns its address.
func (c *Controller) Hit(ctx context.Context, req *v1.HitReq) (res *v1.HitRes, err error) {
	record, err := service.Site().Hit(ctx, req.Id, middleware.IsAdmin(ctx))
	if err != nil {
		return nil, err
	}
	return &v1.HitRes{
		Id:   record.Id,
		Url:  record.Url,
		Hits: record.Hits,
	}, nil
}

// Search looks for navigation entries matching a keyword.
func (c *Controller) Search(ctx context.Context, req *v1.SearchReq) (res *v1.SearchRes, err error) {
	result, err := service.Site().Search(ctx, model.SiteQuery{
		Keyword:           req.Keyword,
		Limit:             req.Limit,
		OnlyEnabled:       true,
		IncludeRestricted: middleware.IsAdmin(ctx),
	})
	if err != nil {
		return nil, err
	}
	return &v1.SearchRes{
		Keyword: result.Keyword,
		Total:   result.Total,
		Sites:   result.Sites,
	}, nil
}
