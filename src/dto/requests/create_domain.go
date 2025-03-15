package requests

// CreateDomain struct for creating a new Domain.
type CreateDomain struct {
	SSL       bool   `json:"ssl"`
	Name      string `json:"name" validate:"required"`
	IpAddress string `json:"ipAddress" validate:"required"`
}
