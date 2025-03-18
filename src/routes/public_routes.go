package routes

import (
	"api-app/main/src/controllers"
	"api-app/main/src/enums"
	"github.com/gofiber/fiber/v2"
)

// PublicRoutes func for describe group of public routes.
func PublicRoutes(a *fiber.App) {
	// Create private routes group.
	route := a.Group("/v1")

	// Register routes for /v1/settings.
	settings := route.Group("/settings")
	settings.Get("/", func(c *fiber.Ctx) error {
		return controllers.GetSettingsByDomainName(c, enums.Public)
	})
	settings.Get("/:id", func(c *fiber.Ctx) error {
		return controllers.GetSettingsByDomainID(c, enums.Public)
	})
}
