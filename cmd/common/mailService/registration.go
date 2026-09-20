package mailService

import (
	"admin-dashboard/cmd/common/helpers"
	"admin-dashboard/cmd/internal/config"
	"admin-dashboard/cmd/internal/constants"
	"admin-dashboard/cmd/internal/models"
	"crypto/sha256"
	"fmt"
	"net/url"
	"time"
)

func (m SendMailService) SendRegisterUserMail(registeredUser *models.User, cfg config.Config, expiryMinutes ...uint64) error {

	token, err := helpers.GenerateVerificationToken(cfg.VerificationTokenLength)

	if err != nil {
		return err
	}

	hash := sha256.Sum256([]byte(token))

	frontendUrl, err := url.Parse(cfg.AppBaseUrl)

	if err != nil {
		return nil
	}

	frontendUrl.Path = fmt.Sprintf("/user-registration/verify-account/%s", registeredUser.UUID)

	expireMin := uint64(cfg.VerificationTokenMinute)
	if len(expiryMinutes) > 0 && expiryMinutes[0] > 0 {
		expireMin = expiryMinutes[0]
	}

	verification := models.EmailVerificationToken{
		UserId:    registeredUser.Id,
		TokenHash: fmt.Sprintf("%x", hash),
		OtpToken:  token,
		ExpiresAt: time.Now().Add(time.Duration(expireMin) * time.Minute),
	}

	m.Db.GetDb().Create(&verification)

	meta := struct {
		Name        string
		Token       string
		ExpiresAt   string
		FrontendUrl string
	}{
		Name:        *registeredUser.Name,
		Token:       token,
		ExpiresAt:   helpers.FormatViewDateTime(verification.ExpiresAt),
		FrontendUrl: frontendUrl.String(),
	}

	mailData := struct {
		AppName string
		Subject string
		Meta    interface{}
	}{
		AppName: cfg.AppName,
		Subject: "Welcome to " + cfg.AppName + " – Verify your email",
		Meta:    meta,
	}

	m.Mailer.SendAsyncMail([]string{*registeredUser.Email}, nil, nil, "user/registerUserMail.html", mailData, nil, registeredUser.Id, constants.SendRegisterUserMail, "users")

	return nil
}

func (m SendMailService) SendVerifyUserMail(user *models.User, cfg config.Config) error {

	meta := struct {
		Name string
		Link string
	}{
		Name: *user.Name,
		Link: cfg.AppBaseUrl,
	}

	mailData := struct {
		AppName string
		Subject string
		Meta    interface{}
	}{
		AppName: cfg.AppName,
		Subject: "Thanks for verifying your email – " + cfg.AppName,
		Meta:    meta,
	}

	m.Mailer.SendAsyncMail([]string{*user.Email}, nil, nil, "user/verifyUserMail.html", mailData, nil, user.Id, constants.SendVerifyUserMail, "users")

	return nil
}
