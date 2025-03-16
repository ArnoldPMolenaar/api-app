package requests

// CreateAppDomain struct for creating a new Domain.
type CreateAppDomain struct {
	SSL       bool   `json:"ssl"`
	Name      string `json:"name" validate:"required"`
	IpAddress string `json:"ipAddress" validate:"required"`
}
