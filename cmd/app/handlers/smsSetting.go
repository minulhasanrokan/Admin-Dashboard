package handlers

import (
	"admin-dashboard/cmd/app/requests"
	"admin-dashboard/cmd/common/response/webResponse"
	"admin-dashboard/cmd/internal/constants"
	"admin-dashboard/cmd/internal/dtos"
	"admin-dashboard/cmd/internal/session"

	"github.com/labstack/echo/v4"
)

func (h Handler) SmsSettingPage(c echo.Context) error {

	smsSetting, err := h.Service.SmsSetting()

	if err != nil {
		return webResponse.SendNotFoundResponse(c, "SMS settings not found", err)
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

	html, err := webResponse.RenderHTML(c, "settings/sms_setting", map[string]any{
		"SmsSetting":   smsSetting,
		"RightDetails": rightDetails,
	})

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "SMS settings", html)
}

func (h Handler) SmsSettingUpdate(c echo.Context) error {

	payload := new(requests.SmsSettingRequest)

	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	smsSetting, err := h.Service.SmsSettingUpdate(payload)

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	err = h.Service.CreateUserActivity(c, constants.ActivityEdit, "sms_settings", smsSetting.Id, "Update SMS Setting")

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "SMS settings updated successfully", nil)
}

func (h Handler) SendTestSms(c echo.Context) error {

	payload := new(requests.SendTestSmsRequest)

	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	resp, err := h.Service.SendTestSms(payload)

	if err != nil {
		_ = h.Service.CreateUserActivity(c, constants.ActivitySendSms, "sms_settings", 0, "Failed to send test SMS to "+payload.Mobile)
		return webResponse.SendBadRequestResponse(c, err.Error(), nil)
	}

	_ = h.Service.CreateUserActivity(c, constants.ActivitySendSms, "sms_settings", 0, "Send test SMS to "+payload.Mobile)

	return webResponse.SendSuccessResponse(c, "Test SMS sent successfully!", resp)
}
