package main

import (
	"api-app/main/src/cache"
	"api-app/main/src/configs"
	"api-app/main/src/database"
	"api-app/main/src/middleware"
	"api-app/main/src/routes"
	"fmt"
	"os"

	routeutil "github.com/ArnoldPMolenaar/api-utils/routes"
	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/log"
)

var (
	version string
)

func main() {
	log.Info("Running version: ", version)

	// Define Fiber config.
	config := configs.FiberConfig()

	// Define a new Fiber app with config.
	app := fiber.New(config)

	// Register Fiber's middleware for app.
	middleware.FiberMiddleware(app)

	// Open database connection.
	if err := database.OpenDBConnection(); err != nil {
		panic(fmt.Sprintf("Could not connect to the database: %v", err))
	}

	// Open Valkey connection.
	if err := cache.OpenValkeyConnection(); err != nil {
		panic(fmt.Sprintf("Could not connect to the cache: %v", err))
	}
	defer cache.Valkey.Close()

	// Register a private routes_util for app.
	routes.PrivateRoutes(app)
	// Register a public routes_util for app.
	routes.PublicRoutes(app)
	// Register k8s routes.
	routeutil.HealthRoute(app)
	routeutil.KubernetesProbeRoutes(app, database.ReadinessCheck, cache.ReadinessCheck, database.MigrationReadinessCheck)
	// Register route for 404 Error.
	routeutil.NotFoundRoute(app)

	// Start server (with or without graceful shutdown).
	if os.Getenv("STAGE_STATUS") == "dev" {
		utils.StartServer(app)
	} else {
		utils.StartServerWithGracefulShutdown(app)
	}
}
