package responses

import (
	"api-app/main/src/models"
	"time"
)

// Domain struct to handle domain response.
type Domain struct {
	ID          uint      `json:"id"`
	AppID       uint      `json:"appId"`
	SSL         bool      `json:"ssl"`
	Name        string    `json:"name"`
	Sub         *string   `json:"sub"`
	SecondLevel string    `json:"secondLevel"`
	TopLevel    string    `json:"topLevel"`
	IpAddress   string    `json:"ipAddress"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

// SetDomain method to set domain data from models.Domain{}.
func (d *Domain) SetDomain(domain *models.Domain) {
	d.ID = domain.ID
	d.AppID = domain.AppID
	d.SSL = domain.SSL
	d.Name = domain.Name
	if domain.Sub.Valid {
		d.Sub = &domain.Sub.String
	}
	d.SecondLevel = domain.SecondLevel
	d.TopLevel = domain.TopLevel
	d.IpAddress = domain.IpAddress
	d.CreatedAt = domain.CreatedAt
	d.UpdatedAt = domain.UpdatedAt
}
