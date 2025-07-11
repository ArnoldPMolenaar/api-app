package controllers

import (
	"api-app/main/src/database"
	"api-app/main/src/dto/requests"
	"api-app/main/src/dto/responses"
	"api-app/main/src/enums"
	"api-app/main/src/errors"
	"api-app/main/src/models"
	"api-app/main/src/services"
	"encoding/json"
	"fmt"
	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/ArnoldPMolenaar/api-utils/pagination"
	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
	"time"
)

// AreAppsAvailable checks if all given apps exist.
func AreAppsAvailable(c *fiber.Ctx) error {
	// Get the appName parameter from the query string.
	query := &requests.AreAppsAvailable{}
	if err := c.QueryParser(query); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Check if the apps exists.
	exist, err := services.AreAppsAvailable(query.AppNames)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.JSON(responses.Exists{Exists: exist})
}

// GetApps function fetches all apps from the database.
func GetApps(c *fiber.Ctx) error {
	apps := make([]models.App, 0)
	values := c.Request().URI().QueryArgs()
	allowedColumns := map[string]bool{
		"id":         true,
		"name":       true,
		"created_at": true,
		"updated_at": true,
	}

	queryFunc := pagination.Query(values, allowedColumns)
	sortFunc := pagination.Sort(values, allowedColumns)
	page := c.QueryInt("page", 1)
	if page < 1 {
		page = 1
	}
	limit := c.QueryInt("limit", 10)
	if limit < 1 {
		limit = 10
	}
	offset := pagination.Offset(page, limit)

	db := database.Pg.Unscoped().Scopes(queryFunc, sortFunc).Limit(limit).Offset(offset).Find(&apps)
	if db.Error != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, db.Error.Error())
	}

	total := int64(0)
	database.Pg.Unscoped().Scopes(queryFunc).Model(&models.App{}).Count(&total)
	pageCount := pagination.Count(int(total), limit)

	paginatedApps := make([]responses.PaginatedApp, len(apps))
	for i := range apps {
		paginatedApp := responses.PaginatedApp{}
		paginatedApp.SetPaginatedApp(&apps[i])
		paginatedApps[i] = paginatedApp
	}

	paginationModel := pagination.CreatePaginationModel(limit, page, pageCount, int(total), paginatedApps)

	return c.Status(fiber.StatusOK).JSON(paginationModel)
}

// GetApp function fetches an app from the database by its ID.
func GetApp(c *fiber.Ctx) error {
	// Get the appID parameter from the URL.
	appIDParam := c.Params("id")
	if appIDParam == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.MissingRequiredParam, "App ID is required.")
	}
	appID, err := utils.StringToUint(appIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "Invalid App ID.")
	}

	// Get the app.
	app, err := services.GetAppById(appID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if app.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errorutil.NotFound, "App not found.")
	}

	// Return the app.
	response := responses.App{}
	response.SetApp(app)

	return c.JSON(response)
}

// CreateApp func to create a app.
func CreateApp(c *fiber.Ctx) error {
	// Parse the request.
	request := requests.CreateApp{}
	if err := c.BodyParser(&request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate app fields.
	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, utils.ValidatorErrors(err))
	}
	if validationErrors := validateAppSettings(&request.Settings); validationErrors != "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.AppSettings, validationErrors)
	}

	// Check if app exists.
	if available, err := services.IsAppAvailable(request.Name); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.AppAvailable, "AppName already available.")
	}

	// Create the app.
	app, err := services.CreateApp(&request)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the app.
	response := responses.App{}
	response.SetApp(app)

	return c.JSON(response)
}

