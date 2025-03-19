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

// GetSettingsByDomainName method to get settings by domain name.
func GetSettingsByDomainName(appName, domainName string, level enums.Level) (*[]models.DomainSetting, error) {
	var settings []models.DomainSetting
	cacheKey := SettingsCacheKeyOnName(appName, domainName)

	if inCache, err := IsSettingsInCache(cacheKey); err != nil {
		return nil, err
	} else if inCache {
		if cacheSettings, err := GetSettingsFromCache(cacheKey); err != nil {
			return nil, err
		} else if cacheSettings != nil && len(*cacheSettings) > 0 {
			settings = *cacheSettings
		}
	}

	if len(settings) == 0 {
		if result := database.Pg.Model(&models.DomainSetting{}).
			Joins("JOIN domains ON domains.id = domain_settings.domain_id").
			Joins("JOIN apps ON apps.id = domains.app_id").
			Where("apps.name = ? AND domains.name = ? AND (level = 'both' OR level = ?)", appName, domainName, level.String()).
			Find(&settings); result.Error != nil {
			return nil, result.Error
		}
		_ = SetSettingsToCache(cacheKey, &settings)
	}

	return &settings, nil
}

// GetSettingsByDomainID method to get settings by domain ID.
func GetSettingsByDomainID(domainID uint, level enums.Level) (*[]models.DomainSetting, error) {
	var settings []models.DomainSetting
	cacheKey := SettingsCacheKeyOnId(domainID)

	if inCache, err := IsSettingsInCache(cacheKey); err != nil {
		return nil, err
	} else if inCache {
		if cacheSettings, err := GetSettingsFromCache(cacheKey); err != nil {
			return nil, err
		} else if cacheSettings != nil && len(*cacheSettings) > 0 {
			settings = *cacheSettings
		}
	}

	if len(settings) == 0 {
		if result := database.Pg.Model(&models.DomainSetting{}).
			Where("domain_id = ? AND (level = 'both' OR level = ?)", domainID, level.String()).
			Find(&settings); result.Error != nil {
			return nil, result.Error
		}
		_ = SetSettingsToCache(cacheKey, &settings)
	}

	return &settings, nil
}

// IsSettingsInCache checks if the settings exists in the cache.
func IsSettingsInCache(key string) (bool, error) {
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

// GetSettingsFromCache gets the settings from the cache.
func GetSettingsFromCache(key string) (*[]models.DomainSetting, error) {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Get().Key(key).Build())
	if result.Error() != nil {
		return nil, result.Error()
	}

	value, err := result.ToString()
	if err != nil {
		return nil, err
	}

	var settings []models.DomainSetting
	if err := json.Unmarshal([]byte(value), &settings); err != nil {
		return nil, err
	}

	return &settings, nil
}

// SetSettingsToCache sets the settings to the cache.
func SetSettingsToCache(key string, settings *[]models.DomainSetting) error {
	value, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	//log.Debug(value)
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

// DeleteSettingsFromCache deletes an existing setting from the cache.
func DeleteSettingsFromCache(key string) error {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Del().Key(key).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// SettingsCacheKeyOnName returns the key for the settings cache with a name.
func SettingsCacheKeyOnName(appName, domainName string) string {
	return fmt.Sprintf("%s:%s:settings", appName, domainName)
}

// SettingsCacheKeyOnId returns the key for the settings cache with an id.
func SettingsCacheKeyOnId(domainID uint) string {
	return fmt.Sprintf("settings:%d", domainID)
}
