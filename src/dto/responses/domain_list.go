package responses

import "api-app/main/src/models"

type DomainList struct {
	Domains []Domain `json:"domains"`
}

// SetDomains method to set domains data from models.Domain{}.
func (d *DomainList) SetDomains(domains *[]models.Domain) {
	d.Domains = make([]Domain, len(*domains))

	for i, domain := range *domains {
		d.Domains[i] = Domain{}
		d.Domains[i].SetDomain(&domain)
	}
}
