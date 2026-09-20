package seed

import (
	"admin-dashboard/cmd/internal/models"

	"gorm.io/gorm"
)

func AuthenticationSettingSeed(db *gorm.DB) error {

	var count int64

	db.Model(&models.AuthenticationSetting{}).Count(&count)

	if count > 0 {
		return nil
	}

	data := models.AuthenticationSetting{
		MainBaseModel:        BaseSeedData(),
		AllowRegistration:    true,
		InviteOnly:           false,
		VerificationRequired: true,
		VerificationExpired:  15,
		ReferralSystem:       false,
		LoginType:            "email",
		PasswordRequired:     true,
		OtpSystem:            false,
		OtpType:              "email",
	}

	if err := db.Create(&data).Error; err != nil {
		return err
	}

	history := models.AuthenticationSettingHistory{
		AuthenticationSetting: data,
		MainId:                data.Id,
	}

	return db.Create(&history).Error
}
