package service

import (
	"context"

	"xyhelper-nav/internal/model"
	"xyhelper-nav/internal/model/entity"
)

// IIcon downloads remote site icons once and stores them inside the project,
// so that a browser only ever requests images from this server.
type IIcon interface {
	// Ensure makes sure the icon of a site is cached locally and returns its
	// cache key. An empty key means no icon is available (yet).
	Ensure(ctx context.Context, site *entity.NavSite) (string, error)
	// Refresh re-downloads icons: a single site when siteId > 0, every site
	// when siteId is 0.
	Refresh(ctx context.Context, siteId int) (*model.IconRefreshResult, error)
	// FilePath resolves a cache key to an absolute file path. It returns an
	// empty string for keys that are malformed or point outside the storage.
	FilePath(ctx context.Context, key string) string
	// Release deletes a cached file once no site references it anymore.
	Release(ctx context.Context, key string) error
}

var localIcon IIcon

// Icon returns the icon cache service implementation.
func Icon() IIcon {
	if localIcon == nil {
		panic("implementation of IIcon is not registered, did you forget to import internal/logic?")
	}
	return localIcon
}

// RegisterIcon registers the icon cache service implementation.
func RegisterIcon(impl IIcon) {
	localIcon = impl
}
