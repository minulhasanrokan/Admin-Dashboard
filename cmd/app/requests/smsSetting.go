package requests

type SmsSettingRequest struct {
	SmsDriver       string `param:"sms_driver" json:"sms_driver" form:"sms_driver" query:"sms_driver" validate:"required"`
	Status          bool   `param:"status" json:"status" form:"status" query:"status"`
	SslApiUrl       string `param:"ssl_api_url" json:"ssl_api_url" form:"ssl_api_url" query:"ssl_api_url" validate:"required,url"`
	SslApiToken     string `param:"ssl_api_token" json:"ssl_api_token" form:"ssl_api_token" query:"ssl_api_token" validate:"required"`
	SslSid          string `param:"ssl_sid" json:"ssl_sid" form:"ssl_sid" query:"ssl_sid" validate:"required"`
	SslCsmsIdPrefix string `param:"ssl_csms_id_prefix" json:"ssl_csms_id_prefix" form:"ssl_csms_id_prefix" query:"ssl_csms_id_prefix" validate:"omitempty,max=20"`
}

type SendTestSmsRequest struct {
	Mobile  string `param:"mobile" json:"mobile" form:"mobile" query:"mobile" validate:"required"`
	Message string `param:"message" json:"message" form:"message" query:"message" validate:"required,min=1,max=500"`
}
