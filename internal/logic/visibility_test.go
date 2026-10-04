// Package logic_test verifies the cross-service rules of the application
// against a throw-away SQLite database.
package logic_test

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/gogf/gf/v2/database/gdb"

	"xyhelper-nav/internal/boot"
	"xyhelper-nav/internal/consts"
	// Register every business implementation.
	_ "xyhelper-nav/internal/logic"
	"xyhelper-nav/internal/model"
	"xyhelper-nav/internal/service"
)

// TestMain points the default database group at a temporary SQLite file,
// moves the working directory inside that temporary directory (so relative
// runtime paths such as the icon cache stay out of the repository) and
// bootstraps the schema.
func TestMain(m *testing.M) {
	dir, err := os.MkdirTemp("", "xyhelper-nav-logic")
	if err != nil {
		panic(err)
	}
	origin, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	if err = os.Chdir(dir); err != nil {
		panic(err)
	}
	err = gdb.SetConfig(gdb.Config{
		"default": gdb.ConfigGroup{
			gdb.ConfigNode{
				Type: "sqlite",
				Link: "sqlite::@file(" + filepath.Join(dir, "test.db") + ")",
			},
		},
	})
	if err != nil {
		panic(err)
	}
	if err = boot.Init(context.Background()); err != nil {
		panic(err)
	}
	code := m.Run()
	if err = os.Chdir(origin); err != nil {
		panic(err)
	}
	if err = os.RemoveAll(dir); err != nil {
		panic(err)
	}
	os.Exit(code)
}

// TestCategoryVisibilityHiddenFromGuests covers the "登录后可见" rule: a
// restricted category, its sites, the search index and the hit counter must all
// stay invisible until the visitor is authenticated.
func TestCategoryVisibilityHiddenFromGuests(t *testing.T) {
	var (
		ctx        = context.Background()
		statusOn   = consts.StatusOn
		restricted = consts.VisibilityLogin
		guest      = model.CategoryQuery{OnlyEnabled: true}
		admin      = model.CategoryQuery{OnlyEnabled: true, IncludeRestricted: true}
	)

	categoryId, err := service.Category().Create(ctx, &model.CategorySaveInput{
		Name:       "私有收藏",
		Status:     &statusOn,
		Visibility: &restricted,
	})
	if err != nil {
		t.Fatalf("create restricted category: %v", err)
	}
	siteId, err := service.Site().Create(ctx, &model.SiteSaveInput{
		CategoryId: categoryId,
		Name:       "Private Wiki",
		Url:        "https://private.example.com",
		Status:     &statusOn,
	})
	if err != nil {
		t.Fatalf("create site in restricted category: %v", err)
	}

	// The category listing must hide it from guests and show it to admins.
	guestCategories, err := service.Category().List(ctx, guest)
	if err != nil {
		t.Fatalf("list categories as guest: %v", err)
	}
	if containsCategory(guestCategories, "私有收藏") {
		t.Error("restricted category is visible to a guest")
	}
	adminCategories, err := service.Category().List(ctx, admin)
	if err != nil {
		t.Fatalf("list categories as admin: %v", err)
	}
	if !containsCategory(adminCategories, "私有收藏") {
		t.Error("restricted category is not visible to an admin")
	}

	// Search must not disclose the site of a restricted category.
	guestSearch, err := service.Site().Search(ctx, model.SiteQuery{
		Keyword:     "Private Wiki",
		OnlyEnabled: true,
	})
	if err != nil {
		t.Fatalf("search as guest: %v", err)
	}
	if guestSearch.Total != 0 {
		t.Errorf("guest search returned %d protected site(s), want 0", guestSearch.Total)
	}
	adminSearch, err := service.Site().Search(ctx, model.SiteQuery{
		Keyword:           "Private Wiki",
		OnlyEnabled:       true,
		IncludeRestricted: true,
	})
	if err != nil {
		t.Fatalf("search as admin: %v", err)
	}
	if adminSearch.Total != 1 {
		t.Errorf("admin search returned %d site(s), want 1", adminSearch.Total)
	}

	// The hit endpoint must not leak the address of a protected site.
	if _, err = service.Site().Hit(ctx, siteId, false); err == nil {
		t.Error("guest was allowed to register a hit on a protected site")
	}
	if _, err = service.Site().Hit(ctx, siteId, true); err != nil {
		t.Errorf("admin was not allowed to register a hit: %v", err)
	}

	// Flipping the category back to public must expose it again.
	public := consts.VisibilityPublic
	if err = service.Category().Update(ctx, categoryId, &model.CategorySaveInput{
		Name:       "私有收藏",
		Status:     &statusOn,
		Visibility: &public,
	}); err != nil {
		t.Fatalf("publish category: %v", err)
	}
	guestCategories, err = service.Category().List(ctx, guest)
	if err != nil {
		t.Fatalf("list categories after publish: %v", err)
	}
	if !containsCategory(guestCategories, "私有收藏") {
		t.Error("published category is still hidden from a guest")
	}
}

// containsCategory reports whether the listing holds a category with that name.
func containsCategory(categories []*model.CategoryItem, name string) bool {
	for _, category := range categories {
		if category.Name == name {
			return true
		}
	}
	return false
}
