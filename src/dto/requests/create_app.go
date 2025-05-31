package requests

// CreateApp struct for creating a new App.
type CreateApp struct {
	Name     string            `json:"name" validate:"required"`
	Settings []AppSetting      `json:"settings" validate:"dive"`
	Domains  []CreateAppDomain `json:"domains" validate:"required,dive"`
}
