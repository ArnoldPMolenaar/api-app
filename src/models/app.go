package models

import "gorm.io/gorm"

type App struct {
	gorm.Model
	Name string `gorm:"uniqueIndex:idx_name,sort:asc;not null"`

	// Relationships.
	Domains []Domain
}
