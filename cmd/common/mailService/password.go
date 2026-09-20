package mailService

import (
	"admin-dashboard/cmd/common/helpers"
	"admin-dashboard/cmd/internal/config"
	"admin-dashboard/cmd/internal/constants"
	"admin-dashboard/cmd/internal/models"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"net/url"
	"time"
)

func (m SendMailService) SendPasswordRestMail(userData *models.User, cfg config.Config) error {

	token, err := helpers.GeneratePasswordResetToken(cfg.PasswordResetTokenLength)

	if err != nil {
		return err
	}

	hash := sha256.Sum256([]byte(token))

	tokenData := models.PasswordResetToken{
		UserId:    userData.Id,
		TokenHash: fmt.Sprintf("%x", hash),
		OtpToken:  token,
		ExpiresAt: time.Now().Add(time.Duration(cfg.PasswordResetTokenMinute) * time.Minute),
	}

	m.Db.GetDb().Create(&tokenData)

	encodedEmail := base64.RawURLEncoding.EncodeToString([]byte(*userData.Email))

	frontendUrl, err := url.Parse(cfg.AppBaseUrl)

	if err != nil {
		return nil
	}

	query := url.Values{}
	query.Set("email", encodedEmail)
	query.Set("token", tokenData.TokenHash)
	frontendUrl.Path = "/user-password/update-password"
	frontendUrl.RawQuery = query.Encode()

	meta := struct {
		Name        string
		Token       string
		ExpiresAt   string
		FrontendUrl string
	}{
		Name:        *userData.Name,
		Token:       token,
		ExpiresAt:   helpers.FormatViewDateTime(tokenData.ExpiresAt),
		FrontendUrl: frontendUrl.String(),
	}

	mailData := struct {
		AppName string
		Subject string
		Meta    interface{}
	}{
		AppName: cfg.AppName,
		Subject: "Request Password Reset",
		Meta:    meta,
	}

	m.Mailer.SendAsyncMail([]string{*userData.Email}, nil, nil, "user/passwordRestMail.html", mailData, nil, userData.Id, constants.SendPasswordRestMail, "users")
	return nil
}

func (m SendMailService) SendPasswordRestSuccessMail(userData *models.User, cfg config.Config) error {

	frontendUrl, err := url.Parse(cfg.AppBaseUrl)

	if err != nil {
		return nil
	}
	
	meta := struct {
		Name        string
		FrontendUrl string
	}{
		Name:        *userData.Name,
		FrontendUrl: frontendUrl.String(),
	}

	mailData := struct {
		AppName string
		Subject string
		Meta    interface{}
	}{
		AppName: cfg.AppName,
		Subject: "Successfully Password Reset",
		Meta:    meta,
	}

	m.Mailer.SendAsyncMail([]string{*userData.Email}, nil, nil, "user/passwordRestSuccessMail.html", mailData, nil, userData.Id, constants.SendPasswordRestSuccessMail, "users")

	return nil
}
