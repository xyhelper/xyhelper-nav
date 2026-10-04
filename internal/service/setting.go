// Package service defines the business interfaces of the application.
// Implementations live in the internal/logic packages and register themselves
// during package initialisation.
package service

import (
	"context"
)

// ISetting manages the key/value site configuration.
type ISetting interface {
	// All returns every setting as a key/value map.
	All(ctx context.Context) (map[string]string, error)
	// Get returns a single setting value, or the optional default when missing.
	Get(ctx context.Context, key string, defaultValue ...string) (string, error)
	// GetBool returns a setting interpreted as a boolean flag.
	GetBool(ctx context.Context, key string, defaultValue ...bool) (bool, error)
	// Set writes a single setting.
	Set(ctx context.Context, key, value string) error
	// SetMap writes multiple settings at once.
	SetMap(ctx context.Context, data map[string]string) error
}

var localSetting ISetting

// Setting returns the setting service implementation.
func Setting() ISetting {
	if localSetting == nil {
		panic("implementation of ISetting is not registered, did you forget to import internal/logic?")
	}
	return localSetting
}

// RegisterSetting registers the setting service implementation.
func RegisterSetting(impl ISetting) {
	localSetting = impl
}
