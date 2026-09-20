package seed

import (
	"admin-dashboard/cmd/common/helpers"
	"admin-dashboard/cmd/internal/models"

	"gorm.io/gorm"
)

func SystemSettingSeed(db *gorm.DB) error {

	var count int64

	db.Model(&models.SystemSetting{}).Count(&count)

	if count > 0 {
		return nil
	}

	setting := models.SystemSetting{

		MainBaseModel:    BaseSeedData(),
		Name:             helpers.Str("Admin-Dashboard"),
		Title:            helpers.Str("Admin-Dashboard"),
		Slogan:           helpers.Str("Admin-Dashboard"),
		Email:            helpers.Str("minulhasanrokan@gmail.com"),
		Mobile:           helpers.Str("+8801627197089"),
		Fax:              helpers.Str("+8801627197089"),
		Address:          helpers.Str("Dhaka, Bangladesh"),
		Country:          helpers.Str("BD"),
		State:            helpers.Str("BD-C"),
		City:             helpers.Str("00821"),
		PostCode:         helpers.Str("1217"),
		Details:          helpers.Str("Default system description"),
		Icon:             helpers.Str("favicon.png"),
		AppleIcon:        helpers.Str("apple_icon.png"),
		Version:          helpers.Str("1.0.0"),
		CopyRight:        helpers.Str("© 2026 Example Co."),
		LoginBg:          helpers.Str("auth-bg.svg"),
		RegistrationBg:   helpers.Str("reg-bg.svg"),
		ForgotPasswordBg: helpers.Str("forgot-password-bg.svg"),
		Logo:             helpers.Str("logo.svg"),
		LogoWhite:        helpers.Str("logo-white.svg"),
		LogoSmall:        helpers.Str("logo-small.svg"),
		LogoDarkSmall:    helpers.Str("logo-small-white.svg"),
		VerificationBg:   helpers.Str("verification-bg.svg"),
	}

	if err := db.Create(&setting).Error; err != nil {
		return err
	}

	history := models.SystemSettingHistory{
		SystemSetting: setting,
		MainId:        setting.Id,
	}

	return db.Create(&history).Error
}
