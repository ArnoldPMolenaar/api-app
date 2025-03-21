package services

import (
	"api-app/main/src/database"
	"api-app/main/src/dto/requests"
	"api-app/main/src/enums"
	"api-app/main/src/models"
	"api-app/main/src/utils"
	"database/sql"
)

func IsDomainNameAvailable(appID uint, name string) (bool, error) {
	var count int64
	if result := database.Pg.Model(&models.Domain{}).
		Where("app_id = ? AND name = ?", appID, name).
		Count(&count); result.Error != nil {
		return false, result.Error
	}

	return count == 1, nil
}

// IsDomainDeleted method to check if a domain is deleted.
func IsDomainDeleted(id uint) (bool, error) {
	var count int64
	if result := database.Pg.Model(&models.Domain{}).
		Unscoped().
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Count(&count); result.Error != nil {
		return false, result.Error
	}

	return count == 1, nil
}

// GetDomainById method to get a domain by its ID.
func GetDomainById(id uint) (*models.Domain, error) {
	domain := &models.Domain{}

	if result := database.Pg.Preload("Settings").Find(domain, "id = ?", id); result.Error != nil {
		return nil, result.Error
	}

	return domain, nil
}

// CreateDomain method to create a domain.
func CreateDomain(appID uint, ssl bool, name, ipAddress string, settings *[]requests.DomainSetting) (*models.Domain, error) {
	subdomain, secondLevelDomain, topLevelDomain := utils.ExtractDomain(name)
	domain := models.Domain{
		AppID:       appID,
		SSL:         ssl,
		Name:        name,
		Sub:         sql.NullString{String: subdomain, Valid: subdomain != ""},
		SecondLevel: secondLevelDomain,
		TopLevel:    topLevelDomain,
		IpAddress:   ipAddress,
		Settings:    make([]models.DomainSetting, len(*settings)),
	}

	for i := range *settings {
		domain.Settings[i] = models.DomainSetting{
			Name:      (*settings)[i].Name,
			Level:     enums.Level((*settings)[i].Level),
			Value:     (*settings)[i].Value,
			ValueType: enums.ValueType((*settings)[i].ValueType),
		}
	}

	if result := database.Pg.Create(&domain); result.Error != nil {
		return nil, result.Error
	}

	return &domain, nil
}

// UpdateDomain method to update a domain.
func UpdateDomain(oldDomain *models.Domain, ssl bool, name, ipAddress string, settings *[]requests.DomainSetting) (*models.Domain, error) {
	subdomain, secondLevelDomain, topLevelDomain := utils.ExtractDomain(name)
	oldDomain.SSL = ssl
	oldDomain.Name = name
	oldDomain.Sub = sql.NullString{String: subdomain, Valid: subdomain != ""}
	oldDomain.SecondLevel = secondLevelDomain
	oldDomain.TopLevel = topLevelDomain
	oldDomain.IpAddress = ipAddress
	oldDomain.Settings = make([]models.DomainSetting, len(*settings))

	for i := range *settings {
		oldDomain.Settings[i] = models.DomainSetting{
			Name:      (*settings)[i].Name,
			Level:     enums.Level((*settings)[i].Level),
			Value:     (*settings)[i].Value,
			ValueType: enums.ValueType((*settings)[i].ValueType),
		}
	}

	if result := database.Pg.Save(oldDomain); result.Error != nil {
		return nil, result.Error
	}

	_ = deleteSettingsCache(oldDomain.ID, oldDomain.Name)

	return oldDomain, nil
}

// DeleteDomain method to delete a domain.
func DeleteDomain(domain *models.Domain) error {
	_ = deleteSettingsCache(domain.ID, domain.Name)

	return database.Pg.Delete(domain).Error
}

// RestoreDomain method to restore a domain.
func RestoreDomain(id uint) error {
	return database.Pg.Unscoped().Model(&models.Domain{}).Where("id = ?", id).Update("deleted_at", nil).Error
}

// deleteSettingsCache method to delete the settings cache.
func deleteSettingsCache(domainID uint, domainName string) error {
	if err := DeleteSettingsFromCache(SettingsCacheKeyOnId(domainID)); err != nil {
		return err
	}

	var appName string
	if result := database.Pg.Model(&models.Domain{}).
		Joins("JOIN apps ON apps.id = domains.app_id").
		Where("domains.id = ?", domainID).
		Select("apps.name").
		Scan(&appName); result.Error != nil {
		return result.Error
	}

	if err := DeleteSettingsFromCache(SettingsCacheKeyOnName(appName, domainName)); err != nil {
		return err
	}

	return nil
}
