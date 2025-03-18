package services

import (
	"api-app/main/src/database"
	"api-app/main/src/enums"
	"api-app/main/src/models"
)

// GetSettingsByDomainName method to get settings by domain name.
func GetSettingsByDomainName(appName, domainName string, level enums.Level) (*[]models.DomainSetting, error) {
	var settings []models.DomainSetting

	if result := database.Pg.Model(&models.DomainSetting{}).
		Joins("JOIN domains ON domains.id = domain_settings.domain_id").
		Joins("JOIN apps ON apps.id = domains.app_id").
		Where("apps.name = ? AND domains.name = ? AND (level = 'both' OR level = ?)", appName, domainName, level.String()).
		Find(&settings); result.Error != nil {
		return nil, result.Error
	}

	return &settings, nil
}

// GetSettingsByDomainID method to get settings by domain ID.
func GetSettingsByDomainID(domainID uint, level enums.Level) (*[]models.DomainSetting, error) {
	var settings []models.DomainSetting

	if result := database.Pg.Model(&models.DomainSetting{}).
		Where("domain_id = ? AND (level = 'both' OR level = ?)", domainID, level.String()).
		Find(&settings); result.Error != nil {
		return nil, result.Error
	}

	return &settings, nil
}
