package handlers

import (
	"admin-dashboard/cmd/app/requests"
	"admin-dashboard/cmd/common/helpers"
	"admin-dashboard/cmd/common/response/webResponse"
	"admin-dashboard/cmd/internal/dtos"
	"admin-dashboard/cmd/internal/session"
	"fmt"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func (h Handler) CreateAccountPage(c echo.Context) error {

	authSetting, _ := h.Service.AuthenticationSetting()

	userRole, err := h.Service.UserRole()

	if err != nil {

		return h.RenderWithErrorCheck(c, nil, err, "500-error")
	}

	data := map[string]interface{}{
		"UserRole":    userRole,
		"AuthSetting": authSetting,
	}

	return h.RenderWithErrorCheck(c, data, err, "create-account")
}

func (h Handler) CreateAccount(c echo.Context) error {

	authSetting, _ := h.Service.AuthenticationSetting()
	if authSetting != nil && !authSetting.AllowRegistration {
		return webResponse.SendBadRequestResponse(c, "Registration is currently disabled by administrator", nil)
	}

	payload := new(requests.UserRegisterRequest)

	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	if authSetting != nil && authSetting.InviteOnly {
		if strings.TrimSpace(payload.InviteCode) == "" {
			return webResponse.SendBadRequestResponse(c, "Valid invitation code is required to register", nil)
		}
	}

	if authSetting != nil && authSetting.PasswordRequired {
		if strings.TrimSpace(payload.Password) == "" {
			return webResponse.SendBadRequestResponse(c, "Password is required", nil)
		}
	}

	if existingUser, err := h.Repository.GetUserByEmail(payload.Email); err != nil {
		if err.Error() != h.Db.NotFoundMessage() {
			return webResponse.SendInternalServerErrorResponse(c, err)
		}
	} else if existingUser != nil {
		return webResponse.SendBadRequestResponse(c, "email already exists", nil)
	}

	if existingUser, err := h.Repository.GetUserByMobile(payload.Mobile); err != nil {
		if err.Error() != h.Db.NotFoundMessage() {
			return webResponse.SendInternalServerErrorResponse(c, err)
		}
	} else if existingUser != nil {
		return webResponse.SendBadRequestResponse(c, "mobile already exists", nil)
	}

	registeredUser, err := h.Service.RegisterUser(payload, c.RealIP())

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	if authSetting != nil && !authSetting.VerificationRequired {
		registeredUser.Verified = true
		_ = h.Service.UpdateUserVerifyStatus(registeredUser)

		return webResponse.SendCreatedSuccessResponse(c, "registration success", map[string]interface{}{
			"uuid":     registeredUser.UUID,
			"verified": true,
		})
	}

	expiryMinutes := uint64(h.Config.VerificationTokenMinute)
	if authSetting != nil && authSetting.VerificationExpired > 0 {
		expiryMinutes = authSetting.VerificationExpired
	}

	if err := h.SendMail.SendRegisterUserMail(registeredUser, h.Config, expiryMinutes); err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendCreatedSuccessResponse(c, "registration success", map[string]interface{}{
		"uuid":     registeredUser.UUID,
		"verified": false,
	})
}

func (h Handler) VerifyUserAccount(c echo.Context) error {

	var uUId requests.UserUUIDRequest

	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil {
		return h.RenderWithErrorCheck(c, nil, err, "500-error")
	}

	user, err := h.Repository.GetUserByUuId(uUId.UUID)

	if err != nil {
		return h.RenderWithErrorCheck(c, nil, nil, "500-error")
	}

	if user.DeleteStatus == 1 {
		return h.RenderWithErrorCheck(c, nil, err, "deleted")
	}

	authSetting, _ := h.Service.AuthenticationSetting()

	if user.Verified && (authSetting == nil || !authSetting.OtpSystem) {
		return h.RenderWithErrorCheck(c, nil, err, "already-verified")
	}

	token, err := h.Repository.GetActiveVerifyTokenUserId(user.Id)

	if err != nil {

		if err.Error() == h.Db.NotFoundMessage() {
			expiryMinutes := uint64(h.Config.VerificationTokenMinute)
			if authSetting != nil && authSetting.VerificationExpired > 0 {
				expiryMinutes = authSetting.VerificationExpired
			}

			if err := h.SendMail.SendRegisterUserMail(user, h.Config, expiryMinutes); err != nil {

				return h.RenderWithErrorCheck(c, nil, err, "500-error")
			}

			token, err = h.Repository.GetActiveVerifyTokenUserId(user.Id)

			if err != nil {
				return h.RenderWithErrorCheck(c, nil, nil, "500-error")
			}

		} else {
			return h.RenderWithErrorCheck(c, nil, err, "500-error")
		}
	}

	if token == nil {
		return h.RenderWithErrorCheck(c, nil, nil, "500-error")
	}

	var email string
	if user.Email != nil {
		email = helpers.MaskEmail(*user.Email)
	}

	data := map[string]interface{}{
		"Email":       email,
		"UuId":        user.UUID,
		"AuthSetting": authSetting,
	}

	return h.RenderWithErrorCheck(c, data, nil, "verification")
}

func (h Handler) ReSendUserVerificationToken(c echo.Context) error {

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

	authSetting, _ := h.Service.AuthenticationSetting()
	if user.Verified && (authSetting == nil || !authSetting.OtpSystem) {
		return webResponse.SendBadRequestResponse(c, "Account already verified", nil)
	}

	expiryMinutes := uint64(h.Config.VerificationTokenMinute)
	if authSetting != nil && authSetting.VerificationExpired > 0 {
		expiryMinutes = authSetting.VerificationExpired
	}

	if err := h.SendMail.SendRegisterUserMail(user, h.Config, expiryMinutes); err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	if authSetting != nil && authSetting.OtpType == "sms" && user.Mobile != nil && *user.Mobile != "" {
		tokenRecord, err := h.Repository.GetActiveVerifyTokenUserId(user.Id)
		if err == nil && tokenRecord != nil {
			_ = h.Service.SendOtpSms(*user.Mobile, tokenRecord.OtpToken)
		}
	}

	return webResponse.SendSuccessResponse(c, "Resent Token Successfully", nil)
}

func (h Handler) VerifyUserAccountStore(c echo.Context) error {

	payload := new(requests.UserVerifyEmailRequest)

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

	authSetting, _ := h.Service.AuthenticationSetting()
	if user.Verified && (authSetting == nil || !authSetting.OtpSystem) {
		return webResponse.SendBadRequestResponse(c, "Account already verified", nil)
	}

	token, err := h.Repository.GetActiveVerifyTokenUserId(user.Id)

	if err != nil {

		if err.Error() == h.Db.NotFoundMessage() {
			return webResponse.SendErrorResponse(c, "Verification code has expired or is invalid. Please request a new OTP.", 400, nil)
		}
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	if token == nil {
		return webResponse.SendErrorResponse(c, "Verification code has expired or is invalid. Please request a new OTP.", 400, nil)
	}

	finalToken := fmt.Sprintf("%s%s%s%s", payload.DigitOne, payload.DigitTwo, payload.DigitThree, payload.DigitFour)

	if finalToken != token.OtpToken {

		return webResponse.SendErrorResponse(c, "Token not matched", 400, nil)
	}

	err = h.Service.UpdateUserVerifyToken(token)

	if err != nil {

		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	if !user.Verified {
		err = h.Service.UpdateUserVerifyStatus(user)

		if err != nil {

			return webResponse.SendInternalServerErrorResponse(c, err)
		}

		_ = h.SendMail.SendVerifyUserMail(user, h.Config)
	}

	if authSetting != nil && authSetting.OtpSystem {
		var userName string
		var userEmail string
		var userPhoto string

		now := time.Now()

		if user.Email != nil {
			userEmail = *user.Email
		}

		if user.Name != nil {
			userName = *user.Name
		}

		if user.Photo != nil {
			userPhoto = *user.Photo
		}

		sessionData := session.UserSession{
			UserId:               user.Id,
			UUID:                 user.UUID,
			Email:                userEmail,
			UserName:             userName,
			UserType:             user.UserType,
			SuperAdmin:           user.SuperAdminStatus,
			IsActive:             user.IsActive,
			Verified:             true,
			LastLoginIP:          c.RealIP(),
			Photo:                userPhoto,
			LoginTime:            &now,
			Group:                user.Group,
			BranchID:             user.BranchId,
			Department:           user.Department,
			PasswordChangeStatus: user.PasswordChangeStatus,
			PasswordExpireDate:   user.PasswordExpireDate,
			PasswordChangeDate:   user.PasswordChangeDate,
			LastActiveTime:       &now,
			RouteAction:          "",
			RouteActionDetails:   dtos.RightInfo{},
			CategoryRouteAction:  "",
		}

		_ = session.SetUserSession(c, sessionData, h.Config)
		_ = h.Service.InsertUserLogHistory(c, user.Id)
		_ = h.Service.UpdateUserLastLoginTime(user)
	}

	return webResponse.SendSuccessResponse(c, "User Verified Successfully", nil)
}
