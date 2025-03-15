package routes

import (
	"api-app/main/src/controllers"
	"github.com/ArnoldPMolenaar/api-utils/middleware"
	"github.com/gofiber/fiber/v2"
)

// PrivateRoutes func for describe group of private routes.
func PrivateRoutes(a *fiber.App) {
	// Create private routes group.
	route := a.Group("/v1")

	// Register CRUD routes for /v1/apps.
	apps := route.Group("/apps", middleware.MachineProtected())
	apps.Get("/", controllers.GetApps)
	apps.Post("/", controllers.CreateApp)
	apps.Get("/:id", controllers.GetApp)
	apps.Put("/:id", controllers.UpdateApp)
	apps.Delete("/:id", controllers.DeleteApp)
	apps.Put("/:id/restore", controllers.RestoreApp)
}
