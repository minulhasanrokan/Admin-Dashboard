package requests

import (
	"testing"

	"github.com/go-playground/validator/v10"
)

func TestEmailSettingRequestsValidation(t *testing.T) {
	validate := validator.New()

	validSendgrid := SendGridSettingRequest{
		SendgridApiKey: "SG.valid_key_12345",
	}
	if err := validate.Struct(validSendgrid); err != nil {
		t.Errorf("expected valid SendGrid request, got error: %v", err)
	}

	invalidSendgrid := SendGridSettingRequest{
		SendgridApiKey: "",
	}
	if err := validate.Struct(invalidSendgrid); err == nil {
		t.Errorf("expected error for empty SendGrid API key, got nil")
	}

	validSMTP := SMTPSettingRequest{
		SmtpHost:       "smtp.mailtrap.io",
		SmtpPort:       587,
		SmtpUsername:   "user123",
		SmtpPassword:   "secretpass",
		SmtpEncryption: "tls",
	}
	if err := validate.Struct(validSMTP); err != nil {
		t.Errorf("expected valid SMTP request, got error: %v", err)
	}

	invalidSMTPPort := SMTPSettingRequest{
		SmtpHost:       "smtp.mailtrap.io",
		SmtpPort:       70000,
		SmtpUsername:   "user123",
		SmtpEncryption: "tls",
	}
	if err := validate.Struct(invalidSMTPPort); err == nil {
		t.Errorf("expected error for invalid SMTP port, got nil")
	}

	validSES := SESSettingRequest{
		AwsAccessKeyId:     "AKIAIOSFODNN7EXAMPLE",
		AwsSecretAccessKey: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY",
		AwsRegion:          "us-east-1",
		AwsConfigSet:       "config-set-1",
	}
	if err := validate.Struct(validSES); err != nil {
		t.Errorf("expected valid SES request, got error: %v", err)
	}

	invalidSES := SESSettingRequest{
		AwsAccessKeyId: "",
		AwsRegion:      "",
	}
	if err := validate.Struct(invalidSES); err == nil {
		t.Errorf("expected error for missing SES access key and region, got nil")
	}

	validTestEmail := SendTestEmailRequest{
		RecipientEmail: "test@example.com",
		Subject:        "Test Subject",
		Message:        "Test email message content",
	}
	if err := validate.Struct(validTestEmail); err != nil {
		t.Errorf("expected valid test email request, got error: %v", err)
	}

	invalidTestEmail := SendTestEmailRequest{
		RecipientEmail: "not-an-email",
		Subject:        "",
		Message:        "",
	}
	if err := validate.Struct(invalidTestEmail); err == nil {
		t.Errorf("expected error for invalid recipient email and empty subject/message, got nil")
	}
}
