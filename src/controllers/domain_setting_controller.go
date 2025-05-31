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

	// Get the app settings.
	appSettings, err := services.GetAppSettingsByName(appName, level)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Get the domain settings.
	domainSettings, err := services.GetDomainSettingsByName(appName, domainName, level)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the settings.
	response, err := toSettingsResponse(appSettings, domainSettings)
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

	// Get the appID with the domainID.
	appID, err := services.GetAppIDByDomainID(domainID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Get the app settings.
	appSettings, err := services.GetAppSettingsByAppID(appID, level)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Get the domain settings.
	domainSettings, err := services.GetDomainSettingsByDomainID(domainID, level)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the settings.
	response, err := toSettingsResponse(appSettings, domainSettings)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errors.DomainSettings, err.Error())
	}

	return c.JSON(response)
}

// toSettingsResponse converts an array of DomainSetting structs to a dynamic JSON object.
func toSettingsResponse(appSettings *[]models.AppSetting, domainSettings *[]models.DomainSetting) (map[string]interface{}, error) {
	response := make(map[string]interface{})

	convertSetting := func(name, valueType, value string) (interface{}, error) {
		switch enums.ValueType(valueType) {
		case enums.Int:
			return strconv.Atoi(value)
		case enums.Float:
			return strconv.ParseFloat(value, 64)
		case enums.String:
			return value, nil
		case enums.Bool:
			return strconv.ParseBool(value)
		case enums.Date:
			return time.Parse(time.DateOnly, value)
		case enums.DateTime:
			return time.Parse(time.DateTime, value)
		case enums.JSON:
			var js json.RawMessage
			err := json.Unmarshal([]byte(value), &js)
			return js, err
		default:
			return nil, fmt.Errorf("unknown ValueType for setting %s", name)
		}
	}

	for i := range *appSettings {
		setting := (*appSettings)[i]
		value, err := convertSetting(setting.Name, string(setting.ValueType), setting.Value)
		if err != nil {
			return nil, fmt.Errorf("error converting setting %s: %v", setting.Name, err)
		}
		response[setting.Name] = value
	}

	if domainSettings != nil {
		for i := range *domainSettings {
			setting := (*domainSettings)[i]
			value, err := convertSetting(setting.Name, string(setting.ValueType), setting.Value)
			if err != nil {
				return nil, fmt.Errorf("error converting setting %s: %v", setting.Name, err)
			}
			response[setting.Name] = value
		}
	}

	return response, nil
}
