package webRoute

import (
	"admin-dashboard/cmd/internal/applicationData"

	"github.com/labstack/echo/v4"
)

func SettingRoutes(web *echo.Group, appData *applicationData.Application) {

	handler := appData.Handler

	settings := web.Group("settings")

	settings.GET("/system", handler.SystemSetting, appData.Middleware.CheckCategoryPermission("settings.system"))

	settings.GET("/business-settings", handler.BusinessSettingPage, appData.Middleware.CheckPermission("settings.system.basic.settings"))
	settings.POST("/business-settings", handler.BusinessSystemSettingUpdate, appData.Middleware.CheckPermission("settings.system.basic.settings"))

	settings.GET("/images-settings", handler.ImagesSettingPage, appData.Middleware.CheckPermission("settings.system.basic.settings"))
	settings.POST("/images-settings", handler.ImagesSystemSettingUpdate, appData.Middleware.CheckPermission("settings.system.basic.settings"))

	settings.GET("/authentication-settings", handler.AuthenticationSettingPage, appData.Middleware.CheckPermission("settings.system.basic.settings"))
	settings.POST("/authentication-settings", handler.AuthenticationSettingUpdate, appData.Middleware.CheckPermission("settings.system.basic.settings"))

	settings.GET("/email-settings", handler.EmailSettingPage, appData.Middleware.CheckPermission("settings.system.system.settings"))
	settings.POST("/email-settings", handler.EmailSettingUpdate, appData.Middleware.CheckPermission("settings.system.system.settings"))
	settings.POST("/email-settings/sendgrid", handler.SendGridSettingUpdate, appData.Middleware.CheckPermission("settings.system.system.settings"))
	settings.POST("/email-settings/smtp", handler.SMTPSettingUpdate, appData.Middleware.CheckPermission("settings.system.system.settings"))
	settings.POST("/email-settings/ses", handler.SESSettingUpdate, appData.Middleware.CheckPermission("settings.system.system.settings"))
	settings.POST("/email-settings/test", handler.SendTestEmail, appData.Middleware.CheckPermission("settings.system.system.settings"))

	settings.GET("/sms-settings", handler.SmsSettingPage, appData.Middleware.CheckPermission("settings.system.system.settings"))
	settings.POST("/sms-settings", handler.SmsSettingUpdate, appData.Middleware.CheckPermission("settings.system.system.settings"))
	settings.POST("/sms-settings/test", handler.SendTestSms, appData.Middleware.CheckPermission("settings.system.system.settings"))

	settings.GET("/user-roles-settings", handler.UserRoleSettingPage, appData.Middleware.CheckPermission("settings.system.system.settings"))
	settings.POST("/user-roles-settings", handler.UserRoleSettingUpdate, appData.Middleware.CheckPermission("settings.system.system.settings"))

	settings.GET("/clear-cache", handler.ClearCachePage, appData.Middleware.CheckPermission("settings.system.other.settings"))
	settings.POST("/clear-cache", handler.ClearCache, appData.Middleware.CheckPermission("settings.system.other.settings"))

	settings.GET("/backup", handler.BackUpPage, appData.Middleware.CheckPermission("settings.system.other.settings"))
	settings.POST("/backup", handler.Backup, appData.Middleware.CheckPermission("settings.system.other.settings"))
	settings.GET("/backup/list", handler.BackupList, appData.Middleware.CheckPermission("settings.system.other.settings"))
	settings.POST("/backup/list", handler.BackupList, appData.Middleware.CheckPermission("settings.system.other.settings"))
	settings.GET("/backup/download/:uuid", handler.BackupDownload, appData.Middleware.CheckPermission("settings.system.other.settings"))
}
