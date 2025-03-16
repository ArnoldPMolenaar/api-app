package responses

import "api-app/main/src/models"

// DomainSetting struct to handle domain setting response.
type DomainSetting struct {
	DomainID  uint   `json:"domainId"`
	Name      string `json:"name"`
	Level     string `json:"level"`
	Value     string `json:"value"`
	ValueType string `json:"valueType"`
}

// SetDomainSetting method to set domain setting data from models.DomainSetting{}.
func (ds *DomainSetting) SetDomainSetting(domainSetting *models.DomainSetting) {
	ds.DomainID = domainSetting.DomainID
	ds.Name = domainSetting.Name
	ds.Level = domainSetting.Level.String()
	ds.Value = domainSetting.Value
	ds.ValueType = domainSetting.ValueType.String()
}
