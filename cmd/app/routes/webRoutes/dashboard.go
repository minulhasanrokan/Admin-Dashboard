package webRoute

import (
	"admin-dashboard/cmd/internal/applicationData"

	"github.com/labstack/echo/v4"
)

func DashboardRoutes(web *echo.Group, appData *applicationData.Application) {

	handler := appData.Handler

	dashboard := web.Group("dashboard")
	dashboard.Use(appData.Middleware.RequireSession())

	dashboard.GET("", handler.Dashboard)
}
