package handlers

import (
	"admin-dashboard/cmd/app/requests"
	"admin-dashboard/cmd/common/response/webResponse"
	"admin-dashboard/cmd/internal/constants"
	"admin-dashboard/cmd/internal/dtos"
	"admin-dashboard/cmd/internal/session"

	"github.com/labstack/echo/v4"
)

func (h Handler) UserRoleSettingPage(c echo.Context) error {

	userRole, err := h.Service.UserRole()

	if err != nil {
		return webResponse.SendNotFoundResponse(c, "User Role not found", err)
	}

	userSess, _, err := session.GetUserSession(c)

	if err != nil || userSess.IsEmpty() || !userSess.IsActive || userSess.CategoryRouteAction == "" {
		_ = session.ClearUserSession(c)
		return webResponse.SendUnauthorizedResponse(c, "Session expired", nil)
	}

	var rightInfo dtos.RightInfo

	rightInfo, err = h.Service.GetRightDetailsByKey(userSess.RouteAction)

	if err != nil {

		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	var rightDetails []dtos.RightInfo

	rightDetails, err = h.Service.GetRightDetailsByCategoryKeyKey(rightInfo.CategoryKey)

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	html, err := webResponse.RenderHTML(c, "settings/user_role", map[string]any{"UserRole": userRole, "RightDetails": rightDetails})

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "User Role settings", html)
}

func (h Handler) UserRoleSettingUpdate(c echo.Context) error {

	payload := new(requests.UserRoleRequest)

	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	emailSetting, err := h.Service.UserRoleSettingUpdate(payload)

	if err != nil {

		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	err = h.Service.CreateUserActivity(c, constants.ActivityEdit, "user_roles", emailSetting.Id, "Update User Role Setting")

	if err != nil {

		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "User Role settings Updated", nil)
}
