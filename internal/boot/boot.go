// Package boot prepares the application on startup: it makes sure the SQLite
// data directory exists, applies the schema and seeds first-run data.
//
// Everything here is idempotent, so it is safe to run on every boot and on a
// completely empty deployment.
package boot

import (
	"context"
	_ "embed"
	"strings"

	// Register the pure-Go SQLite driver used by the default configuration.
	_ "github.com/gogf/gf/contrib/drivers/sqlite/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gfile"
	"github.com/gogf/gf/v2/os/gtime"
	"golang.org/x/crypto/bcrypt"

	"xyhelper-nav/internal/consts"
	"xyhelper-nav/internal/dao"
	"xyhelper-nav/internal/model/do"
)

// schemaSQL is the embedded database structure, see schema.sql.
//
//go:embed schema.sql
var schemaSQL string

// Init prepares database, schema and seed data.
func Init(ctx context.Context) error {
	if err := ensureDataDir(ctx); err != nil {
		return err
	}
	if err := migrate(ctx); err != nil {
		return err
	}
	if err := patchSchema(ctx); err != nil {
		return err
	}
	if err := seedSettings(ctx); err != nil {
		return err
	}
	if err := seedAdmin(ctx); err != nil {
		return err
	}
	return seedDemoData(ctx)
}

// schemaPatches lists columns that were introduced after the first release.
// CREATE TABLE IF NOT EXISTS cannot change a table that already exists, so
// these statements bring an existing database up to date.
var schemaPatches = []struct {
	Table  string
	Column string
	DDL    string
}{
	{
		Table:  "nav_category",
		Column: "visibility",
		DDL:    "ALTER TABLE `nav_category` ADD COLUMN `visibility` INTEGER NOT NULL DEFAULT 0",
	},
	{
		Table:  "nav_site",
		Column: "icon_key",
		DDL:    "ALTER TABLE `nav_site` ADD COLUMN `icon_key` TEXT NOT NULL DEFAULT ''",
	},
	{
		Table:  "nav_site",
		Column: "icon_fetched_at",
		DDL:    "ALTER TABLE `nav_site` ADD COLUMN `icon_fetched_at` DATETIME",
	},
}

// patchSchema adds every column that is missing from an existing database.
func patchSchema(ctx context.Context) error {
	db := g.DB()
	for _, patch := range schemaPatches {
		fields, err := db.TableFields(ctx, patch.Table)
		if err != nil {
			return err
		}
		if _, exists := fields[patch.Column]; exists {
			continue
		}
		if _, err = db.Exec(ctx, patch.DDL); err != nil {
			return err
		}
		g.Log().Infof(ctx, "schema updated: added column %s.%s", patch.Table, patch.Column)
	}
	return nil
}

// ensureDataDir creates the directory of the SQLite file when it does not exist.
// A missing configuration file is not fatal: the database may also be
// configured programmatically, in which case its directory already exists.
func ensureDataDir(ctx context.Context) error {
	linkValue, err := g.Cfg().Get(ctx, "database.default.link")
	if err != nil {
		return nil
	}
	link := linkValue.String()
	start := strings.Index(link, "@file(")
	if start < 0 {
		return nil
	}
	path := link[start+len("@file("):]
	if end := strings.Index(path, ")"); end >= 0 {
		path = path[:end]
	}
	if path == "" || path == ":memory:" {
		return nil
	}
	dir := gfile.Dir(path)
	if dir == "" || gfile.Exists(dir) {
		return nil
	}
	return gfile.Mkdir(dir)
}

// migrate executes every statement of the embedded schema file.
func migrate(ctx context.Context) error {
	db := g.DB()
	for _, statement := range splitStatements(schemaSQL) {
		if _, err := db.Exec(ctx, statement); err != nil {
			return err
		}
	}
	return nil
}

// splitStatements strips SQL comments and splits a script into single statements.
func splitStatements(script string) []string {
	lines := make([]string, 0, 64)
	for _, line := range strings.Split(script, "\n") {
		if strings.HasPrefix(strings.TrimSpace(line), "--") {
			continue
		}
		lines = append(lines, line)
	}
	statements := make([]string, 0, 16)
	for _, part := range strings.Split(strings.Join(lines, "\n"), ";") {
		if statement := strings.TrimSpace(part); statement != "" {
			statements = append(statements, statement)
		}
	}
	return statements
}

// seedSettings inserts every missing default setting.
func seedSettings(ctx context.Context) error {
	for key, value := range consts.DefaultSettings {
		count, err := dao.NavSetting.Ctx(ctx).Where(dao.NavSetting.Columns().K, key).Count()
		if err != nil {
			return err
		}
		if count > 0 {
			continue
		}
		now := gtime.Now()
		if _, err = dao.NavSetting.Ctx(ctx).Data(do.NavSetting{
			K:         key,
			V:         value,
			CreatedAt: now,
			UpdatedAt: now,
		}).Insert(); err != nil {
			return err
		}
	}
	return nil
}

// seedAdmin creates the default administrator when the table is empty.
func seedAdmin(ctx context.Context) error {
	count, err := dao.NavAdmin.Ctx(ctx).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		return nil
	}
	hashed, err := bcrypt.GenerateFromPassword([]byte(consts.DefaultAdminPass), bcrypt.DefaultCost)
	if err != nil {
		return err
	}
	now := gtime.Now()
	_, err = dao.NavAdmin.Ctx(ctx).Data(do.NavAdmin{
		Username:  consts.DefaultAdminUser,
		Password:  string(hashed),
		Nickname:  "管理员",
		CreatedAt: now,
		UpdatedAt: now,
	}).Insert()
	return err
}

// seedDemoData fills a brand new installation with a small, useful example
// catalogue so the navigation page never starts out empty.
func seedDemoData(ctx context.Context) error {
	categoryCount, err := dao.NavCategory.Ctx(ctx).Count()
	if err != nil {
		return err
	}
	siteCount, err := dao.NavSite.Ctx(ctx).Count()
	if err != nil {
		return err
	}
	if categoryCount > 0 || siteCount > 0 {
		return nil
	}

	now := gtime.Now()
	for categoryIndex, category := range demoCategories {
		categoryId, err := dao.NavCategory.Ctx(ctx).Data(do.NavCategory{
			Name:        category.Name,
			Icon:        category.Icon,
			Description: category.Description,
			Sort:        (categoryIndex + 1) * 10,
			Status:      consts.StatusOn,
			Visibility:  consts.VisibilityPublic,
			CreatedAt:   now,
			UpdatedAt:   now,
		}).InsertAndGetId()
		if err != nil {
			return err
		}
		for siteIndex, site := range category.Sites {
			if _, err = dao.NavSite.Ctx(ctx).Data(do.NavSite{
				CategoryId:  int(categoryId),
				Name:        site.Name,
				Url:         site.Url,
				Icon:        site.Icon,
				Description: site.Description,
				Sort:        (siteIndex + 1) * 10,
				Status:      consts.StatusOn,
				CreatedAt:   now,
				UpdatedAt:   now,
			}).Insert(); err != nil {
				return err
			}
		}
	}
	return nil
}
