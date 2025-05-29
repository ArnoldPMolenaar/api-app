package services

import (
	"api-app/main/src/database"
	"api-app/main/src/dto/requests"
	"api-app/main/src/models"
	"api-app/main/src/utils"
	"database/sql"
)

// IsAppAvailable method to check if an app is available.
func IsAppAvailable(app string) (bool, error) {
	if result := database.Pg.Limit(1).Find(&models.App{}, "name = ?", app); result.Error != nil {
		return false, result.Error
	} else {
		return result.RowsAffected == 1, nil
	}
}

// AreAppsAvailable checks if all the given app names exist.
func AreAppsAvailable(apps []string) (bool, error) {
	var foundApps int64
	result := database.Pg.Model(&models.App{}).Where("name IN ?", apps).Count(&foundApps)
	if result.Error != nil {
		return false, result.Error
	}
	return int64(len(apps)) == foundApps, nil
}

// IsAppDeleted method to check if an app is deleted.
func IsAppDeleted(id uint) (bool, error) {
	var count int64
	if result := database.Pg.Model(&models.App{}).
		Unscoped().
		Where("id = ? AND deleted_at IS NOT NULL", id).
		Count(&count); result.Error != nil {
		return false, result.Error
	}

	return count == 1, nil
}

// GetAppById method to get an app by its ID.
func GetAppById(id uint, unscoped ...bool) (*models.App, error) {
	app := &models.App{}
	query := database.Pg

	if len(unscoped) > 0 && unscoped[0] {
		query = query.Unscoped()
	}

	if result := query.Preload("Domains").Find(app, "id = ?", id); result.Error != nil {
		return nil, result.Error
	}

	return app, nil
}

// CreateApp method to create an app.
func CreateApp(name string, domains *[]requests.CreateAppDomain) (*models.App, error) {
	app := models.App{
		Name:    name,
		Domains: make([]models.Domain, len(*domains)),
	}

	for i := range *domains {
		subdomain, secondLevelDomain, topLevelDomain := utils.ExtractDomain((*domains)[i].Name)

		app.Domains[i] = models.Domain{
			SSL:         (*domains)[i].SSL,
			Name:        (*domains)[i].Name,
			Sub:         sql.NullString{String: subdomain, Valid: subdomain != ""},
			SecondLevel: secondLevelDomain,
			TopLevel:    topLevelDomain,
			IpAddress:   (*domains)[i].IpAddress,
		}
	}

	if result := database.Pg.Create(&app); result.Error != nil {
		return nil, result.Error
	}

	return &app, nil
}

// UpdateApp method to update an app.
func UpdateApp(oldApp *models.App, name string, domains *[]requests.UpdateAppDomain) (*models.App, error) {
	// Start a new transaction
	tx := database.Pg.Begin()
	if tx.Error != nil {
		return nil, tx.Error
	}

	oldApp.Name = name

	// Create a map for quick lookup of new domains by name.
	newDomainsMap := make(map[uint]requests.UpdateAppDomain)
	for i := range *domains {
		if (*domains)[i].ID != 0 {
			newDomainsMap[(*domains)[i].ID] = (*domains)[i]
		}
	}

	// Iterate through old domains to update or mark as deleted.
	for i := range oldApp.Domains {
		oldDomain := &oldApp.Domains[i]
		if newDomain, exists := newDomainsMap[oldDomain.ID]; exists {
			// Update existing domain.
			oldDomain.SSL = newDomain.SSL
			subdomain, secondLevelDomain, topLevelDomain := utils.ExtractDomain(newDomain.Name)
			oldDomain.Sub = sql.NullString{String: subdomain, Valid: subdomain != ""}
			oldDomain.SecondLevel = secondLevelDomain
			oldDomain.TopLevel = topLevelDomain
			oldDomain.IpAddress = newDomain.IpAddress

			// Restore if it was previously deleted.
			if oldDomain.DeletedAt.Valid {
				oldDomain.DeletedAt.Valid = false
			}

			if result := tx.Save(&oldDomain); result.Error != nil {
				tx.Rollback()
				return nil, result.Error
			}

			// Remove from newDomainsMap as it is already processed.
			delete(newDomainsMap, oldDomain.ID)
		} else {
			// Mark as deleted if not in new domains
			if result := tx.Delete(&oldDomain); result.Error != nil {
				tx.Rollback()
				return nil, result.Error
			}
		}
	}

	// Add new domains that were not in old domains.
	for _, newDomain := range newDomainsMap {
		subdomain, secondLevelDomain, topLevelDomain := utils.ExtractDomain(newDomain.Name)
		oldApp.Domains = append(oldApp.Domains, models.Domain{
			SSL:         newDomain.SSL,
			Name:        newDomain.Name,
			Sub:         sql.NullString{String: subdomain, Valid: subdomain != ""},
			SecondLevel: secondLevelDomain,
			TopLevel:    topLevelDomain,
			IpAddress:   newDomain.IpAddress,
		})
	}

	if result := tx.Save(oldApp); result.Error != nil {
		tx.Rollback()
		return nil, result.Error
	}

	// Commit the transaction
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return nil, err
	}

	// Retrieve the updated app. Because new domains are added and now have IDs.
	newApp, err := GetAppById(oldApp.ID)
	if err != nil {
		return nil, err
	}

	return newApp, nil
}

// DeleteApp method to delete an app.
func DeleteApp(app *models.App) error {
	return database.Pg.Delete(app).Error
}

// RestoreApp method to restore a deleted app.
func RestoreApp(id uint) error {
	return database.Pg.Unscoped().Model(&models.App{}).Where("id = ?", id).Update("deleted_at", nil).Error
}
