package handlers

import (
	"admin-dashboard/cmd/app/requests"
	"admin-dashboard/cmd/common/response/webResponse"
	"admin-dashboard/cmd/internal/constants"
	"admin-dashboard/cmd/internal/dtos"
	"admin-dashboard/cmd/internal/session"

	"github.com/labstack/echo/v4"
)

func (h Handler) AuthenticationSettingPage(c echo.Context) error {

	authSetting, err := h.Service.AuthenticationSetting()

	if err != nil {
		return webResponse.SendNotFoundResponse(c, "Authentication settings not found", err)
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

	html, err := webResponse.RenderHTML(c, "settings/authentication-settings", map[string]any{
		"AuthSetting":  authSetting,
		"RightDetails": rightDetails,
		"ConfigData":   h.Config,
	})

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Authentication Settings", html)
}

func (h Handler) AuthenticationSettingUpdate(c echo.Context) error {

	payload := new(requests.AuthenticationSettingRequest)

	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	authSetting, err := h.Service.AuthenticationSettingUpdate(payload)

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	err = h.Service.CreateUserActivity(c, constants.ActivityEdit, "authentication_settings", authSetting.Id, "Update Authentication Setting")

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Authentication settings updated", nil)
}
