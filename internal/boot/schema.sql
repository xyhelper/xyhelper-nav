-- ----------------------------------------------------------------------------
-- xyhelper-nav schema (SQLite)
-- This file is the single source of truth for the database structure.
-- It is embedded into the binary and executed automatically on startup,
-- so a fresh deployment needs no manual migration step.
-- ----------------------------------------------------------------------------

CREATE TABLE IF NOT EXISTS `nav_category`
(
    `id`          INTEGER PRIMARY KEY AUTOINCREMENT,
    `name`        TEXT    NOT NULL DEFAULT '',
    `icon`        TEXT    NOT NULL DEFAULT '',
    `description` TEXT    NOT NULL DEFAULT '',
    `sort`        INTEGER NOT NULL DEFAULT 0,
    `status`      INTEGER NOT NULL DEFAULT 1,
    -- 0: visible to everyone, 1: visible to signed-in administrators only.
    `visibility`  INTEGER NOT NULL DEFAULT 0,
    `created_at`  DATETIME,
    `updated_at`  DATETIME
);

CREATE INDEX IF NOT EXISTS `idx_nav_category_sort` ON `nav_category` (`sort`);

CREATE TABLE IF NOT EXISTS `nav_site`
(
    `id`          INTEGER PRIMARY KEY AUTOINCREMENT,
    `category_id` INTEGER NOT NULL DEFAULT 0,
    `name`        TEXT    NOT NULL DEFAULT '',
    `url`         TEXT    NOT NULL DEFAULT '',
    `icon`        TEXT    NOT NULL DEFAULT '',
    `description` TEXT    NOT NULL DEFAULT '',
    `sort`        INTEGER NOT NULL DEFAULT 0,
    `status`      INTEGER NOT NULL DEFAULT 1,
    `hits`        INTEGER NOT NULL DEFAULT 0,
    -- Local icon cache: `icon_key` is the file name inside the icon storage
    -- directory, `icon_fetched_at` is the last download attempt.
    `icon_key`        TEXT NOT NULL DEFAULT '',
    `icon_fetched_at` DATETIME,
    `created_at`  DATETIME,
    `updated_at`  DATETIME
);

CREATE INDEX IF NOT EXISTS `idx_nav_site_category` ON `nav_site` (`category_id`);
CREATE INDEX IF NOT EXISTS `idx_nav_site_sort` ON `nav_site` (`sort`);

CREATE TABLE IF NOT EXISTS `nav_setting`
(
    `id`         INTEGER PRIMARY KEY AUTOINCREMENT,
    `k`          TEXT NOT NULL DEFAULT '',
    `v`          TEXT NOT NULL DEFAULT '',
    `created_at` DATETIME,
    `updated_at` DATETIME
);

CREATE UNIQUE INDEX IF NOT EXISTS `idx_nav_setting_k` ON `nav_setting` (`k`);

CREATE TABLE IF NOT EXISTS `nav_admin`
(
    `id`            INTEGER PRIMARY KEY AUTOINCREMENT,
    `username`      TEXT NOT NULL DEFAULT '',
    `password`      TEXT NOT NULL DEFAULT '',
    `nickname`      TEXT NOT NULL DEFAULT '',
    `last_login_at` DATETIME,
    `last_login_ip` TEXT NOT NULL DEFAULT '',
    `created_at`    DATETIME,
    `updated_at`    DATETIME
);

CREATE UNIQUE INDEX IF NOT EXISTS `idx_nav_admin_username` ON `nav_admin` (`username`);
