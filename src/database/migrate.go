package database

import (
	"api-app/main/src/models"
	"gorm.io/gorm"
)

// Migrate the database schema.
// See: https://gorm.io/docs/migration.html#Auto-Migration
func Migrate(db *gorm.DB) error {
	// Adds the level enum type to the database.
	if tx := db.Exec(`DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'level') THEN 
			CREATE TYPE level AS ENUM ('public', 'private', 'both'); 
		END IF; 
	END $$;`); tx.Error != nil {
		return tx.Error
	}

	// Adds the value_type enum type to the database.
	if tx := db.Exec(`DO $$ 
	BEGIN 
		IF NOT EXISTS (SELECT 1 FROM pg_type WHERE typname = 'value_type') THEN 
			CREATE TYPE value_type AS ENUM ('int', 'float', 'string', 'bool', 'date', 'datetime', 'json'); 
		END IF; 
	END $$;`); tx.Error != nil {
		return tx.Error
	}

	err := db.AutoMigrate(&models.App{}, &models.Domain{}, &models.DomainSetting{})
	if err != nil {
		return err
	}

	// Seed App.
	apps := []string{"Admin"}
	for _, app := range apps {
		if err := db.FirstOrCreate(&models.App{}, models.App{Name: app}).Error; err != nil {
			return err
		}
	}

	return nil
}
