package handlers

import (
	"admin-dashboard/cmd/app/requests"
	"admin-dashboard/cmd/common/response/webResponse"
	"admin-dashboard/cmd/internal/constants"
	"admin-dashboard/cmd/internal/dtos"
	"admin-dashboard/cmd/internal/session"

	"github.com/labstack/echo/v4"
)

func (h Handler) EmailSettingPage(c echo.Context) error {
	emailSetting, err := h.Service.EmailSetting()
	if err != nil {
		return webResponse.SendNotFoundResponse(c, "Email settings not found", err)
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

	html, err := webResponse.RenderHTML(c, "settings/email_setting", map[string]any{
		"EmailSetting":            emailSetting,
		"RightDetails":            rightDetails,
		"AwsRegions":              constants.AwsRegions,
		"SmtpEncryptionProtocols": constants.SmtpEncryptionProtocols,
	})
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "System settings", html)
}

func (h Handler) EmailSettingUpdate(c echo.Context) error {
	payload := new(requests.EmailSettingRequest)
	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	emailSetting, err := h.Service.EmailSettingUpdate(payload)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	err = h.Service.CreateUserActivity(c, constants.ActivityEdit, "email_settings", emailSetting.Id, "Update Email General Setting")
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Email settings updated successfully", nil)
}

func (h Handler) SendGridSettingUpdate(c echo.Context) error {
	payload := new(requests.SendGridSettingRequest)
	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	emailSetting, err := h.Service.SendGridSettingUpdate(payload)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	_ = h.Service.CreateUserActivity(c, constants.ActivityEdit, "email_settings", emailSetting.Id, "Update SendGrid Setting")

	return webResponse.SendSuccessResponse(c, "SendGrid configuration saved successfully", nil)
}

func (h Handler) SMTPSettingUpdate(c echo.Context) error {
	payload := new(requests.SMTPSettingRequest)
	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	emailSetting, err := h.Service.SMTPSettingUpdate(payload)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	_ = h.Service.CreateUserActivity(c, constants.ActivityEdit, "email_settings", emailSetting.Id, "Update SMTP Setting")

	return webResponse.SendSuccessResponse(c, "SMTP configuration saved successfully", nil)
}

func (h Handler) SESSettingUpdate(c echo.Context) error {
	payload := new(requests.SESSettingRequest)
	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	emailSetting, err := h.Service.SESSettingUpdate(payload)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	_ = h.Service.CreateUserActivity(c, constants.ActivityEdit, "email_settings", emailSetting.Id, "Update AWS SES Setting")

	return webResponse.SendSuccessResponse(c, "AWS SES configuration saved successfully", nil)
}

func (h Handler) SendTestEmail(c echo.Context) error {
	payload := new(requests.SendTestEmailRequest)
	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	if err := h.Service.SendTestEmail(payload); err != nil {
		_ = h.Service.CreateUserActivity(c, constants.ActivitySendEmail, "email_settings", 0, "Failed to send test email to "+payload.RecipientEmail)
		return webResponse.SendBadRequestResponse(c, err.Error(), nil)
	}

	_ = h.Service.CreateUserActivity(c, constants.ActivitySendEmail, "email_settings", 0, "Send test email to "+payload.RecipientEmail)

	return webResponse.SendSuccessResponse(c, "Test email sent successfully!", nil)
}
