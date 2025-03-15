package requests

// CreateApp struct for creating a new App.
type CreateApp struct {
	Name    string         `json:"name" validate:"required"`
	Domains []CreateDomain `json:"domains" validate:"required,dive"`
}
