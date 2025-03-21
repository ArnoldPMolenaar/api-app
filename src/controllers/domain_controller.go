package controllers

import (
	"api-app/main/src/dto/requests"
	"api-app/main/src/dto/responses"
	"api-app/main/src/enums"
	"api-app/main/src/errors"
	"api-app/main/src/services"
	"encoding/json"
	"fmt"
	errorutil "github.com/ArnoldPMolenaar/api-utils/errors"
	"github.com/ArnoldPMolenaar/api-utils/utils"
	"github.com/gofiber/fiber/v2"
	"strconv"
	"strings"
	"time"
)

// GetDomain function fetches a domain from the database by its ID.
func GetDomain(c *fiber.Ctx) error {
	// Get the domainID parameter from the URL.
	domainIDParam := c.Params("id")
	if domainIDParam == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.MissingRequiredParam, "Domain ID is required.")
	}
	domainID, err := utils.StringToUint(domainIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, "Invalid Domain ID.")
	}

	// Get the domain.
	domain, err := services.GetDomainById(domainID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if domain.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errorutil.NotFound, "Domain not found.")
	}

	// Return the domain.
	response := responses.Domain{}
	response.SetDomain(domain)

	return c.JSON(response)
}

// CreateDomain func to create a domain.
func CreateDomain(c *fiber.Ctx) error {
	// Parse the request.
	request := requests.CreateDomain{}
	if err := c.BodyParser(&request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate domain fields.
	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, utils.ValidatorErrors(err))
	}
	if validationErrors := validateSettings(&request.Settings); validationErrors != "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.DomainSettings, validationErrors)
	}

	// Check if domain exists.
	if available, err := services.IsDomainNameAvailable(request.AppID, request.Name); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if available {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.DomainAvailable, "DomainName already available.")
	}

	// Create the domain.
	domain, err := services.CreateDomain(request.AppID, request.SSL, request.Name, request.IpAddress, &request.Settings)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the domain.
	response := responses.Domain{}
	response.SetDomain(domain)

	return c.JSON(response)
}

// UpdateDomain func to update a domain.
func UpdateDomain(c *fiber.Ctx) error {
	// Get the ID from the URL.
	domainIDParam := c.Params("id")
	if domainIDParam == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.MissingRequiredParam, "Domain ID is required.")
	}
	domainID, err := utils.StringToUint(domainIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Parse the request.
	request := requests.UpdateDomain{}
	if err := c.BodyParser(&request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.BodyParse, err.Error())
	}

	// Validate domain fields.
	validate := utils.NewValidator()
	if err := validate.Struct(request); err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.Validator, utils.ValidatorErrors(err))
	}
	if validationErrors := validateSettings(&request.Settings); validationErrors != "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errors.DomainSettings, validationErrors)
	}

	// Get the domain.
	domain, err := services.GetDomainById(domainID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if domain.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.DomainExists, "Domain not found.")
	}

	// Check if domain exists.
	if request.Name != domain.Name {
		if available, err := services.IsDomainNameAvailable(domain.AppID, request.Name); err != nil {
			return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
		} else if available {
			return errorutil.Response(c, fiber.StatusBadRequest, errors.DomainAvailable, "DomainName already available.")
		}
	}

	// Check if the domain data has been modified since it was last fetched.
	if request.UpdatedAt.Unix() < domain.UpdatedAt.Unix() {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.OutOfSync, "Data is out of sync.")
	}

	// Update the domain.
	domain, err = services.UpdateDomain(domain, request.SSL, request.Name, request.IpAddress, &request.Settings)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	// Return the domain.
	response := responses.Domain{}
	response.SetDomain(domain)

	return c.JSON(response)
}

// DeleteDomain func to delete a domain.
func DeleteDomain(c *fiber.Ctx) error {
	// Get the ID from the URL.
	domainIDParam := c.Params("id")
	if domainIDParam == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.MissingRequiredParam, "Domain ID is required.")
	}
	domainID, err := utils.StringToUint(domainIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Get the domain.
	domain, err := services.GetDomainById(domainID)
	if err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if domain.ID == 0 {
		return errorutil.Response(c, fiber.StatusNotFound, errors.DomainExists, "Domain not found.")
	}

	// Delete the domain.
	if err := services.DeleteDomain(domain); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// RestoreDomain func to restore a domain.
func RestoreDomain(c *fiber.Ctx) error {
	// Get the ID from the URL.
	domainIDParam := c.Params("id")
	if domainIDParam == "" {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.MissingRequiredParam, "Domain ID is required.")
	}
	domainID, err := utils.StringToUint(domainIDParam)
	if err != nil {
		return errorutil.Response(c, fiber.StatusBadRequest, errorutil.InvalidParam, err.Error())
	}

	// Find the domain.
	if deleted, err := services.IsDomainDeleted(domainID); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	} else if !deleted {
		return errorutil.Response(c, fiber.StatusNotFound, errors.DomainExists, "Domain does not exist.")
	}

	// Restore the domain.
	if err := services.RestoreDomain(domainID); err != nil {
		return errorutil.Response(c, fiber.StatusInternalServerError, errorutil.QueryError, err.Error())
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// validateSettings validates an array of DomainSetting structs.
// It checks if the Value field of each DomainSetting is valid based on its ValueType.
// If any validation errors occur, it returns a comma-separated string of error messages.
// If the string is empty, it means all validations passed.
func validateSettings(settings *[]requests.DomainSetting) string {
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
