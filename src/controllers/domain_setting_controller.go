package controllers

import (
	"api-app/main/src/enums"
	"api-app/main/src/errors"
	"api-app/main/src/models"
	"api-app/main/src/services"
	"encoding/json"
	"fmt"
	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"time"
)

// GetSettingsByDomainName function to get settings by domain name.
func GetSettingsByDomainName(c *fiber.Ctx, level enums.Level) error {
	// Get the AppName and DomainName parameter from the URL.
	appName := c.Query("app")
	if appName == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.MissingRequiredParam, "App Name is required.")
	}
	domainName := c.Query("domain")
	if domainName == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.MissingRequiredParam, "Domain Name is required.")
	}

	// Get the settings.
	settings, err := services.GetSettingsByDomainName(appName, domainName, level)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the settings.
	response, err := toSettingsResponse(settings)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errors.DomainSettings, err.Error())
	}

	return c.JSON(response)
}

// GetSettingsByDomainID function to get settings by domain ID.
func GetSettingsByDomainID(c *fiber.Ctx, level enums.Level) error {
	// Get the domainID parameter from the URL.
	domainIDParam := c.Params("id")
	if domainIDParam == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.MissingRequiredParam, "Domain ID is required.")
	}
	domainID, err := utils.StringToUint(domainIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "Invalid Domain ID.")
	}

	// Get the settings.
	settings, err := services.GetSettingsByDomainID(domainID, level)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the settings.
	response, err := toSettingsResponse(settings)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errors.DomainSettings, err.Error())
	}

	return c.JSON(response)
}

// toSettingsResponse converts an array of DomainSetting structs to a dynamic JSON object.
func toSettingsResponse(settings *[]models.DomainSetting) (map[string]interface{}, error) {
	response := make(map[string]interface{})

	for i := range *settings {
		var setting = (*settings)[i]
		var value interface{}
		var err error

		switch setting.ValueType {
		case enums.Int:
			value, err = strconv.Atoi(setting.Value)
		case enums.Float:
			value, err = strconv.ParseFloat(setting.Value, 64)
		case enums.String:
			value = setting.Value
		case enums.Bool:
			value, err = strconv.ParseBool(setting.Value)
		case enums.Date:
			value, err = time.Parse(time.DateOnly, setting.Value)
		case enums.DateTime:
			value, err = time.Parse(time.DateTime, setting.Value)
		case enums.JSON:
			var js json.RawMessage
			err = json.Unmarshal([]byte(setting.Value), &js)
			value = js
		default:
			err = fmt.Errorf("unknown ValueType for setting %s", setting.Name)
		}

		if err != nil {
			return nil, fmt.Errorf("error converting setting %s: %v", setting.Name, err)
		}

		response[setting.Name] = value
	}

	return response, nil
}
