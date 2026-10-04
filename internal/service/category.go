package service

import (
	"context"

	"xyhelper-nav/internal/model"
	"xyhelper-nav/internal/model/entity"
)

// ICategory manages navigation categories.
type ICategory interface {
	// List returns every category the caller may see, optionally with its
	// sites attached.
	List(ctx context.Context, query model.CategoryQuery) ([]*model.CategoryItem, error)
	// ListEntities returns the raw category records the caller may see.
	ListEntities(ctx context.Context, query model.CategoryQuery) ([]*entity.NavCategory, error)
	// RestrictedIds returns the ids of the categories that require a login.
	// It is the single source of truth for the "登录后可见" rule.
	RestrictedIds(ctx context.Context) ([]int, error)
	// Get returns one category by id.
	Get(ctx context.Context, id int) (*entity.NavCategory, error)
	// Create stores a new category and returns its id.
	Create(ctx context.Context, input *model.CategorySaveInput) (int, error)
	// Update modifies an existing category.
	Update(ctx context.Context, id int, input *model.CategorySaveInput) error
	// Delete removes a category and every site that belongs to it.
	Delete(ctx context.Context, id int) error
}

var localCategory ICategory

// Category returns the category service implementation.
func Category() ICategory {
	if localCategory == nil {
		panic("implementation of ICategory is not registered, did you forget to import internal/logic?")
	}
	return localCategory
}

// RegisterCategory registers the category service implementation.
func RegisterCategory(impl ICategory) {
	localCategory = impl
}
