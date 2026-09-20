package services

import (
	"context"
	"errors"
	"fmt"

	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/app/requests"
	"admin-dashboard/cmd/common/helpers"
	"admin-dashboard/cmd/internal/constants"
	"admin-dashboard/cmd/internal/mailer"
	"admin-dashboard/cmd/internal/models"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"gopkg.in/gomail.v2"
)

func (s *Service) EmailSetting() (*models.EmailSetting, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	setting, err := repository.EmailSetting()
	if err != nil {
		return nil, err
	}

	setting.SendgridApiKey = helpers.DecryptPtr(setting.SendgridApiKey)
	setting.SmtpPassword = helpers.DecryptPtr(setting.SmtpPassword)
	setting.AwsAccessKeyId = helpers.DecryptPtr(setting.AwsAccessKeyId)
	setting.AwsSecretAccessKey = helpers.DecryptPtr(setting.AwsSecretAccessKey)

	return &setting, nil
}

func (s *Service) EmailSettingUpdate(payload *requests.EmailSettingRequest) (*models.EmailSetting, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	emailSetting, err := repository.EmailSetting()
	if err != nil {
		return nil, err
	}

	emailSetting.FromName = &payload.FromName
	emailSetting.FromEmail = &payload.FromEmail
	emailSetting.MailDriver = payload.MailDriver

	if err := repository.UpdateEmailSetting(&emailSetting); err != nil {
		return nil, err
	}

	return &emailSetting, nil
}

func (s *Service) SendGridSettingUpdate(payload *requests.SendGridSettingRequest) (*models.EmailSetting, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	emailSetting, err := repository.EmailSetting()
	if err != nil {
		return nil, err
	}

	emailSetting.SendgridApiKey = helpers.EncryptPtr(payload.SendgridApiKey)

	if err := repository.UpdateEmailSetting(&emailSetting); err != nil {
		return nil, err
	}

	return &emailSetting, nil
}

func (s *Service) SMTPSettingUpdate(payload *requests.SMTPSettingRequest) (*models.EmailSetting, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	emailSetting, err := repository.EmailSetting()
	if err != nil {
		return nil, err
	}

	emailSetting.SmtpHost = &payload.SmtpHost
	emailSetting.SmtpPort = &payload.SmtpPort
	emailSetting.SmtpUsername = &payload.SmtpUsername
	if payload.SmtpPassword != "" {
		emailSetting.SmtpPassword = helpers.EncryptPtr(payload.SmtpPassword)
	}
	emailSetting.SmtpEncryption = &payload.SmtpEncryption

	if err := repository.UpdateEmailSetting(&emailSetting); err != nil {
		return nil, err
	}

	return &emailSetting, nil
}

func (s *Service) SESSettingUpdate(payload *requests.SESSettingRequest) (*models.EmailSetting, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	emailSetting, err := repository.EmailSetting()
	if err != nil {
		return nil, err
	}

	emailSetting.AwsAccessKeyId = helpers.EncryptPtr(payload.AwsAccessKeyId)
	if payload.AwsSecretAccessKey != "" {
		emailSetting.AwsSecretAccessKey = helpers.EncryptPtr(payload.AwsSecretAccessKey)
	}
	emailSetting.AwsRegion = &payload.AwsRegion
	emailSetting.AwsConfigSet = &payload.AwsConfigSet

	if err := repository.UpdateEmailSetting(&emailSetting); err != nil {
		return nil, err
	}

	return &emailSetting, nil
}

