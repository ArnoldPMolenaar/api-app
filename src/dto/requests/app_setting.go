package requests

// AppSetting struct for creating or updating a AppSetting.
type AppSetting struct {
	Name      string `json:"name" validate:"required"`
	Level     string `json:"level" validate:"required"`
	Value     string `json:"value" validate:"required"`
	ValueType string `json:"valueType" validate:"required"`
}
