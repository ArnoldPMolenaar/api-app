package database

import (
	"api-app/main/src/models"
	"context"
	"errors"
	"fmt"
	"sort"
	"time"

	utilsdb "github.com/ArnoldPMolenaar/api-utils/database"
	"gorm.io/gorm"
)

var Pg *gorm.DB

// OpenDBConnection Start a new database connection.
// Also tries to migrate the database schema.
func OpenDBConnection() error {
	// Open connection to database.
	db, err := utilsdb.PostgresSQLConnection()
	if err != nil {
		return err
	}

	// Migrate the database schema.
	err = Migrate(db)
	if err != nil {
		return err
	}

	// Set the global DB variable.
	Pg = db

	return nil
}

// ReadinessCheck verifies that the database connection is initialized and reachable.
func ReadinessCheck() error {
	if Pg == nil {
		return errors.New("database connection is not initialized")
	}

	sqlDB, err := Pg.DB()
	if err != nil {
		return fmt.Errorf("database sql handle unavailable: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return fmt.Errorf("database ping failed: %w", err)
	}

	return nil
}

// MigrationReadinessCheck verifies that required tables and seed data exist.
func MigrationReadinessCheck() error {
	if Pg == nil {
		return errors.New("database connection is not initialized")
	}

	requiredTables := []any{
		&models.App{},
		&models.AppSetting{},
		&models.Domain{},
		&models.DomainSetting{},
	}
	for _, table := range requiredTables {
		if !Pg.Migrator().HasTable(table) {
			return fmt.Errorf("missing required table for %T", table)
		}
	}

	type enumCheck struct {
		typeName string
		labels   []string
	}

	enumChecks := []enumCheck{
		{typeName: "level", labels: []string{"public", "private", "both"}},
		{typeName: "value_type", labels: []string{"int", "float", "string", "bool", "date", "datetime", "json"}},
	}

	for _, check := range enumChecks {
		labels, err := getEnumLabels(check.typeName)
		if err != nil {
			return err
		}

		sort.Strings(labels)
		sort.Strings(check.labels)
		if len(labels) != len(check.labels) {
			return fmt.Errorf("%s enum labels mismatch: have %v, want %v", check.typeName, labels, check.labels)
		}

		for i := range labels {
			if labels[i] != check.labels[i] {
				return fmt.Errorf("%s enum labels mismatch: have %v, want %v", check.typeName, labels, check.labels)
			}
		}
	}

	return nil
}

func getEnumLabels(typeName string) ([]string, error) {
	type row struct {
		Label string
	}

	rows := make([]row, 0)
	err := Pg.Raw(`
		SELECT e.enumlabel AS label
		FROM pg_type t
		JOIN pg_enum e ON t.oid = e.enumtypid
		WHERE t.typname = ?
		ORDER BY e.enumsortorder
	`, typeName).Scan(&rows).Error
	if err != nil {
		return nil, fmt.Errorf("failed to validate %s enum: %w", typeName, err)
	}

	if len(rows) == 0 {
		return nil, fmt.Errorf("missing required enum type: %s", typeName)
	}

	labels := make([]string, len(rows))
	for i := range rows {
		labels[i] = rows[i].Label
	}

	return labels, nil
}
