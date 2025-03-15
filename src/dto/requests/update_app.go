package requests

import "time"

// UpdateApp struct for updating a existing App.
type UpdateApp struct {
	Name      string         `json:"name" validate:"required"`
	Domains   []UpdateDomain `json:"domains" validate:"required,dive"`
	UpdatedAt time.Time      `json:"updatedAt" validate:"required"`
}
