package webRoute

import (
	"admin-dashboard/cmd/internal/applicationData"

	"github.com/labstack/echo/v4"
)

func ReferenceRoutes(web *echo.Group, appData *applicationData.Application) {

	handler := appData.Handler

	reference := web.Group("reference")

	// Category route for Designation
	reference.GET("/designation", handler.DesignationPage, appData.Middleware.CheckCategoryPermission("reference.designation"))

	// Detail actions for Designation
	reference.GET("/designation/view", handler.DesignationPage, appData.Middleware.CheckPermission("reference.designation.view"))
	reference.GET("/designation/list", handler.DesignationList, appData.Middleware.CheckPermission("reference.designation.view"))
	reference.POST("/designation/list", handler.DesignationList, appData.Middleware.CheckPermission("reference.designation.view"))
	reference.POST("/designation/add", handler.DesignationAdd, appData.Middleware.CheckPermission("reference.designation.add"))
	reference.POST("/designation/edit/:id", handler.DesignationEdit, appData.Middleware.CheckPermission("reference.designation.edit"))
	reference.POST("/designation/delete/:id", handler.DesignationDelete, appData.Middleware.CheckPermission("reference.designation.delete"))
	reference.POST("/designation/active/:id", handler.DesignationActive, appData.Middleware.CheckPermission("reference.designation.active"))
	reference.POST("/designation/deactivate/:id", handler.DesignationDeactivate, appData.Middleware.CheckPermission("reference.designation.deactivate"))
	reference.POST("/designation/stc/:id", handler.DesignationSTC, appData.Middleware.CheckPermission("reference.designation.stc"))
	reference.POST("/designation/verify/:id", handler.DesignationVerify, appData.Middleware.CheckPermission("reference.designation.verify"))
}
