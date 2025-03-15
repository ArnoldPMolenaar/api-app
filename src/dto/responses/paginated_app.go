package responses

import (
	"api-app/main/src/models"
	"time"
)

// PaginatedApp struct to hold paginated app data.
type PaginatedApp struct {
	ID        uint      `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt"`
}

// SetPaginatedApp method to set app data from models.App{}.
func (a *PaginatedApp) SetPaginatedApp(app *models.App) {
	a.ID = app.ID
	a.Name = app.Name
	a.CreatedAt = app.CreatedAt
	a.UpdatedAt = app.UpdatedAt
}
