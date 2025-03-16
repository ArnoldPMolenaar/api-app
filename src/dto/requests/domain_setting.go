package requests

// DomainSetting struct for creating a new DomainSetting.
type DomainSetting struct {
	DomainID  uint   `json:"domainId" validate:"required"`
	Name      string `json:"name" validate:"required"`
	Level     string `json:"level" validate:"required"`
	Value     string `json:"value" validate:"required"`
	ValueType string `json:"valueType" validate:"required"`
}
