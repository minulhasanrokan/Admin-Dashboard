package seed

import (
	"admin-dashboard/cmd/internal/models"

	"gorm.io/gorm"
)

func SmsSettingSeed(db *gorm.DB) error {

	var count int64

	db.Model(&models.SmsSetting{}).Count(&count)

	if count > 0 {
		return nil
	}

	data := models.SmsSetting{
		MainBaseModel:   BaseSeedData(),
		SmsDriver:       "sslcommerz",
		Status:          true,
		SslApiUrl:       "https://smsplus.sslwireless.com/api/v3/send-sms",
		SslApiToken:     "",
		SslSid:          "",
		SslCsmsIdPrefix: "ADM",
	}

	if err := db.Create(&data).Error; err != nil {
		return err
	}

	history := models.SmsSettingHistory{
		SmsSetting: data,
		MainId:     data.Id,
	}

	return db.Create(&history).Error
}
