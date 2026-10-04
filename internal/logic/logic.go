// Package logic registers every business implementation of the application.
// Importing this package wires the internal/service interfaces to the
// internal/logic implementations through their init functions.
package logic

import (
	_ "xyhelper-nav/internal/logic/admin"
	_ "xyhelper-nav/internal/logic/backup"
	_ "xyhelper-nav/internal/logic/category"
	_ "xyhelper-nav/internal/logic/icon"
	_ "xyhelper-nav/internal/logic/setting"
	_ "xyhelper-nav/internal/logic/site"
)
