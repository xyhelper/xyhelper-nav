package admin

import (
	"context"

	v1 "xyhelper-nav/api/nav/v1"
	"xyhelper-nav/internal/model"
	"xyhelper-nav/internal/service"
)

// CategoryList returns every category together with its sites, including the
// ones that are only visible after a login.
func (c *Controller) CategoryList(ctx context.Context, req *v1.CategoryListReq) (res *v1.CategoryListRes, err error) {
	list, err := service.Category().List(ctx, model.CategoryQuery{IncludeRestricted: true})
	if err != nil {
		return nil, err
	}
	return &v1.CategoryListRes{List: list, Total: len(list)}, nil
}

// CategoryCreate stores a new category.
func (c *Controller) CategoryCreate(ctx context.Context, req *v1.CategoryCreateReq) (res *v1.CategoryCreateRes, err error) {
	id, err := service.Category().Create(ctx, &model.CategorySaveInput{
		Name:        req.Name,
		Icon:        req.Icon,
		Description: req.Description,
		Sort:        req.Sort,
		Status:      req.Status,
		Visibility:  req.Visibility,
	})
	if err != nil {
		return nil, err
	}
	return &v1.CategoryCreateRes{Id: id}, nil
}

// CategoryUpdate modifies an existing category.
func (c *Controller) CategoryUpdate(ctx context.Context, req *v1.CategoryUpdateReq) (res *v1.CategoryUpdateRes, err error) {
	if err = service.Category().Update(ctx, req.Id, &model.CategorySaveInput{
		Name:        req.Name,
		Icon:        req.Icon,
		Description: req.Description,
		Sort:        req.Sort,
		Status:      req.Status,
		Visibility:  req.Visibility,
	}); err != nil {
		return nil, err
	}
	return &v1.CategoryUpdateRes{}, nil
}

// CategoryDelete removes a category and all of its sites.
func (c *Controller) CategoryDelete(ctx context.Context, req *v1.CategoryDeleteReq) (res *v1.CategoryDeleteRes, err error) {
	if err = service.Category().Delete(ctx, req.Id); err != nil {
		return nil, err
	}
	return &v1.CategoryDeleteRes{}, nil
}
