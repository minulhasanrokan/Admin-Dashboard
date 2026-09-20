package mailer

import (
	"admin-dashboard/cmd/internal/config"
	"admin-dashboard/cmd/internal/databases"
	"bytes"
	"context"
	"fmt"
	"html/template"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsConfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ses"
	"github.com/aws/aws-sdk-go-v2/service/ses/types"
	"github.com/labstack/echo/v4"
)

type SESMailer struct {
	client     *ses.Client
	Sender     string
	SenderName string
	Logger     echo.Logger
	Config     config.Config
	Db         databases.DbConnection
}

func ConnectSESMailer(cfg *config.Config, logger echo.Logger, db databases.DbConnection) (MailConnection, error) {

	ctx := context.TODO()

	awsCfg, err := awsConfig.LoadDefaultConfig(ctx, awsConfig.WithRegion(cfg.MailAWSRegion))

	if err != nil {
		return nil, err
	}

	client := ses.NewFromConfig(awsCfg)

	return &SESMailer{client: client, Sender: cfg.MailSESFrom, SenderName: cfg.MailSenderName, Logger: logger, Db: db}, nil
}

func (m *SESMailer) SendMail(to []string, cc []string, bcc []string, templateFile string, data EmailData, attachmentPath []string, tableId uint64, mailType string, tableName string) error {
	_ = attachmentPath
	absolutePath := "templates/" + templateFile

	tmpl, err := template.ParseFS(templateFs, absolutePath)

	if err != nil {

		return err
	}

	subject := new(bytes.Buffer)

	if err := tmpl.ExecuteTemplate(subject, "subject", data); err != nil {

		return err
	}

	body := new(bytes.Buffer)

	if err := tmpl.ExecuteTemplate(body, "htmlBody", data); err != nil {

		return err
	}

	input := &ses.SendEmailInput{

		Destination: &types.Destination{
			ToAddresses:  to,
			CcAddresses:  cc,
			BccAddresses: bcc,
		},

		Message: &types.Message{

			Subject: &types.Content{

				Data: aws.String(subject.String()),
			},
			Body: &types.Body{

				Html: &types.Content{

					Data: aws.String(body.String()),
				},
				Text: &types.Content{

					Data: aws.String("Your client does not support HTML"),
				},
			},
		},
		Source: aws.String(fmt.Sprintf("%s <%s>", m.SenderName, m.Sender)),
	}
	_, err = m.client.SendEmail(context.TODO(), input)

	sentStatus := 1
	errorStatus := 1
	ErrorMsg := ""

	if err != nil {

		sentStatus = 0
		errorStatus = 1
		ErrorMsg = err.Error()
	}

	logger := NewMailServiceLogger(m.Db)

	now := time.Now()
	log := MailLog{
		FromAddress:    m.Config.MailSMTPFrom,
		ToAddress:      strings.Join(to, ","),
		BccAddress:     strings.Join(bcc, ","),
		AttachmentPath: strings.Join(attachmentPath, ","),
		CcAddress:      strings.Join(cc, ","),
		TableName:      tableName,
		TableId:        tableId,
		Subject:        subject.String(),
		Body:           body.String(),
		SentStatus:     sentStatus,
		TryStatus:      1,
		ErrorStatus:    errorStatus,
		ErrorMsg:       ErrorMsg,
		MailType:       mailType,
		CreatedAt:      &now,
	}

	_ = logger.InsertMailLog(log)

	return err
}

func (m *SESMailer) SendAsyncMail(to []string, cc []string, bcc []string, templateFile string, data EmailData, attachmentPath []string, tableId uint64, mailType string, tableName string) {

	go func() {

		_ = m.SendMail(to, cc, bcc, templateFile, data, attachmentPath, tableId, mailType, tableName)
	}()
}
