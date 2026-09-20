package requests

type EmailSettingRequest struct {
	MailDriver string `param:"mail_driver" json:"mail_driver" form:"mail_driver" query:"mail_driver" validate:"required,oneof=smtp ses sendgrid"`
	FromName   string `param:"from_name" json:"from_name" form:"from_name" query:"from_name" validate:"required,min=2,max=200"`
	FromEmail  string `param:"from_email" json:"from_email" form:"from_email" query:"from_email" validate:"required,email"`
}

type SendGridSettingRequest struct {
	SendgridApiKey string `param:"sendgrid_api_key" json:"sendgrid_api_key" form:"sendgrid_api_key" query:"sendgrid_api_key" validate:"required"`
}

type SMTPSettingRequest struct {
	SmtpHost       string `param:"smtp_host" json:"smtp_host" form:"smtp_host" query:"smtp_host" validate:"required"`
	SmtpPort       uint64 `param:"smtp_port" json:"smtp_port" form:"smtp_port" query:"smtp_port" validate:"required,min=1,max=65535"`
	SmtpUsername   string `param:"smtp_username" json:"smtp_username" form:"smtp_username" query:"smtp_username" validate:"required"`
	SmtpPassword   string `param:"smtp_password" json:"smtp_password" form:"smtp_password" query:"smtp_password"`
	SmtpEncryption string `param:"smtp_encryption" json:"smtp_encryption" form:"smtp_encryption" query:"smtp_encryption" validate:"required,oneof=tls ssl none starttls"`
}

type SESSettingRequest struct {
	AwsAccessKeyId     string `param:"aws_access_key_id" json:"aws_access_key_id" form:"aws_access_key_id" query:"aws_access_key_id" validate:"required"`
	AwsSecretAccessKey string `param:"aws_secret_access_key" json:"aws_secret_access_key" form:"aws_secret_access_key" query:"aws_secret_access_key"`
	AwsRegion          string `param:"aws_region" json:"aws_region" form:"aws_region" query:"aws_region" validate:"required"`
	AwsConfigSet       string `param:"aws_config_set" json:"aws_config_set" form:"aws_config_set" query:"aws_config_set"`
}

type SendTestEmailRequest struct {
	RecipientEmail string `param:"recipient_email" json:"recipient_email" form:"recipient_email" query:"recipient_email" validate:"required,email"`
	Subject        string `param:"subject" json:"subject" form:"subject" query:"subject" validate:"required,max=200"`
	Message        string `param:"message" json:"message" form:"message" query:"message" validate:"required"`
}
