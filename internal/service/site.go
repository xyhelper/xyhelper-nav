package service

import (
	"context"

	"xyhelper-nav/internal/model"
	"xyhelper-nav/internal/model/entity"
)

// ISite manages navigation entries.
type ISite interface {
	// List returns the sites matching the query. Sites that belong to a
	// login-protected category are only returned when the query asks for them.
	List(ctx context.Context, query model.SiteQuery) ([]*entity.NavSite, error)
	// Get returns one site by id.
	Get(ctx context.Context, id int) (*entity.NavSite, error)
	// Create stores a new site and returns its id.
	Create(ctx context.Context, input *model.SiteSaveInput) (int, error)
	// Update modifies an existing site.
	Update(ctx context.Context, id int, input *model.SiteSaveInput) error
	// Delete removes a single site.
	Delete(ctx context.Context, id int) error
	// Hit records one visit and returns the visited site. Sites inside a
	// login-protected category are rejected unless includeRestricted is set.
	Hit(ctx context.Context, id int, includeRestricted bool) (*entity.NavSite, error)
	// IsVisible reports whether a site may be disclosed to the caller. It is
	// the single place where the category visibility rule is evaluated for a
	// single record.
	IsVisible(ctx context.Context, site *entity.NavSite, includeRestricted bool) (bool, error)
	// Search looks for sites matching the query keyword.
	Search(ctx context.Context, query model.SiteQuery) (*model.SearchResult, error)
}

var localSite ISite

// Site returns the site service implementation.
func Site() ISite {
	if localSite == nil {
		panic("implementation of ISite is not registered, did you forget to import internal/logic?")
	}
	return localSite
}

// RegisterSite registers the site service implementation.
func RegisterSite(impl ISite) {
	localSite = impl
}
