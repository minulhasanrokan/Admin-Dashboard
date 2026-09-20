package webRoute

import (
	"admin-dashboard/cmd/internal/applicationData"

	"github.com/labstack/echo/v4"
)

func ProfileRoutes(web *echo.Group, appData *applicationData.Application) {

	handler := appData.Handler

	profile := web.Group("profile")
	profile.Use(appData.Middleware.RequireLoginSession())

	profile.GET("", handler.ProfilePage)
	profile.POST("/update", handler.ProfileUpdate)
	profile.POST("/change-password", handler.ProfileChangePassword)
	profile.POST("/change-photo", handler.ProfileChangePhoto)
}
