package seed

import (
	"os"
	"strconv"

	"admin-dashboard/cmd/common/helpers"
	"admin-dashboard/cmd/internal/models"

	"gorm.io/gorm"
)

func EmailSettingSeed(db *gorm.DB) error {

	var count int64

	db.Model(&models.EmailSetting{}).Count(&count)
	if count > 0 {
		return nil
	}

	var smtpPort *uint64

	if portStr := os.Getenv("MAIL_SMTP_PORT"); portStr != "" {
		if p, err := strconv.ParseUint(portStr, 10, 32); err == nil {
			smtpPort = &p
		}
	}

	driver := getEnvOrDefault("MAIL_DRIVER", "smtp")

	var fromEmail string
	switch driver {
	case "sendgrid":
		fromEmail = os.Getenv("MAIL_SENDGRID_FROM")
	case "ses":
		fromEmail = os.Getenv("MAIL_SES_FROM")
	default:
		fromEmail = os.Getenv("MAIL_SMTP_FROM")
	}

	data := models.EmailSetting{
		MainBaseModel:      BaseSeedData(),
		MailDriver:         driver,
		FromEmail:          helpers.Str(fromEmail),
		FromName:           helpers.Str(os.Getenv("MAIL_SENDER_NAME")),
		SendgridApiKey:     helpers.EncryptPtr(os.Getenv("MAIL_SENDGRID_API_KEY")),
		SmtpHost:           helpers.Str(os.Getenv("MAIL_SMTP_HOST")),
		SmtpPort:           smtpPort,
		SmtpUsername:       helpers.Str(os.Getenv("MAIL_SMTP_USERNAME")),
		SmtpPassword:       helpers.EncryptPtr(os.Getenv("MAIL_SMTP_PASSWORD")),
		SmtpEncryption:     helpers.Str(os.Getenv("MAIL_SMTP_ENCRYPTION")),
		AwsAccessKeyId:     helpers.EncryptPtr(os.Getenv("AWS_ACCESS_KEY_ID")),
		AwsSecretAccessKey: helpers.EncryptPtr(os.Getenv("AWS_SECRET_ACCESS_KEY")),
		AwsRegion:          helpers.Str(os.Getenv("MAIL_AWS_REGION")),
		AwsConfigSet:       helpers.Str(os.Getenv("MAIL_AWS_CONFIG_SET")),
	}

	if err := db.Create(&data).Error; err != nil {
		return err
	}

	history := models.EmailSettingHistory{
		EmailSetting: data,
		MainId:       data.Id,
	}

	return db.Create(&history).Error
}

func getEnvOrDefault(key, fallback string) string {

	if value := os.Getenv(key); value != "" {
		return value
	}

	return fallback
}
