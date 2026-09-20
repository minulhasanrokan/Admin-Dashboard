package mailer

import (
	"admin-dashboard/cmd/internal/config"
	"admin-dashboard/cmd/internal/databases"
	"embed"
	"errors"

	"github.com/labstack/echo/v4"
)

//go:embed templates/*
var templateFs embed.FS

type EmailData struct {
	AppName string
	Subject string
	Meta    interface{}
}

type MailConnection interface {
	SendMail(to []string, cc []string, bcc []string, templateFile string, data EmailData, attachmentPath []string, tableId uint64, mailType string, tableName string) error
	SendAsyncMail(to []string, cc []string, bcc []string, templateFile string, data EmailData, attachmentPath []string, tableId uint64, mailType string, tableName string)
}

func ConnectMailer(cfg *config.Config, logger echo.Logger, db databases.DbConnection) (MailConnection, error) {

	activeMailer := cfg.MailDriver

	if activeMailer == "" {

		return nil, errors.New("MAILER_DRIVER environment variable not set")
	}

	var mailer MailConnection
	var err error

	switch activeMailer {
	case "smtp":
		mailer, err = ConnectSMTPMailer(cfg, logger, db)
	case "sendgrid":
		mailer, err = ConnectSendGrid(cfg, logger, db)
	case "ses":
		mailer, err = ConnectSESMailer(cfg, logger, db)
	default:
		err = errors.New(activeMailer + " is not supported")
	}

	if err != nil {

		logger.Errorf("ConnectMailer: %v", err.Error())

		return nil, err
	}

	return mailer, nil
}
