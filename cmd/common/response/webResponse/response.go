package webResponse

import (
	"admin-dashboard/cmd/common/response"
	"bytes"
	"fmt"
	"log"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func extractAndSetCSRF(c echo.Context) string {

	var csrfTokenStr string

	if csrfVal := c.Get("_csrf"); csrfVal != nil {
		if str, ok := csrfVal.(string); ok {
			csrfTokenStr = str
		}
	}

	if csrfTokenStr == "" {
		if forcedCsrf := c.Get(middleware.DefaultCSRFConfig.ContextKey); forcedCsrf != nil {
			if str, ok := forcedCsrf.(string); ok {
				csrfTokenStr = str
			}
		}
	}

	c.Response().Header().Set("X-CSRF-Token", csrfTokenStr)
	return csrfTokenStr
}

func SendSuccessResponse(c echo.Context, message string, data interface{}) error {

	csrfTokenStr := extractAndSetCSRF(c)

	return c.JSON(http.StatusOK, response.WebJsonSuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
		CSRF:    csrfTokenStr,
	})
}

func SendCreatedSuccessResponse(c echo.Context, message string, data interface{}) error {

	csrfTokenStr := extractAndSetCSRF(c)

	return c.JSON(http.StatusCreated, response.WebJsonSuccessResponse{
		Success: true,
		Message: message,
		Data:    data,
		CSRF:    csrfTokenStr,
	})
}

func SendFailedValidationResponse(c echo.Context, errors []*response.ValidationError) error {
	csrfTokenStr := extractAndSetCSRF(c)

	return c.JSON(http.StatusUnprocessableEntity, response.WebJsonFailedValidationResponse{
		Success: false,
		Errors:  errors,
		Message: "Validation Failed",
		CSRF:    csrfTokenStr,
	})
}

func SendErrorResponse(c echo.Context, message string, statusCode int, err error) error {

	csrfTokenStr := extractAndSetCSRF(c)

	var errText string

	if err != nil {

		log.Printf("[ERROR] %s: %v", message, err.Error())

		errText = fmt.Sprintf("[ERROR] %s: %v", message, err.Error())
	} else {
		errText = message
	}

	return c.JSON(statusCode, response.WebJsonErrorResponse{
		Success: false,
		Message: message,
		CSRF:    csrfTokenStr,
		Error:   errText,
	})
}

func SendBadRequestResponse(c echo.Context, message string, err error) error {
	return SendErrorResponse(c, message, http.StatusBadRequest, err)
}

func SendNotFoundResponse(c echo.Context, message string, err error) error {
	return SendErrorResponse(c, message, http.StatusNotFound, err)
}

func SendCsrfErrorResponse(c echo.Context, message string, err error) error {
	return SendErrorResponse(c, message, http.StatusForbidden, err)
}

func SendInternalServerErrorResponse(c echo.Context, err error) error {
	return SendErrorResponse(c, "Internal Server Error", http.StatusInternalServerError, err)
}

func SendUnauthorizedResponse(c echo.Context, message string, err error) error {
	if message == "" {
		message = "Unauthorized"
	}
	return SendErrorResponse(c, message, http.StatusUnauthorized, err)
}

func RenderHTML(c echo.Context, template string, data any) (string, error) {

	var buf bytes.Buffer

	err := c.Echo().Renderer.Render(&buf, template+".gohtml", data, c)

	if err != nil {
		return "", err
	}

	return buf.String(), nil
}
