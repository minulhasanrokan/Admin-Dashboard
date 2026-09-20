package webRoute

import (
	"admin-dashboard/cmd/internal/applicationData"

	"github.com/labstack/echo/v4"
)

func UserRoutes(web *echo.Group, appData *applicationData.Application) {

	handler := appData.Handler

	user := web.Group("user-registration")
	user.Use(appData.Middleware.ClearSession())

	user.GET("/create-account", handler.CreateAccountPage)
	user.POST("/store", handler.CreateAccount)

	user.GET("/verify-account/:uuid", handler.VerifyUserAccount)
	user.POST("/verify-store", handler.VerifyUserAccountStore)
	user.POST("/resend-token", handler.ReSendUserVerificationToken)

	password := web.Group("user-password")
	password.Use(appData.Middleware.ClearSession())

	password.GET("/forgot-password", handler.ForgotPassword)
	password.GET("/forgot-password/:uuid", handler.ForgotPasswordPage)
	password.GET("/update-password", handler.UpdatePasswordPage)

	password.POST("/forgot-password-store", handler.ForgotPasswordStore)
	password.POST("/resend-password-token", handler.ReSendPasswordToken)
	password.POST("/reset-password", handler.ResetPasswordStore)
}
