package controllers

import (
	"api-app/main/src/enums"
	"api-app/main/src/errors"
	"api-app/main/src/services"
	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v2"
)

// GetSettingsByAppName function to get settings by app name.
func GetSettingsByAppName(c *fiber.Ctx, level enums.Level) error {
	// Get the AppName and DomainName parameter from the URL.
	appName := c.Query("app")
	if appName == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.MissingRequiredParam, "App Name is required.")
	}

	// Get the app settings.
	appSettings, err := services.GetAppSettingsByName(appName, level)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the settings.
	response, err := toSettingsResponse(appSettings, nil)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errors.DomainSettings, err.Error())
	}

	return c.JSON(response)
}

// GetSettingsByAppID function to get settings by app ID.
func GetSettingsByAppID(c *fiber.Ctx, level enums.Level) error {
	// Get the domainID parameter from the URL.
	appIDParam := c.Params("id")
	if appIDParam == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.MissingRequiredParam, "App ID is required.")
	}
	appID, err := utils.StringToUint(appIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "Invalid App ID.")
	}

	// Get the app settings.
	appSettings, err := services.GetAppSettingsByAppID(appID, level)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the settings.
	response, err := toSettingsResponse(appSettings, nil)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errors.DomainSettings, err.Error())
	}

	return c.JSON(response)
}
