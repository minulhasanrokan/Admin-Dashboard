package main

import (
	"admin-dashboard/cmd/internal/config"
	"admin-dashboard/cmd/internal/databases"
	"admin-dashboard/cmd/internal/databases/seed"
	"admin-dashboard/cmd/internal/models"
	"log"
)

func main() {

	cfg, err := config.Init()

	if err != nil {

		panic(err)
	}

	dbConnPtr, err := databases.ConnectDatabase(cfg)

	if err != nil {
		panic(err)
	}

	dbConn := dbConnPtr

	db := dbConn.GetDb()

	err = db.AutoMigrate(&models.SystemSetting{}, &models.SystemSettingHistory{}, &models.PrivacyPolicy{}, &models.PrivacyPolicyHistory{}, &models.UserRole{}, &models.UserRoleHistory{}, &models.User{}, &models.UserHistory{}, &models.EmailVerificationToken{}, &models.PasswordResetToken{}, &models.UserLoginHistory{}, &models.RightGroup{}, &models.RightCategory{}, &models.RightDetail{}, &models.UserRight{}, &models.UserRightHistory{}, &models.MailHistory{}, &models.SmsHistory{}, &models.Country{}, &models.CountryHistory{}, &models.Region{}, &models.RegionHistory{}, &models.SubRegion{}, &models.SubRegionHistory{}, &models.State{}, &models.StateHistory{}, &models.City{}, &models.CityHistory{}, &models.UserActivity{}, models.ActivityList{}, &models.EmailSetting{}, &models.EmailSettingHistory{}, &models.BackupType{}, &models.BackupHistory{}, &models.AuthenticationSetting{}, &models.AuthenticationSettingHistory{}, &models.SmsSetting{}, &models.SmsSettingHistory{}, &models.Session{})

	if err != nil {

		panic(err)
	}

	err = seed.AllSeed(db)

	if err != nil {

		panic(err)
	}

	log.Println("Database migrated successfully")
}
