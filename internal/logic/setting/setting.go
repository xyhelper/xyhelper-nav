// Package setting implements the site configuration service.
package setting

import (
	"context"
	"strings"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"

	"xyhelper-nav/internal/dao"
	"xyhelper-nav/internal/model/do"
	"xyhelper-nav/internal/model/entity"
	"xyhelper-nav/internal/service"
)

type sSetting struct{}

func init() {
	service.RegisterSetting(New())
}

// New creates and returns the setting service implementation.
func New() service.ISetting {
	return &sSetting{}
}

// All returns every setting as a key/value map.
func (s *sSetting) All(ctx context.Context) (map[string]string, error) {
	var records []*entity.NavSetting
	if err := dao.NavSetting.Ctx(ctx).Scan(&records); err != nil {
		return nil, err
	}
	result := make(map[string]string, len(records))
	for _, record := range records {
		result[record.K] = record.V
	}
	return result, nil
}

// Get returns a single setting value, or the optional default when missing.
func (s *sSetting) Get(ctx context.Context, key string, defaultValue ...string) (string, error) {
	record, err := dao.NavSetting.Ctx(ctx).Where(dao.NavSetting.Columns().K, key).One()
	if err != nil {
		return "", err
	}
	if record.IsEmpty() {
		if len(defaultValue) > 0 {
			return defaultValue[0], nil
		}
		return "", nil
	}
	return record[dao.NavSetting.Columns().V].String(), nil
}

// GetBool returns a setting interpreted as a boolean flag.
func (s *sSetting) GetBool(ctx context.Context, key string, defaultValue ...bool) (bool, error) {
	fallback := false
	if len(defaultValue) > 0 {
		fallback = defaultValue[0]
	}
	value, err := s.Get(ctx, key)
	if err != nil {
		return fallback, err
	}
	switch strings.ToLower(strings.TrimSpace(value)) {
	case "1", "true", "on", "yes", "enable", "enabled":
		return true, nil
	case "0", "false", "off", "no", "disable", "disabled":
		return false, nil
	case "":
		return fallback, nil
	default:
		return fallback, nil
	}
}

// Set writes a single setting.
func (s *sSetting) Set(ctx context.Context, key, value string) error {
	if strings.TrimSpace(key) == "" {
		return gerror.New("配置项名称不能为空")
	}
	count, err := dao.NavSetting.Ctx(ctx).Where(dao.NavSetting.Columns().K, key).Count()
	if err != nil {
		return err
	}
	if count > 0 {
		_, err = dao.NavSetting.Ctx(ctx).
			Where(dao.NavSetting.Columns().K, key).
			Data(g.Map{dao.NavSetting.Columns().V: value, dao.NavSetting.Columns().UpdatedAt: gtime.Now()}).
			Update()
		return err
	}
	now := gtime.Now()
	_, err = dao.NavSetting.Ctx(ctx).Data(do.NavSetting{
		K:         key,
		V:         value,
		CreatedAt: now,
		UpdatedAt: now,
	}).Insert()
	return err
}

// SetMap writes multiple settings at once.
func (s *sSetting) SetMap(ctx context.Context, data map[string]string) error {
	for key, value := range data {
		if err := s.Set(ctx, key, value); err != nil {
			return err
		}
	}
	return nil
}
