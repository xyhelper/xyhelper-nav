package model

import (
	"testing"
	"time"

	"xyhelper-nav/internal/consts"
	"xyhelper-nav/internal/model/entity"
)

func TestStatusValue(t *testing.T) {
	on := consts.StatusOn
	off := consts.StatusOff
	cases := []struct {
		name   string
		status *int
		want   int
	}{
		{name: "missing status defaults to visible", status: nil, want: consts.StatusOn},
		{name: "explicit on", status: &on, want: consts.StatusOn},
		{name: "explicit off", status: &off, want: consts.StatusOff},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := StatusValue(c.status); got != c.want {
				t.Fatalf("StatusValue(%v) = %d, want %d", c.status, got, c.want)
			}
		})
	}
}

func TestVisibilityValue(t *testing.T) {
	login := consts.VisibilityLogin
	public := consts.VisibilityPublic
	cases := []struct {
		name       string
		visibility *int
		want       int
	}{
		{name: "missing visibility defaults to public", visibility: nil, want: consts.VisibilityPublic},
		{name: "explicit public", visibility: &public, want: consts.VisibilityPublic},
		{name: "explicit login", visibility: &login, want: consts.VisibilityLogin},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			if got := VisibilityValue(c.visibility); got != c.want {
				t.Fatalf("VisibilityValue(%v) = %d, want %d", c.visibility, got, c.want)
			}
		})
	}
}

func TestNewSiteItem(t *testing.T) {
	if NewSiteItem(nil) != nil {
		t.Fatal("NewSiteItem(nil) should return nil")
	}
	now := time.Now()
	record := &entity.NavSite{
		Id:          7,
		CategoryId:  3,
		Name:        "GitHub",
		Url:         "https://github.com",
		Description: "code hosting",
		Sort:        20,
		Status:      consts.StatusOn,
		Hits:        42,
		CreatedAt:   now,
	}
	item := NewSiteItem(record)
	if item == nil {
		t.Fatal("NewSiteItem returned nil for a valid record")
	}
	if item.Id != record.Id || item.CategoryId != record.CategoryId || item.Name != record.Name {
		t.Fatalf("unexpected identity mapping: %+v", item)
	}
	if item.Url != record.Url || item.Hits != record.Hits || !item.CreatedAt.Equal(now) {
		t.Fatalf("unexpected field mapping: %+v", item)
	}
}

func TestNewSiteItems(t *testing.T) {
	items := NewSiteItems(nil)
	if items == nil || len(items) != 0 {
		t.Fatalf("NewSiteItems(nil) = %#v, want empty non-nil slice", items)
	}
	items = NewSiteItems([]*entity.NavSite{{Id: 1}, {Id: 2}})
	if len(items) != 2 || items[1].Id != 2 {
		t.Fatalf("unexpected slice conversion: %+v", items)
	}
}

func TestNewCategoryItem(t *testing.T) {
	if NewCategoryItem(nil) != nil {
		t.Fatal("NewCategoryItem(nil) should return nil")
	}
	item := NewCategoryItem(&entity.NavCategory{Id: 5, Name: "开发工具", Status: consts.StatusOn})
	if item.Id != 5 || item.Name != "开发工具" {
		t.Fatalf("unexpected mapping: %+v", item)
	}
	if item.Sites == nil {
		t.Fatal("Sites should be initialised to an empty slice, not nil")
	}
}
