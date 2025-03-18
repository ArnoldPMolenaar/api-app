package routes

import (
	"api-app/main/src/controllers"
	"api-app/main/src/enums"
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

	// Register CRUD routes for /v1/domains.
	domains := route.Group("/domains", middleware.MachineProtected())
	domains.Post("/", controllers.CreateDomain)
	domains.Get("/:id", controllers.GetDomain)
	domains.Put("/:id", controllers.UpdateDomain)
	domains.Delete("/:id", controllers.DeleteDomain)
	domains.Put("/:id/restore", controllers.RestoreDomain)

	// Register routes for /v1/domains/settings.
	domainSettings := domains.Group("/settings")
	domainSettings.Get("/", func(c *fiber.Ctx) error {
		return controllers.GetSettingsByDomainName(c, enums.Private)
	})
	domainSettings.Get("/:id", func(c *fiber.Ctx) error {
		return controllers.GetSettingsByDomainID(c, enums.Private)
	})
}
