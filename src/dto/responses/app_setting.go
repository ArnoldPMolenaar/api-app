package responses

import "api-app/main/src/models"

// AppSetting struct to handle app setting response.
type AppSetting struct {
	Name      string `json:"name"`
	Level     string `json:"level"`
	Value     string `json:"value"`
	ValueType string `json:"valueType"`
}

// SetAppSetting method to set app setting data from models.AppSetting{}.
func (as *AppSetting) SetAppSetting(appSetting *models.AppSetting) {
	as.Name = appSetting.Name
	as.Level = appSetting.Level.String()
	as.Value = appSetting.Value
	as.ValueType = appSetting.ValueType.String()
}
