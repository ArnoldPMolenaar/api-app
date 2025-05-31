package services

import (
	"api-app/main/src/cache"
	"api-app/main/src/database"
	"api-app/main/src/enums"
	"api-app/main/src/models"
	"context"
	"encoding/json"
	"fmt"
	"github.com/valkey-io/valkey-go"
	"os"
	"time"
)

// GetAppSettingsByName method to get settings by app name.
func GetAppSettingsByName(appName string, level enums.Level) (*[]models.AppSetting, error) {
	var settings []models.AppSetting
	cacheKey := AppSettingsCacheKeyOnName(appName, level)

	if inCache, err := IsAppSettingsInCache(cacheKey); err != nil {
		return nil, err
	} else if inCache {
		if cacheSettings, err := GetAppSettingsFromCache(cacheKey); err != nil {
			return nil, err
		} else if cacheSettings != nil && len(*cacheSettings) > 0 {
			settings = *cacheSettings
		}
	}

	if len(settings) == 0 {
		if result := database.Pg.Model(&models.AppSetting{}).
			Joins("JOIN apps ON apps.id = app_settings.app_id").
			Where("apps.name = ? AND (level = 'both' OR level = ?)", appName, level.String()).
			Find(&settings); result.Error != nil {
			return nil, result.Error
		}
		_ = SetAppSettingsToCache(cacheKey, &settings)
	}

	return &settings, nil
}

// GetAppSettingsByAppID method to get settings by app ID.
func GetAppSettingsByAppID(appID uint, level enums.Level) (*[]models.AppSetting, error) {
	var settings []models.AppSetting
	cacheKey := AppSettingsCacheKeyOnId(appID, level)

	if inCache, err := IsAppSettingsInCache(cacheKey); err != nil {
		return nil, err
	} else if inCache {
		if cacheSettings, err := GetAppSettingsFromCache(cacheKey); err != nil {
			return nil, err
		} else if cacheSettings != nil && len(*cacheSettings) > 0 {
			settings = *cacheSettings
		}
	}

	if len(settings) == 0 {
		if result := database.Pg.Model(&models.AppSetting{}).
			Where("app_id = ? AND (level = 'both' OR level = ?)", appID, level.String()).
			Find(&settings); result.Error != nil {
			return nil, result.Error
		}
		_ = SetAppSettingsToCache(cacheKey, &settings)
	}

	return &settings, nil
}

// IsAppSettingsInCache checks if the settings exists in the cache.
func IsAppSettingsInCache(key string) (bool, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Exists().Key(key).Build())
	if result.Error() != nil {
		return false, result.Error()
	}

	value, err := result.ToInt64()
	if err != nil {
		return false, err
	}

	return value == 1, nil
}

// GetAppSettingsFromCache gets the settings from the cache.
func GetAppSettingsFromCache(key string) (*[]models.AppSetting, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Get().Key(key).Build())
	if result.Error() != nil {
		return nil, result.Error()
	}

	value, err := result.ToString()
	if err != nil {
		return nil, err
	}

	var settings []models.AppSetting
	if err := json.Unmarshal([]byte(value), &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

// SetAppSettingsToCache sets the settings to the cache.
func SetAppSettingsToCache(key string, settings *[]models.AppSetting) error {
	value, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	expiration := os.Getenv("VALKEY_EXPIRATION")
	duration, err := time.ParseDuration(expiration)
	if err != nil {
		return err
	}

	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Set().Key(key).Value(valkey.BinaryString(value)).Ex(duration).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// DeleteAppSettingsFromCache deletes an existing setting from the cache.
func DeleteAppSettingsFromCache(key string) error {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Del().Key(key).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// AppSettingsCacheKeyOnName returns the key for the settings cache with a name.
func AppSettingsCacheKeyOnName(appName string, level enums.Level) string {
	return fmt.Sprintf("%s:settings:%s", appName, level.String())
}

// AppSettingsCacheKeyOnId returns the key for the settings cache with an id.
func AppSettingsCacheKeyOnId(appID uint, level enums.Level) string {
	return fmt.Sprintf("settings:apps:%d:%s", appID, level.String())
}
