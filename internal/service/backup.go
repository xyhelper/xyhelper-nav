package service

import (
	"context"

	"xyhelper-nav/internal/model"
)

// IBackup imports and exports the navigation content as a portable document.
type IBackup interface {
	// Export returns the current content as a backup document.
	Export(ctx context.Context) (*model.BackupData, error)
	// Import replaces the current content with the given document.
	Import(ctx context.Context, data *model.BackupData) error
}

var localBackup IBackup

// Backup returns the backup service implementation.
func Backup() IBackup {
	if localBackup == nil {
		panic("implementation of IBackup is not registered, did you forget to import internal/logic?")
	}
	return localBackup
}

// RegisterBackup registers the backup service implementation.
func RegisterBackup(impl IBackup) {
	localBackup = impl
}
