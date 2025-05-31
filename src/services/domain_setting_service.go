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

// GetDomainSettingsByName method to get settings by domain name.
func GetDomainSettingsByName(appName, domainName string, level enums.Level) (*[]models.DomainSetting, error) {
	var settings []models.DomainSetting
	cacheKey := DomainSettingsCacheKeyOnName(appName, domainName, level)

	if inCache, err := IsDomainSettingsInCache(cacheKey); err != nil {
		return nil, err
	} else if inCache {
		if cacheSettings, err := GetDomainSettingsFromCache(cacheKey); err != nil {
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
		_ = SetDomainSettingsToCache(cacheKey, &settings)
	}

	return &settings, nil
}

// GetDomainSettingsByDomainID method to get settings by domain ID.
func GetDomainSettingsByDomainID(domainID uint, level enums.Level) (*[]models.DomainSetting, error) {
	var settings []models.DomainSetting
	cacheKey := DomainSettingsCacheKeyOnId(domainID, level)

	if inCache, err := IsDomainSettingsInCache(cacheKey); err != nil {
		return nil, err
	} else if inCache {
		if cacheSettings, err := GetDomainSettingsFromCache(cacheKey); err != nil {
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
		_ = SetDomainSettingsToCache(cacheKey, &settings)
	}

	return &settings, nil
}

// IsDomainSettingsInCache checks if the settings exists in the cache.
func IsDomainSettingsInCache(key string) (bool, error) {
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

// GetDomainSettingsFromCache gets the settings from the cache.
func GetDomainSettingsFromCache(key string) (*[]models.DomainSetting, error) {
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

// SetDomainSettingsToCache sets the settings to the cache.
func SetDomainSettingsToCache(key string, settings *[]models.DomainSetting) error {
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

// DeleteDomainSettingsFromCache deletes an existing setting from the cache.
func DeleteDomainSettingsFromCache(key string) error {
	result := cache.Valkey.Do(context.Background(), cache.Valkey.B().Del().Key(key).Build())
	if result.Error() != nil {
		return result.Error()
	}

	return nil
}

// DomainSettingsCacheKeyOnName returns the key for the settings cache with a name.
func DomainSettingsCacheKeyOnName(appName, domainName string, level enums.Level) string {
	return fmt.Sprintf("%s:%s:settings:%s", appName, domainName, level.String())
}

// DomainSettingsCacheKeyOnId returns the key for the settings cache with an id.
func DomainSettingsCacheKeyOnId(domainID uint, level enums.Level) string {
	return fmt.Sprintf("settings:domains:%d:%s", domainID, level.String())
}
