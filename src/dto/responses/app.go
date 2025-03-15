package responses

import (
	"api-app/main/src/models"
	"time"
)

// App struct to hold app data.
type App struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
	Domains   []Domain  `json:"domains"`
}

// SetApp method to set app data from models.App{}.
func (a *App) SetApp(app *models.App) {
	a.ID = app.ID
	a.Name = app.Name
	a.CreatedAt = app.CreatedAt
	a.UpdatedAt = app.UpdatedAt

	a.Domains = make([]Domain, len(app.Domains))
	for i := range app.Domains {
		a.Domains[i].SetDomain(&app.Domains[i])
	}
}
