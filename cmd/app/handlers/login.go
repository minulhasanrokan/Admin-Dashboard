package handlers

import (
	"admin-dashboard/cmd/app/requests"
	"admin-dashboard/cmd/common/helpers"
	"admin-dashboard/cmd/common/response/webResponse"
	"admin-dashboard/cmd/internal/dtos"
	"admin-dashboard/cmd/internal/models"
	"admin-dashboard/cmd/internal/session"
	"net/http"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
)

func (h Handler) LoginPage(c echo.Context) error {

	err := session.ClearUserSession(c)

	if err != nil {

		return h.RenderWithErrorCheck(c, nil, err, "500-error")
	}

	userRole, err := h.Service.UserRole()

	if err != nil {

		return h.RenderWithErrorCheck(c, nil, err, "500-error")
	}

	authSetting, _ := h.Service.AuthenticationSetting()

	data := map[string]interface{}{
		"UserRole":    userRole,
		"AuthSetting": authSetting,
	}

	return h.RenderWithErrorCheck(c, data, err, "login")
}

func (h Handler) LoginStore(c echo.Context) error {

	authSetting, _ := h.Service.AuthenticationSetting()

	payload := new(requests.UserLoginRequest)

	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {

		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	var userData *models.User
	var err error

	loginType := "email"
	if authSetting != nil && authSetting.LoginType != "" {
		loginType = authSetting.LoginType
	}

	if loginType == "mobile" {
		mobile := strings.TrimSpace(payload.Mobile)
		if mobile == "" {
			mobile = strings.TrimSpace(payload.Email)
		}
		if mobile == "" {
			return webResponse.SendBadRequestResponse(c, "Mobile number is required", nil)
		}
		userData, err = h.Repository.GetUserByMobile(mobile)
	} else {
		email := strings.TrimSpace(payload.Email)
		if email == "" {
			email = strings.TrimSpace(payload.Mobile)
		}
		if email == "" {
			return webResponse.SendBadRequestResponse(c, "Email address is required", nil)
		}
		userData, err = h.Repository.GetUserByEmail(email)
	}

	if err != nil {

		if err.Error() == h.Db.NotFoundMessage() {

			return webResponse.SendBadRequestResponse(c, "Invalid username or password", nil)
		}

		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	passwordRequired := true
	if authSetting != nil {
		passwordRequired = authSetting.PasswordRequired
	}

	if passwordRequired {
		if userData.Password == nil || payload.Password == "" {
			return webResponse.SendBadRequestResponse(c, "Invalid username or password", nil)
		}

		passwordStatus, err := helpers.CheckPasswordHash(payload.Password, *userData.Password)

		if err != nil || !passwordStatus || userData.DeleteStatus == 1 {

			return webResponse.SendBadRequestResponse(c, "Invalid username or password", nil)
		}
	} else if userData.DeleteStatus == 1 {
		return webResponse.SendBadRequestResponse(c, "Invalid user account", nil)
	}

	verificationRequired := true
	if authSetting != nil {
		verificationRequired = authSetting.VerificationRequired
	}

	if verificationRequired && !userData.Verified {

		expiryMinutes := uint64(h.Config.VerificationTokenMinute)
		if authSetting != nil && authSetting.VerificationExpired > 0 {
			expiryMinutes = authSetting.VerificationExpired
		}

		if err := h.SendMail.SendRegisterUserMail(userData, h.Config, expiryMinutes); err != nil {

			return webResponse.SendInternalServerErrorResponse(c, err)
		}

		return webResponse.SendBadRequestResponse(c, "Please verify your account before logging in. Check your email and complete the verification process.", nil)
	}

	if authSetting != nil && authSetting.OtpSystem {
		expiryMinutes := uint64(h.Config.VerificationTokenMinute)
		if authSetting.VerificationExpired > 0 {
			expiryMinutes = authSetting.VerificationExpired
		}

		if err := h.SendMail.SendRegisterUserMail(userData, h.Config, expiryMinutes); err != nil {
			return webResponse.SendInternalServerErrorResponse(c, err)
		}

		if authSetting.OtpType == "sms" && userData.Mobile != nil && *userData.Mobile != "" {
			tokenRecord, err := h.Repository.GetActiveVerifyTokenUserId(userData.Id)
			if err == nil && tokenRecord != nil {
				_ = h.Service.SendOtpSms(*userData.Mobile, tokenRecord.OtpToken)
			}
		}

		return webResponse.SendSuccessResponse(c, "OTP sent successfully. Please verify to proceed.", map[string]interface{}{
			"require_otp": true,
			"uuid":        userData.UUID,
			"otp_type":    authSetting.OtpType,
		})
	}

	var userName string
	var userEmail string
	var userPhoto string

	now := time.Now()

	if userData.Email != nil {
		userEmail = *userData.Email
	}

	if userData.Name != nil {
		userName = *userData.Name
	}

	if userData.Photo != nil {
		userPhoto = *userData.Photo
	}

	sessionData := session.UserSession{
		UserId:               userData.Id,
		UUID:                 userData.UUID,
		Email:                userEmail,
		UserName:             userName,
		UserType:             userData.UserType,
		SuperAdmin:           userData.SuperAdminStatus,
		IsActive:             userData.IsActive,
		Verified:             userData.Verified,
		LastLoginIP:          c.RealIP(),
		Photo:                userPhoto,
		LoginTime:            &now,
		Group:                userData.Group,
		BranchID:             userData.BranchId,
		Department:           userData.Department,
		PasswordChangeStatus: userData.PasswordChangeStatus,
		PasswordExpireDate:   userData.PasswordExpireDate,
		PasswordChangeDate:   userData.PasswordChangeDate,
		LastActiveTime:       &now,
		RouteAction:          "",
		RouteActionDetails:   dtos.RightInfo{},
		CategoryRouteAction:  "",
	}

	err = session.SetUserSession(c, sessionData, h.Config)

	if err != nil {

		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	err = h.Service.InsertUserLogHistory(c, userData.Id)

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	err = h.Service.UpdateUserLastLoginTime(userData)

	if err != nil {

		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "User login Success", nil)
}

func (h Handler) Logout(c echo.Context) error {
	userSess, _, err := session.GetUserSession(c)
	if err == nil && !userSess.IsEmpty() {
		_ = h.Service.UpdateUserLogTime(userSess.UserId, 1, 0)
	}

	_ = session.ClearUserSession(c)

	if c.Request().Header.Get("X-Requested-With") == "XMLHttpRequest" || c.Request().Header.Get("Accept") == "application/json" {
		return webResponse.SendSuccessResponse(c, "Logout successful", map[string]string{
			"redirect": "/login",
		})
	}

	return c.Redirect(http.StatusSeeOther, "/login")
}
