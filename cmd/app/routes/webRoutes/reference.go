package webRoute

import (
	"admin-dashboard/cmd/internal/applicationData"

	"github.com/labstack/echo/v4"
)

func ReferenceRoutes(web *echo.Group, appData *applicationData.Application) {

	handler := appData.Handler

	reference := web.Group("reference")

	// Category route for Designation (sidebar link: /reference/designation)
	reference.GET("/designation", handler.DesignationPage, appData.Middleware.CheckCategoryPermission("reference.designation"))
	reference.GET("/designation/add", handler.DesignationAddPage, appData.Middleware.CheckPermission("reference.designation.add"))
	reference.POST("/designation/add", handler.DesignationAdd, appData.Middleware.CheckPermission("reference.designation.add"))
	reference.GET("/designation/view", handler.DesignationPage, appData.Middleware.CheckPermission("reference.designation.view"))
	reference.POST("/designation/view", handler.DesignationList, appData.Middleware.CheckPermission("reference.designation.view"))
	reference.GET("/designation/view/:uuid", handler.DesignationDetailsPage, appData.Middleware.CheckPermission("reference.designation.view"))
	reference.GET("/designation/edit", handler.DesignationPage, appData.Middleware.CheckPermission("reference.designation.edit"))
	reference.POST("/designation/edit", handler.DesignationList, appData.Middleware.CheckPermission("reference.designation.edit"))
	reference.GET("/designation/edit/:uuid", handler.DesignationEditPage, appData.Middleware.CheckPermission("reference.designation.edit"))
	reference.POST("/designation/edit/:uuid", handler.DesignationEdit, appData.Middleware.CheckPermission("reference.designation.edit"))
	reference.GET("/designation/delete", handler.DesignationPage, appData.Middleware.CheckPermission("reference.designation.delete"))
	reference.POST("/designation/delete", handler.DesignationList, appData.Middleware.CheckPermission("reference.designation.delete"))
	reference.GET("/designation/delete/:uuid", handler.DesignationDeletePage, appData.Middleware.CheckPermission("reference.designation.delete"))
	reference.POST("/designation/delete/:uuid", handler.DesignationDelete, appData.Middleware.CheckPermission("reference.designation.delete"))
	reference.GET("/designation/stc", handler.DesignationPage, appData.Middleware.CheckPermission("reference.designation.stc"))
	reference.POST("/designation/stc", handler.DesignationList, appData.Middleware.CheckPermission("reference.designation.stc"))
	reference.GET("/designation/stc/:uuid", handler.DesignationStcPage, appData.Middleware.CheckPermission("reference.designation.stc"))
	reference.POST("/designation/stc/:uuid", handler.DesignationStc, appData.Middleware.CheckPermission("reference.designation.stc"))
	reference.GET("/designation/verify", handler.DesignationPage, appData.Middleware.CheckPermission("reference.designation.verify"))
	reference.POST("/designation/verify", handler.DesignationList, appData.Middleware.CheckPermission("reference.designation.verify"))
	reference.GET("/designation/verify/:uuid", handler.DesignationVerifyPage, appData.Middleware.CheckPermission("reference.designation.verify"))
	reference.POST("/designation/verify/:uuid", handler.DesignationVerify, appData.Middleware.CheckPermission("reference.designation.verify"))
	reference.GET("/designation/active", handler.DesignationPage, appData.Middleware.CheckPermission("reference.designation.active"))
	reference.POST("/designation/active", handler.DesignationList, appData.Middleware.CheckPermission("reference.designation.active"))
	reference.GET("/designation/active/:uuid", handler.DesignationActivePage, appData.Middleware.CheckPermission("reference.designation.active"))
	reference.POST("/designation/active/:uuid", handler.DesignationActive, appData.Middleware.CheckPermission("reference.designation.active"))
	reference.GET("/designation/deactivate", handler.DesignationPage, appData.Middleware.CheckPermission("reference.designation.deactivate"))
	reference.POST("/designation/deactivate", handler.DesignationList, appData.Middleware.CheckPermission("reference.designation.deactivate"))
	reference.GET("/designation/deactivate/:uuid", handler.DesignationDeactivatePage, appData.Middleware.CheckPermission("reference.designation.deactivate"))
	reference.POST("/designation/deactivate/:uuid", handler.DesignationDeactivate, appData.Middleware.CheckPermission("reference.designation.deactivate"))
	reference.GET("/designation/history/:uuid", handler.DesignationHistory, appData.Middleware.CheckPermission("reference.designation.view"))
	reference.GET("/designation/activity/:uuid", handler.DesignationActivity, appData.Middleware.CheckPermission("reference.designation.view"))
}
