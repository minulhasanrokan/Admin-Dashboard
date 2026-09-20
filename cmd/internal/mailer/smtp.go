package mailer

import (
	"admin-dashboard/cmd/internal/config"
	"admin-dashboard/cmd/internal/databases"
	"bytes"
	"html/template"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"gopkg.in/gomail.v2"
)

type SMTPMailer struct {
	Dialer     *gomail.Dialer
	Sender     string
	SenderName string
	Logger     echo.Logger
	Config     config.Config
	Db         databases.DbConnection
}

func ConnectSMTPMailer(cfg *config.Config, logger echo.Logger, db databases.DbConnection) (MailConnection, error) {

	dialer := gomail.NewDialer(cfg.MailSMTPHost, cfg.MailSMTPPort, cfg.MailSMTPUsername, cfg.MailSMTPPassword)

	return &SMTPMailer{Dialer: dialer, Sender: cfg.MailSMTPFrom, SenderName: cfg.MailSenderName, Logger: logger, Config: *cfg, Db: db}, nil
}

func (m *SMTPMailer) SendMail(to []string, cc []string, bcc []string, templateFile string, data EmailData, attachmentPath []string, tableId uint64, mailType string, tableName string) error {

	absolutePath := "templates/" + templateFile

	tmpl, err := template.ParseFS(templateFs, absolutePath)

	if err != nil {
		return err
	}

	subject := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(subject, "subject", data)

	if err != nil {
		return err
	}

	body := new(bytes.Buffer)

	err = tmpl.ExecuteTemplate(body, "htmlBody", data)

	if err != nil {
		return err
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", msg.FormatAddress(m.Sender, m.SenderName))
	msg.SetHeader("To", to...)
	if len(cc) > 0 {
		msg.SetHeader("Cc", cc...)
	}
	if len(bcc) > 0 {
		msg.SetHeader("Bcc", bcc...)
	}

	for _, file := range attachmentPath {

		if file == "" {
			continue
		}
		msg.Attach(file)
	}

	msg.SetHeader("Subject", subject.String())
	msg.SetBody("text/html", body.String())

	err = m.Dialer.DialAndSend(msg)

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

func (m *SMTPMailer) SendAsyncMail(to []string, cc []string, bcc []string, templateFile string, data EmailData, attachmentPath []string, tableId uint64, mailType string, tableName string) {

	go func() {

		_ = m.SendMail(to, cc, bcc, templateFile, data, attachmentPath, tableId, mailType, tableName)
	}()
}
