package seed

import (
	"admin-dashboard/cmd/internal/models"

	"gorm.io/gorm"
)

func UserRoleSeed(db *gorm.DB) error {

	var count int64

	db.Model(&models.UserRole{}).Count(&count)

	if count > 0 {
		return nil
	}

	data := models.UserRole{

		MainBaseModel:      BaseSeedData(),
		PasswordMinLength:  8,
		PasswordMaxLength:  20,
		SessionTime:        2,
		SessionCheckTime:   1,
		PasswordExpireDays: 30,
	}

	if err := db.Create(&data).Error; err != nil {
		return err
	}

	history := models.UserRoleHistory{
		UserRole: data,
		MainId:   data.Id,
	}

	return db.Create(&history).Error
}