func (s *Service) SendEmailDirect(to string, subject string, bodyHTML string) error {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	emailSetting, err := repository.EmailSetting()
	if err != nil {
		return err
	}

	fromEmail := "noreply@example.com"
	if s.Config.MailSMTPFrom != "" {
		fromEmail = s.Config.MailSMTPFrom
	} else if s.Config.MailSendGridFrom != "" {
		fromEmail = s.Config.MailSendGridFrom
	} else if s.Config.MailSESFrom != "" {
		fromEmail = s.Config.MailSESFrom
	}
	if emailSetting.FromEmail != nil && *emailSetting.FromEmail != "" {
		fromEmail = *emailSetting.FromEmail
	}

	fromName := "Admin Dashboard"
	if s.Config.AppName != "" {
		fromName = s.Config.AppName
	} else if s.Config.MailSenderName != "" {
		fromName = s.Config.MailSenderName
	}
	if emailSetting.FromName != nil && *emailSetting.FromName != "" {
		fromName = *emailSetting.FromName
	}

	switch emailSetting.MailDriver {
	case "sendgrid":
		apiKey := ""
		if emailSetting.SendgridApiKey != nil {
			if dec := helpers.DecryptPtr(emailSetting.SendgridApiKey); dec != nil {
				apiKey = *dec
			}
		}
		if apiKey == "" {
			return errors.New("SendGrid API Key is not configured. Please open SendGrid Settings to enter your API key.")
		}

		from := mail.NewEmail(fromName, fromEmail)
		recipient := mail.NewEmail(to, to)
		content := mail.NewContent("text/html", bodyHTML)
		m := mail.NewV3MailInit(from, subject, recipient, content)

		client := sendgrid.NewSendClient(apiKey)
		resp, err := client.Send(m)
		if err != nil {
			return fmt.Errorf("SendGrid error: %w", err)
		} else if resp.StatusCode >= 400 {
			return fmt.Errorf("SendGrid error (status %d): %s", resp.StatusCode, resp.Body)
		}
		return nil

	case "ses":
		accessKey := ""
		if emailSetting.AwsAccessKeyId != nil {
			if dec := helpers.DecryptPtr(emailSetting.AwsAccessKeyId); dec != nil {
				accessKey = *dec
			}
		}
		secretKey := ""
		if emailSetting.AwsSecretAccessKey != nil {
			if dec := helpers.DecryptPtr(emailSetting.AwsSecretAccessKey); dec != nil {
				secretKey = *dec
			}
		}
		region := "us-east-1"
		if emailSetting.AwsRegion != nil && *emailSetting.AwsRegion != "" {
			region = *emailSetting.AwsRegion
		}

		if accessKey == "" || secretKey == "" {
			return errors.New("AWS SES Access Key or Secret Key is not configured. Please open AWS SES Settings to enter your credentials.")
		}

		creds := credentials.NewStaticCredentialsProvider(accessKey, secretKey, "")
		awsCfg, err := awsConfig.LoadDefaultConfig(
			context.Background(),
			awsConfig.WithRegion(region),
			awsConfig.WithCredentialsProvider(creds),
		)
		if err != nil {
			return fmt.Errorf("AWS config error: %w", err)
		}

		client := ses.NewFromConfig(awsCfg)
		input := &ses.SendEmailInput{
			Destination: &types.Destination{
				ToAddresses: []string{to},
			},
			Message: &types.Message{
				Body: &types.Body{
					Html: &types.Content{
						Data:    aws.String(bodyHTML),
						Charset: aws.String("UTF-8"),
					},
				},
				Subject: &types.Content{
					Data:    aws.String(subject),
					Charset: aws.String("UTF-8"),
				},
			},
			Source: aws.String(fmt.Sprintf("%s <%s>", fromName, fromEmail)),
		}

		if emailSetting.AwsConfigSet != nil && *emailSetting.AwsConfigSet != "" {
			input.ConfigurationSetName = emailSetting.AwsConfigSet
		}

		_, err = client.SendEmail(context.Background(), input)
		if err != nil {
			return fmt.Errorf("AWS SES error: %w", err)
		}
		return nil

	default:
		host := ""
		if emailSetting.SmtpHost != nil {
			host = *emailSetting.SmtpHost
		}
		port := 587
		if emailSetting.SmtpPort != nil && *emailSetting.SmtpPort > 0 {
			port = int(*emailSetting.SmtpPort)
		}
		username := ""
		if emailSetting.SmtpUsername != nil {
			username = *emailSetting.SmtpUsername
		}
		password := ""
		if emailSetting.SmtpPassword != nil {
			if dec := helpers.DecryptPtr(emailSetting.SmtpPassword); dec != nil {
				password = *dec
			}
		}

		if host == "" {
			return errors.New("SMTP Host is not configured. Please open SMTP Settings to enter your mail server host.")
		}

		msg := gomail.NewMessage()
		msg.SetHeader("From", msg.FormatAddress(fromEmail, fromName))
		msg.SetHeader("To", to)
		msg.SetHeader("Subject", subject)
		msg.SetBody("text/html", bodyHTML)

		dialer := gomail.NewDialer(host, port, username, password)
		if emailSetting.SmtpEncryption != nil && *emailSetting.SmtpEncryption == "ssl" {
			dialer.SSL = true
		}

		if err := dialer.DialAndSend(msg); err != nil {
			return fmt.Errorf("SMTP error: %w", err)
		}
		return nil
	}
}

func (s *Service) SendTestEmail(payload *requests.SendTestEmailRequest) error {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	emailSetting, err := repository.EmailSetting()
	if err != nil {
		return err
	}

	fromEmail := "noreply@example.com"
	if s.Config.MailSMTPFrom != "" {
		fromEmail = s.Config.MailSMTPFrom
	} else if s.Config.MailSendGridFrom != "" {
		fromEmail = s.Config.MailSendGridFrom
	} else if s.Config.MailSESFrom != "" {
		fromEmail = s.Config.MailSESFrom
	}
	if emailSetting.FromEmail != nil && *emailSetting.FromEmail != "" {
		fromEmail = *emailSetting.FromEmail
	}

	fromName := "Admin Dashboard"
	if s.Config.AppName != "" {
		fromName = s.Config.AppName
	} else if s.Config.MailSenderName != "" {
		fromName = s.Config.MailSenderName
	}
	if emailSetting.FromName != nil && *emailSetting.FromName != "" {
		fromName = *emailSetting.FromName
	}

	bodyHTML := fmt.Sprintf(
		`<div style="font-family: Arial, sans-serif; max-width: 600px; margin: 0 auto; padding: 20px; border: 1px solid #e0e0e0; rounded: 8px;">
			<h2 style="color: #333; margin-bottom: 10px;">Test Email from %s</h2>
			<p style="color: #666; font-size: 15px; line-height: 1.5;">%s</p>
			<hr style="border: 0; border-top: 1px solid #eee; margin: 20px 0;">
			<p style="color: #999; font-size: 12px;">Active Driver: <strong>%s</strong> | Timestamp: %s</p>
		</div>`,
		fromName,
		payload.Message,
		emailSetting.MailDriver,
		helpers.NowView(),
	)

	sendErr := s.SendEmailDirect(payload.RecipientEmail, payload.Subject, bodyHTML)

	sentStatus := 1
	errorStatus := 0
	errorMsg := ""
	if sendErr != nil {
		sentStatus = 0
		errorStatus = 1
		errorMsg = sendErr.Error()
	}

	mailLogger := mailer.NewMailServiceLogger(s.Db)
	_ = mailLogger.InsertMailLog(mailer.MailLog{
		FromAddress: fromEmail,
		ToAddress:   payload.RecipientEmail,
		Subject:     payload.Subject,
		Body:        bodyHTML,
		SentStatus:  sentStatus,
		TryStatus:   1,
		ErrorStatus: errorStatus,
		ErrorMsg:    errorMsg,
		MailType:    constants.SendTestEmailMail,
		TableName:   "email_settings",
		TableId:     emailSetting.Id,
	})

	return sendErr
}

