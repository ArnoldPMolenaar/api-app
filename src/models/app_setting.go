package models

import "api-app/main/src/enums"

type AppSetting struct {
	AppID     uint            `gorm:"primaryKey;autoIncrement:false"`
	Name      string          `gorm:"primaryKey;autoIncrement:false"`
	Level     enums.Level     `gorm:"primaryKey;autoIncrement:false;type:level"`
	Value     string          `gorm:"not null"`
	ValueType enums.ValueType `gorm:"not null;type:value_type"`

	// Relationships.
	App App `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;foreignKey:AppID;references:ID"`
}