// UpdateApp func to update a app.
func UpdateApp(c *fiber.Ctx) error {
	// Get the ID from the URL.
	appIDParam := c.Params("id")
	if appIDParam == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.MissingRequiredParam, "App ID is required.")
	}
	appID, err := utils.StringToUint(appIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Parse the request.
	request := requests.UpdateApp{}
	if err := c.BodyParser(&request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate app fields.
	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, utils.ValidatorErrors(err))
	}
	if validationErrors := validateAppSettings(&request.Settings); validationErrors != "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.AppSettings, validationErrors)
	}

	// Check if app exists.
	app, err := services.GetAppById(appID, true)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if app.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errorutil.NotFound, "App not found.")
	}

	// Check if app name is unique.
	if app.Name != request.Name {
		if available, err := services.IsAppAvailable(request.Name); err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
		} else if available {
			return errorutil.Response(c, fiber.StatusBadRequest, errors.AppAvailable, "AppName already available.")
		}
	}

	// Check if the app data has been modified since it was last fetched.
	if request.UpdatedAt.Unix() < app.UpdatedAt.Unix() {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.OutOfSync, "Data is out of sync.")
	}
	for i := range app.Domains {
		if app.Domains[i].DeletedAt.Valid {
			continue
		}

		var newDomain *requests.UpdateAppDomain = nil
		for j := range request.Domains {
			if app.Domains[i].ID == request.Domains[j].ID {
				newDomain = &request.Domains[j]
				break
			}
		}
		if newDomain != nil && newDomain.UpdatedAt.Unix() < app.Domains[i].UpdatedAt.Unix() {
			return errorutil.Response(c, fiber.StatusBadRequest, errorutil.OutOfSync, "Data is out of sync.")
		}
	}

	// Update the app.
	updatedApp, err := services.UpdateApp(app, &request)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the app.
	response := responses.App{}
	response.SetApp(updatedApp)

	return c.JSON(response)
}

// DeleteApp func to delete a app.
func DeleteApp(c *fiber.Ctx) error {
	// Get the ID from the URL.
	appIDParam := c.Params("id")
	if appIDParam == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.MissingRequiredParam, "App ID is required.")
	}
	appID, err := utils.StringToUint(appIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the app.
	app, err := services.GetAppById(appID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if app.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.AppExists, "App does not exist.")
	}

	// Delete the app.
	if err := services.DeleteApp(app); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// RestoreApp func to restore a app.
func RestoreApp(c *fiber.Ctx) error {
	// Get the ID from the URL.
	appIDParam := c.Params("id")
	if appIDParam == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.MissingRequiredParam, "App ID is required.")
	}
	appID, err := utils.StringToUint(appIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the app.
	if deleted, err := services.IsAppDeleted(appID); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !deleted {
		return errorutil.Response(c, fiber.StatusNotFound, errors.AppExists, "App does not exist.")
	}

	// Restore the app.
	if err := services.RestoreApp(appID); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// validateAppSettings validates an array of DomainSetting structs.
// It checks if the Value field of each DomainSetting is valid based on its ValueType.
// If any validation errors occur, it returns a comma-separated string of error messages.
// If the string is empty, it means all validations passed.
func validateAppSettings(settings *[]requests.AppSetting) string {
	var validateErrors []string

	for i := range *settings {
		setting := &(*settings)[i]
		switch enums.ValueType(setting.ValueType) {
		case enums.Int:
			if _, err := strconv.Atoi(setting.Value); err != nil {
				validateErrors = append(validateErrors, fmt.Sprintf("Invalid int value for setting %s", setting.Name))
			}
		case enums.Float:
			if _, err := strconv.ParseFloat(setting.Value, 64); err != nil {
				validateErrors = append(validateErrors, fmt.Sprintf("Invalid float value for setting %s", setting.Name))
			}
		case enums.String:
			// No validation needed for string type.
		case enums.Bool:
			if _, err := strconv.ParseBool(setting.Value); err != nil {
				validateErrors = append(validateErrors, fmt.Sprintf("Invalid bool value for setting %s", setting.Name))
			}
		case enums.Date:
			if _, err := time.Parse(time.DateOnly, setting.Value); err != nil {
				validateErrors = append(validateErrors, fmt.Sprintf("Invalid date value for setting %s", setting.Name))
			}
		case enums.DateTime:
			if _, err := time.Parse(time.DateTime, setting.Value); err != nil {
				validateErrors = append(validateErrors, fmt.Sprintf("Invalid datetime value for setting %s", setting.Name))
			}
		case enums.JSON:
			var js json.RawMessage
			if err := json.Unmarshal([]byte(setting.Value), &js); err != nil {
				validateErrors = append(validateErrors, fmt.Sprintf("Invalid JSON value for setting %s", setting.Name))
			}
		default:
			validateErrors = append(validateErrors, fmt.Sprintf("Unknown ValueType for setting %s", setting.Name))
		}
	}

	return strings.Join(validateErrors, ", ")
}
