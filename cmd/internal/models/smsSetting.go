package models

type SmsSetting struct {
	SmsDriver       string `gorm:"column:sms_driver;type:varchar(50);default:'sslcommerz'" json:"sms_driver"`
	Status          bool   `gorm:"column:status;default:true" json:"status"`
	SslApiUrl       string `gorm:"column:ssl_api_url;type:varchar(255);default:'https://smsplus.sslwireless.com/api/v3/send-sms'" json:"ssl_api_url"`
	SslApiToken     string `gorm:"column:ssl_api_token;type:varchar(255)" json:"ssl_api_token"`
	SslSid          string `gorm:"column:ssl_sid;type:varchar(100)" json:"ssl_sid"`
	SslCsmsIdPrefix string `gorm:"column:ssl_csms_id_prefix;type:varchar(50);default:'ADM'" json:"ssl_csms_id_prefix"`
	MainBaseModel
}

type SmsSettingHistory struct {
	SmsSetting
	MainId uint64 `gorm:"column:main_id;type:int unsigned" json:"main_id"`
}
