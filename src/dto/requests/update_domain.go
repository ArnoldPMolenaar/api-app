package requests

import "time"

// UpdateDomain struct for updating a existing Domain.
type UpdateDomain struct {
	ID        uint      `json:"id" validate:"required"`
	SSL       bool      `json:"ssl"`
	Name      string    `json:"name" validate:"required"`
	IpAddress string    `json:"ipAddress" validate:"required"`
	UpdatedAt time.Time `json:"updatedAt" validate:"required"`
}
