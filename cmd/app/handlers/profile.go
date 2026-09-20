package handlers

import (
	"admin-dashboard/cmd/app/requests"
	"admin-dashboard/cmd/common/helpers"
	"admin-dashboard/cmd/common/response"
	"admin-dashboard/cmd/common/response/webResponse"
	"admin-dashboard/cmd/internal/constants"
	"admin-dashboard/cmd/internal/session"
	"fmt"
	"path/filepath"
	"time"

	"github.com/labstack/echo/v4"
)

func (h Handler) ProfilePage(c echo.Context) error {

	userSess, _, err := session.GetUserSession(c)

	if err != nil || userSess.IsEmpty() || !userSess.IsActive {
		_ = session.ClearUserSession(c)
		return webResponse.SendUnauthorizedResponse(c, "Session expired", nil)
	}

	user, err := h.Service.GetUserByID(userSess.UserId)

	if err != nil {
		return webResponse.SendNotFoundResponse(c, "User not found", err)
	}

	userRole, err := h.Service.UserRole()
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	data := map[string]any{
		"User":       user,
		"UserRole":   userRole,
		"ConfigData": h.Config,
	}

	html, err := webResponse.RenderHTML(c, "profile/profile_view", data)

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "My Profile", html)
}

func (h Handler) ProfileUpdate(c echo.Context) error {

	userSess, sess, err := session.GetUserSession(c)

	if err != nil || userSess.IsEmpty() || !userSess.IsActive {
		_ = session.ClearUserSession(c)
		return webResponse.SendUnauthorizedResponse(c, "Session expired", nil)
	}

	payload := new(requests.ProfileUpdateRequest)

	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	updatedUser, err := h.Service.UpdateUserProfile(userSess.UserId, payload.Name, payload.Mobile, payload.Address, payload.Details)

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	if updatedUser.Name != nil && *updatedUser.Name != "" {
		sess.Values[constants.SessionKeyUserName] = *updatedUser.Name
		_ = sess.Save(c.Request(), c.Response())
	}

	_ = h.Service.CreateUserActivity(c, constants.ActivityEdit, "users", userSess.UserId, "Update Profile Information")

	return webResponse.SendSuccessResponse(c, "Profile updated successfully", nil)
}

func (h Handler) ProfileChangePassword(c echo.Context) error {

	userSess, _, err := session.GetUserSession(c)

	if err != nil || userSess.IsEmpty() || !userSess.IsActive {
		_ = session.ClearUserSession(c)
		return webResponse.SendUnauthorizedResponse(c, "Session expired", nil)
	}

	payload := new(requests.ProfilePasswordRequest)

	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	user, err := h.Service.GetUserByID(userSess.UserId)

	if err != nil {
		return webResponse.SendNotFoundResponse(c, "User not found", err)
	}

	if user.Password == nil {
		return webResponse.SendBadRequestResponse(c, "Password not set for user", nil)
	}

	isValid, err := helpers.CheckPasswordHash(payload.CurrentPassword, *user.Password)
	if err != nil || !isValid {
		valErr := []*response.ValidationError{
			{
				Key:       "current_password",
				Error:     "The current password entered is incorrect",
				Condition: "invalid",
			},
		}
		return webResponse.SendFailedValidationResponse(c, valErr)
	}

	err = h.Service.UpdateUserPassword(user, payload.NewPassword)

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	_ = h.Service.CreateUserActivity(c, constants.ActivityEdit, "users", userSess.UserId, "Change User Password")

	return webResponse.SendSuccessResponse(c, "Password changed successfully", nil)
}

func (h Handler) ProfileChangePhoto(c echo.Context) error {

	userSess, sess, err := session.GetUserSession(c)

	if err != nil || userSess.IsEmpty() || !userSess.IsActive {
		_ = session.ClearUserSession(c)
		return webResponse.SendUnauthorizedResponse(c, "Session expired", nil)
	}

	payload := new(requests.ProfilePhotoRequest)

	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	if payload.Photo == nil {
		return webResponse.SendBadRequestResponse(c, "No photo uploaded", nil)
	}

	photoName := fmt.Sprintf("profile_%d_%s%s", userSess.UserId, time.Now().Format("20060102150405"), filepath.Ext(payload.Photo.Filename))

	_, err = h.Storage.UploadFile("user/"+photoName, payload.Photo)

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	_, err = h.Service.UpdateUserPhoto(userSess.UserId, photoName)

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	sess.Values[constants.SessionKeyPhoto] = photoName
	_ = sess.Save(c.Request(), c.Response())

	_ = h.Service.CreateUserActivity(c, constants.ActivityEdit, "users", userSess.UserId, "Update Profile Photo")

	return webResponse.SendSuccessResponse(c, "Profile photo updated successfully", map[string]string{
		"photo": photoName,
		"url":   fmt.Sprintf("%s/file/user/%s", h.Config.AppBaseUrl, photoName),
	})
}
