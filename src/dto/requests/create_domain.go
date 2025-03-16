package requests

// CreateDomain struct for creating a new Domain.
type CreateDomain struct {
	AppID     uint            `json:"appId" validate:"required"`
	SSL       bool            `json:"ssl"`
	Name      string          `json:"name" validate:"required"`
	IpAddress string          `json:"ipAddress" validate:"required"`
	Settings  []DomainSetting `json:"settings" validate:"dive"`
}
