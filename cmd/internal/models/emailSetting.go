package models

import "admin-dashboard/cmd/internal/constants"

type EmailSetting struct {
	MailDriver         string  `gorm:"column:mail_driver;type:varchar(20);default:'smtp'" json:"mail_driver"`
	FromEmail          *string `gorm:"column:from_email;type:varchar(100)" json:"from_email"`
	FromName           *string `gorm:"column:from_name;type:varchar(100)" json:"from_name"`
	SendgridApiKey     *string `gorm:"column:sendgrid_api_key;type:varchar(255)" json:"sendgrid_api_key,omitempty"`
	SmtpHost           *string `gorm:"column:smtp_host;type:varchar(255)" json:"smtp_host,omitempty"`
	SmtpPort           *uint64 `gorm:"column:smtp_port;type:int unsigned" json:"smtp_port,omitempty"`
	SmtpUsername       *string `gorm:"column:smtp_username;type:varchar(255)" json:"smtp_username,omitempty"`
	SmtpPassword       *string `gorm:"column:smtp_password;type:varchar(255)" json:"-"`
	SmtpEncryption     *string `gorm:"column:smtp_encryption;type:varchar(10);default:'tls'" json:"smtp_encryption,omitempty"`
	AwsAccessKeyId     *string `gorm:"column:aws_access_key_id;type:varchar(255)" json:"aws_access_key_id,omitempty"`
	AwsSecretAccessKey *string `gorm:"column:aws_secret_access_key;type:varchar(255)" json:"-"`
	AwsRegion          *string `gorm:"column:aws_region;type:varchar(50);default:'us-east-1'" json:"aws_region,omitempty"`
	AwsConfigSet       *string `gorm:"column:aws_config_set;type:varchar(100)" json:"aws_config_set,omitempty"`
	MainBaseModel
}

func (e EmailSetting) HasSendgrid() bool {
	return e.SendgridApiKey != nil && *e.SendgridApiKey != ""
}

func (e EmailSetting) HasSmtp() bool {
	return e.SmtpHost != nil && *e.SmtpHost != ""
}

func (e EmailSetting) HasSes() bool {
	return e.AwsAccessKeyId != nil && *e.AwsAccessKeyId != ""
}

func (e EmailSetting) GetFromEmail() string {
	if e.FromEmail != nil {
		return *e.FromEmail
	}
	return ""
}

func (e EmailSetting) GetFromName() string {
	if e.FromName != nil {
		return *e.FromName
	}
	return ""
}

func (e EmailSetting) GetSendgridApiKey() string {
	if e.SendgridApiKey != nil {
		return *e.SendgridApiKey
	}
	return ""
}

func (e EmailSetting) GetSmtpHost() string {
	if e.SmtpHost != nil {
		return *e.SmtpHost
	}
	return ""
}

func (e EmailSetting) GetSmtpPort() uint64 {
	if e.SmtpPort != nil && *e.SmtpPort > 0 {
		return *e.SmtpPort
	}
	return 587
}

func (e EmailSetting) GetSmtpUsername() string {
	if e.SmtpUsername != nil {
		return *e.SmtpUsername
	}
	return ""
}

func (e EmailSetting) GetSmtpPassword() string {
	if e.SmtpPassword != nil {
		return *e.SmtpPassword
	}
	return ""
}

func (e EmailSetting) GetSmtpEncryption() string {
	if e.SmtpEncryption != nil && *e.SmtpEncryption != "" {
		return *e.SmtpEncryption
	}
	return ""
}

func (e EmailSetting) GetAwsAccessKeyId() string {
	if e.AwsAccessKeyId != nil {
		return *e.AwsAccessKeyId
	}
	return ""
}

func (e EmailSetting) GetAwsSecretAccessKey() string {
	if e.AwsSecretAccessKey != nil {
		return *e.AwsSecretAccessKey
	}
	return ""
}

func (e EmailSetting) GetAwsRegion() string {
	if e.AwsRegion != nil && *e.AwsRegion != "" {
		return *e.AwsRegion
	}
	return ""
}

func (e EmailSetting) GetAwsConfigSet() string {
	if e.AwsConfigSet != nil {
		return *e.AwsConfigSet
	}
	return ""
}

func (e EmailSetting) GetAwsRegions() []constants.SelectOption {
	return constants.AwsRegions
}

func (e EmailSetting) GetSmtpEncryptionProtocols() []constants.SelectOption {
	return constants.SmtpEncryptionProtocols
}

type EmailSettingHistory struct {
	EmailSetting
	MainId uint64 `gorm:"column:main_id;type:int unsigned" json:"main_id"`
}
