package renderer

import (
	"bytes"
	"os"
	"testing"

	"admin-dashboard/cmd/internal/models"
)

func TestEmailSettingTemplateRender(t *testing.T) {
	root := "../../../web/html"
	if _, err := os.Stat(root); err != nil {
		root = "web/html"
	}
	tmpl := loadTemplates(root)

	emptySetting := &models.EmailSetting{}
	var buf1 bytes.Buffer
	data1 := map[string]interface{}{
		"EmailSetting": emptySetting,
		"ConfigData": map[string]string{
			"AppBaseUrl":   "http://localhost:3000",
			"AssetBaseUrl": "http://localhost:3000",
			"AppName":      "Admin-Dashboard",
		},
	}
	if err := tmpl.ExecuteTemplate(&buf1, "settings/email_setting.gohtml", data1); err != nil {
		t.Fatalf("failed to render empty EmailSetting template: %v", err)
	}

	host := "smtp.example.com"
	port := uint64(587)
	user := "testuser"
	pass := "secret"
	key := "SG.12345"
	awsKey := "AKIA123"
	awsSec := "SEC123"
	fromEmail := "admin@example.com"
	fromName := "Admin"
	populatedSetting := &models.EmailSetting{
		MailDriver:         "sendgrid",
		FromEmail:          &fromEmail,
		FromName:           &fromName,
		SmtpHost:           &host,
		SmtpPort:           &port,
		SmtpUsername:       &user,
		SmtpPassword:       &pass,
		SendgridApiKey:     &key,
		AwsAccessKeyId:     &awsKey,
		AwsSecretAccessKey: &awsSec,
	}
	var buf2 bytes.Buffer
	data2 := map[string]interface{}{
		"EmailSetting": populatedSetting,
		"ConfigData": map[string]string{
			"AppBaseUrl":   "http://localhost:3000",
			"AssetBaseUrl": "http://localhost:3000",
			"AppName":      "Admin-Dashboard",
		},
	}
	if err := tmpl.ExecuteTemplate(&buf2, "settings/email_setting.gohtml", data2); err != nil {
		t.Fatalf("failed to render populated EmailSetting template: %v", err)
	}
}
