package webRoute

import (
	"admin-dashboard/cmd/internal/applicationData"

	"github.com/labstack/echo/v4"
)

func LoginRoutes(web *echo.Group, appData *applicationData.Application) {

	handler := appData.Handler

	login := web.Group("")
	login.Use(appData.Middleware.LoginCheck())

	login.GET("", handler.LoginPage)
	login.GET("/", handler.LoginPage)
	login.GET("login", handler.LoginPage)
	login.POST("login", handler.LoginStore)

	web.GET("logout", handler.Logout)
	web.GET("logout/", handler.Logout)
	web.POST("logout", handler.Logout)
}
