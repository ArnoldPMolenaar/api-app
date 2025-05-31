package requests

import "time"

// UpdateApp struct for updating a existing App.
type UpdateApp struct {
	Name      string            `json:"name" validate:"required"`
	Settings  []AppSetting      `json:"settings" validate:"dive"`
	Domains   []UpdateAppDomain `json:"domains" validate:"required,dive"`
	UpdatedAt time.Time         `json:"updatedAt" validate:"required"`
}
