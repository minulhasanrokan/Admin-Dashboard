package mailer

import (
	"admin-dashboard/cmd/internal/config"
	"admin-dashboard/cmd/internal/databases"
	"bytes"
	"encoding/base64"
	"html/template"
	"mime"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
)

type SendGridMailer struct {
	ApiKey     string
	Sender     string
	SenderName string
	Logger     echo.Logger
	Config     config.Config
	Db         databases.DbConnection
}

func ConnectSendGrid(cfg *config.Config, logger echo.Logger, db databases.DbConnection) (MailConnection, error) {

	return &SendGridMailer{ApiKey: cfg.MailSendGridAPIKey, Sender: cfg.MailSendGridFrom, SenderName: cfg.MailSenderName, Logger: logger, Db: db}, nil
}

func (m *SendGridMailer) SendMail(to []string, cc []string, bcc []string, templateFile string, data EmailData, attachmentPath []string, tableId uint64, mailType string, tableName string) error {

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

	from := mail.NewEmail(m.SenderName, m.Sender)

	message := mail.NewV3Mail()
	message.SetFrom(from)
	message.Subject = subject.String()

	p := mail.NewPersonalization()

	for _, email := range to {
		p.AddTos(mail.NewEmail("", email))
	}

	for _, email := range cc {
		p.AddCCs(mail.NewEmail("", email))
	}

	for _, email := range bcc {
		p.AddBCCs(mail.NewEmail("", email))
	}

	message.AddPersonalizations(p)

	message.AddContent(mail.NewContent("text/html", body.String()))

	for _, file := range attachmentPath {

		data, err := os.ReadFile(file)

		if err != nil {
			return err
		}

		attachment := mail.NewAttachment()
		attachment.SetContent(base64.StdEncoding.EncodeToString(data))
		attachment.SetType(mime.TypeByExtension(filepath.Ext(file)))
		attachment.SetFilename(filepath.Base(file))
		attachment.SetDisposition("attachment")

		message.AddAttachment(attachment)
	}

	client := sendgrid.NewSendClient(m.ApiKey)

	_, err = client.Send(message)

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
func (m *SendGridMailer) SendAsyncMail(to []string, cc []string, bcc []string, templateFile string, data EmailData, attachmentPath []string, tableId uint64, mailType string, tableName string) {

	go func() {

		_ = m.SendMail(to, cc, bcc, templateFile, data, attachmentPath, tableId, mailType, tableName)
	}()
}
