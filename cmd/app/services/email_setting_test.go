package services

import (
	"admin-dashboard/cmd/internal/models"
	"testing"
)

func TestEmailSettingModelGetters(t *testing.T) {
	// 1. Empty EmailSetting
	empty := models.EmailSetting{}
	if enc := empty.GetSmtpEncryption(); enc != "" {
		t.Errorf("expected empty string for unset SmtpEncryption, got '%s'", enc)
	}
	if reg := empty.GetAwsRegion(); reg != "" {
		t.Errorf("expected empty string for unset AwsRegion, got '%s'", reg)
	}
	if empty.GetSmtpHost() != "" {
		t.Errorf("expected empty string for unset SmtpHost")
	}
	if empty.GetSmtpPort() != 587 {
		t.Errorf("expected default port 587 for unset SmtpPort, got %d", empty.GetSmtpPort())
	}
	if empty.HasSendgrid() {
		t.Errorf("expected HasSendgrid false for empty setting")
	}
	if empty.HasSmtp() {
		t.Errorf("expected HasSmtp false for empty setting")
	}
	if empty.HasSes() {
		t.Errorf("expected HasSes false for empty setting")
	}

	// 2. Populated EmailSetting
	enc := "ssl"
	reg := "eu-west-1"
	host := "smtp.example.com"
	port := uint64(465)
	user := "myuser"
	pass := "mypass"
	sgKey := "SG.123456"
	awsKey := "AKIA12345"
	awsSec := "SECRET123"
	fromEmail := "info@example.com"
	fromName := "Support"

	populated := models.EmailSetting{
		MailDriver:         "smtp",
		FromEmail:          &fromEmail,
		FromName:           &fromName,
		SmtpHost:           &host,
		SmtpPort:           &port,
		SmtpUsername:       &user,
		SmtpPassword:       &pass,
		SmtpEncryption:     &enc,
		SendgridApiKey:     &sgKey,
		AwsAccessKeyId:     &awsKey,
		AwsSecretAccessKey: &awsSec,
		AwsRegion:          &reg,
	}

	if populated.GetSmtpEncryption() != "ssl" {
		t.Errorf("expected 'ssl', got '%s'", populated.GetSmtpEncryption())
	}
	if populated.GetAwsRegion() != "eu-west-1" {
		t.Errorf("expected 'eu-west-1', got '%s'", populated.GetAwsRegion())
	}
	if populated.GetSmtpHost() != "smtp.example.com" {
		t.Errorf("expected 'smtp.example.com', got '%s'", populated.GetSmtpHost())
	}
	if populated.GetSmtpPort() != 465 {
		t.Errorf("expected 465, got %d", populated.GetSmtpPort())
	}
	if !populated.HasSendgrid() {
		t.Errorf("expected HasSendgrid true")
	}
	if !populated.HasSmtp() {
		t.Errorf("expected HasSmtp true")
	}
	if !populated.HasSes() {
		t.Errorf("expected HasSes true")
	}
}

func TestEmailSettingProtocolsAndRegions(t *testing.T) {
	setting := models.EmailSetting{}
	protocols := setting.GetSmtpEncryptionProtocols()
	if len(protocols) < 3 {
		t.Errorf("expected at least 3 encryption protocols (tls, ssl, starttls), got %d", len(protocols))
	}

	regions := setting.GetAwsRegions()
	if len(regions) < 5 {
		t.Errorf("expected at least 5 AWS regions, got %d", len(regions))
	}
}
