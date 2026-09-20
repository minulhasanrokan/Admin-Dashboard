package handlers

import (
	"admin-dashboard/cmd/app/requests"
	"admin-dashboard/cmd/common/helpers"
	"admin-dashboard/cmd/common/response/webResponse"
	"encoding/base64"
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
)

func (h Handler) ForgotPassword(c echo.Context) error {

	return h.RenderWithErrorCheck(c, nil, nil, "forgot-password")
}

func (h Handler) ForgotPasswordStore(c echo.Context) error {

	payload := new(requests.EmailRequest)

	validationErrors := h.PayloadValidate(c, payload)

	if validationErrors != nil {

		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	user, err := h.Repository.GetUserByEmail(payload.Email)

	if err != nil {

		if err.Error() == h.Db.NotFoundMessage() {

			return webResponse.SendBadRequestResponse(c, "Invalid email address", nil)
		}

		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	if user.DeleteStatus == 1 {
		return webResponse.SendBadRequestResponse(c, "Invalid email address", nil)
	}

	if !user.Verified {

		return webResponse.SendErrorResponse(c, "Account not verified", http.StatusForbidden, nil)
	}

	_, err = h.Repository.GetActivePasswordTokenUserId(user.Id)

	if err != nil {

		if err.Error() == h.Db.NotFoundMessage() {

			if err := h.SendMail.SendPasswordRestMail(user, h.Config); err != nil {

				return webResponse.SendInternalServerErrorResponse(c, err)
			}

		} else {
			return webResponse.SendInternalServerErrorResponse(c, err)
		}
	}

	return webResponse.SendSuccessResponse(c, "Password reset email has been sent successfully", map[string]interface{}{
		"uuid": user.UUID,
	})
}

func (h Handler) ForgotPasswordPage(c echo.Context) error {

	var uUId requests.UserUUIDRequest

	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil {

		return h.RenderWithErrorCheck(c, nil, err, "500-error")
	}

	userRole, err := h.Service.UserRole()

	if err != nil {

		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	user, err := h.Repository.GetUserByUuId(uUId.UUID)

	if err != nil {
		return h.RenderWithErrorCheck(c, nil, nil, "500-error")
	}

	if user.DeleteStatus == 1 {
		return h.RenderWithErrorCheck(c, nil, err, "deleted")
	}

	_, err = h.Repository.GetActivePasswordTokenUserId(user.Id)

	if err != nil {

		if err.Error() == h.Db.NotFoundMessage() {

			if err := h.SendMail.SendPasswordRestMail(user, h.Config); err != nil {

				return h.RenderWithErrorCheck(c, nil, err, "500-error")
			}
		} else {
			return h.RenderWithErrorCheck(c, nil, err, "500-error")
		}
	}

	var email string
	if user.Email != nil {
		email = helpers.MaskEmail(*user.Email)
	}

	data := map[string]interface{}{
		"Email":    email,
		"UuId":     user.UUID,
		"UserRole": userRole,
	}

	return h.RenderWithErrorCheck(c, data, nil, "forgot-password-page")
}

func (h Handler) UpdatePasswordPage(c echo.Context) error {

	var request requests.PasswordResetQuery

	if err := (&echo.DefaultBinder{}).BindQueryParams(c, &request); err != nil {

		return h.RenderWithErrorCheck(c, nil, err, "500-error")
	}

	userRole, err := h.Service.UserRole()

	if err != nil {

		return h.RenderWithErrorCheck(c, nil, nil, "500-error")
	}

	decodeEmail, err := base64.RawURLEncoding.DecodeString(request.RawEmail)

	if err != nil {

		return h.RenderWithErrorCheck(c, nil, err, "500-error")
	}

	user, err := h.Repository.GetUserByEmail(string(decodeEmail))

	if err != nil {
		return h.RenderWithErrorCheck(c, nil, nil, "500-error")
	}

	if user.DeleteStatus == 1 {
		return h.RenderWithErrorCheck(c, nil, err, "deleted")
	}

	var email string
	if user.Email != nil {
		email = helpers.MaskEmail(*user.Email)
	}

	data := map[string]interface{}{
		"Email":    email,
		"UuId":     user.UUID,
		"UserRole": userRole,
	}

	return h.RenderWithErrorCheck(c, data, err, "forgot-password-page")
}

func (h Handler) ReSendPasswordToken(c echo.Context) error {

	var uUId requests.UserUUIDRequest

	if err := (&echo.DefaultBinder{}).BindBody(c, &uUId); err != nil {

		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	user, err := h.Repository.GetUserByUuId(uUId.UUID)

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	if user.DeleteStatus == 1 {
		return webResponse.SendBadRequestResponse(c, "Invalid user account", nil)
	}

	if err := h.SendMail.SendPasswordRestMail(user, h.Config); err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Resent Token Successfully", nil)
}

func (h Handler) ResetPasswordStore(c echo.Context) error {

	payload := new(requests.UserPasswordResetRequest)

	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	user, err := h.Repository.GetUserByUuId(payload.UUID)

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	if user.DeleteStatus == 1 {
		return webResponse.SendBadRequestResponse(c, "Invalid user account", nil)
	}

	token, err := h.Repository.GetActivePasswordTokenUserId(user.Id)

	if err != nil {

		if err.Error() == h.Db.NotFoundMessage() {
			return webResponse.SendErrorResponse(c, "Password reset code has expired or is invalid. Please request a new OTP.", 400, nil)
		}
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	if token == nil {
		return webResponse.SendErrorResponse(c, "Password reset code has expired or is invalid. Please request a new OTP.", 400, nil)
	}

	finalToken := fmt.Sprintf("%s%s%s%s", payload.DigitOne, payload.DigitTwo, payload.DigitThree, payload.DigitFour)

	if finalToken != token.OtpToken {

		return webResponse.SendErrorResponse(c, "Token not matched", 400, nil)
	}

	err = h.Service.UpdateUserResetPasswordToken(token)

	if err != nil {

		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	err = h.Service.UpdateUserPassword(user, payload.Password)

	if err != nil {

		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	_ = h.SendMail.SendPasswordRestSuccessMail(user, h.Config)

	return webResponse.SendSuccessResponse(c, "User Password Reset Successfully", nil)
}
