package models

import (
	"database/sql"
	"gorm.io/gorm"
)

type Domain struct {
	gorm.Model
	AppID       uint   `gorm:"uniqueIndex:idx_app_name;not null"`
	SSL         bool   `gorm:"default:false;not null"`
	Name        string `gorm:"uniqueIndex:idx_app_name;not null"`
	Sub         sql.NullString
	SecondLevel string `gorm:"not null"`
	TopLevel    string `gorm:"not null"`
	IpAddress   string `gorm:"not null"`

	// Relationships.
	App      App `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppID;references:ID"`
	Settings []DomainSetting
}
