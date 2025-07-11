package utils

import (
	"caisse-app-scaled/caisse_app_scaled/logger"
	"caisse-app-scaled/caisse_app_scaled/models"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
)

func FAILURE_ERR(err error) string { return "Failed Request : " + err.Error() }
func SYNTAX_ERR(source string, value any) string {
	return fmt.Sprintf("%s : %v could not be parsed, check syntax", source, value)
}
func NOTFOUND_ERR(source string, value any) string {
	return fmt.Sprintf("%s : %v was not found", source, value)
}

func GetApiError(c *fiber.Ctx, message string, status int) error {
	errorMsg := ""
	switch status {
	case 400:
		errorMsg = "Bad Request"
	case 401:
		errorMsg = "Unauthorized"
	case 403:
		errorMsg = "Forbidden"
	case 404:
		errorMsg = "Not Found "
	case 500:
		errorMsg = "Internal Server Error"
	}
	logger.Error(message)
	return c.Status(status).JSON(models.ApiError{
		Timestamp: time.Now(),
		Status:    status,
		Error:     errorMsg,
		Message:   message,
		Success:   false,
		Path:      c.Route().Path,
	})
}

func GetApiSuccess(cache map[string]any, c *fiber.Ctx, status int) error {
	message := ""
	switch status {
	case 200:
		message = "Ok"
	case 201:
		message = "Created"
	}
	path := c.Method() + " " + c.Path()
	now := time.Now()
	cache["t - "+path] = now
	cache[path] = models.ApiSuccess{
		Timestamp: time.Now(),
		Success:   true,
		Status:    status,
		Message:   message,
		Path:      c.Route().Path,
	}
	return c.Status(status).JSON(cache[path].(models.ApiSuccess))
}

func API_MAGASIN() string {
	if os.Getenv("ENVTEST") == "TRUE" {
		return "http://" + os.Getenv("GATEWAY") + ":8080/magasin"
	}
	return "http://" + os.Getenv("GATEWAY") + "/magasin"
}
func API_MERE() string {
	if os.Getenv("ENVTEST") == "TRUE" {
		return "http://" + os.Getenv("GATEWAY") + ":8090/mere"
	}
	return "http://" + os.Getenv("GATEWAY") + "/mere"
}
func API_LOGISTIC() string {
	if os.Getenv("ENVTEST") == "TRUE" {
		return "http://" + os.Getenv("GATEWAY") + ":8091/logistique"
	}
	return "http://" + os.Getenv("GATEWAY") + "/logistique"
}

func Errnotnil(err error) {
	if err != nil {
		logger.Error(err.Error())
	}
}
