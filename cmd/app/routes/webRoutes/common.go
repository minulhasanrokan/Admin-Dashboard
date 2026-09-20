package webRoute

import (
	"admin-dashboard/cmd/internal/applicationData"

	"github.com/labstack/echo/v4"
)

func CommonRoutes(web *echo.Group, appData *applicationData.Application) {

	handler := appData.Handler

	common := web.Group("common")
	common.Use(appData.Middleware.RequireSession())

	common.GET("/get-dropdown", handler.GetDropdown)
}
