package services

import (
	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/app/requests"
	"admin-dashboard/cmd/internal/models"
	"path/filepath"
	"time"
)

func (s *Service) SystemSetting() (*models.SystemSetting, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	setting, err := repository.SystemSetting()

	if err != nil {

		return nil, err
	}

	return &setting, nil
}

func (s *Service) BusinessSystemSettingUpdate(payload *requests.BusinessSettingRequest) (*models.SystemSetting, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	systemSetting, err := repository.SystemSetting()
	if err != nil {
		return nil, err
	}

	systemSetting.Name = &payload.Name
	systemSetting.Title = &payload.Title
	systemSetting.Slogan = &payload.Slogan
	systemSetting.Email = &payload.Email
	systemSetting.Mobile = &payload.Mobile
	systemSetting.Fax = &payload.Fax
	systemSetting.Address = &payload.Address
	systemSetting.Country = &payload.Country
	systemSetting.State = &payload.State
	systemSetting.City = &payload.City
	systemSetting.PostCode = &payload.PostCode
	systemSetting.Details = &payload.Details
	systemSetting.Version = &payload.Version
	systemSetting.CopyRight = &payload.CopyRight

	if err := repository.UpdateSystemSetting(&systemSetting); err != nil {
		return nil, err
	}

	return &systemSetting, nil
}

func (s *Service) ImagesSystemSettingUpdate(payload *requests.ImagesSettingRequest) (*models.SystemSetting, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	systemSetting, err := repository.SystemSetting()

	if err != nil {
		return nil, err
	}

	if payload.Logo != nil {

		logoName := "logo_" + time.Now().Format("20060102150405.000000000") + filepath.Ext(payload.Logo.Filename)

		_, err = s.Storage.UploadFile("system/"+logoName, payload.Logo)

		if err != nil {
			return nil, err
		}

		systemSetting.Logo = &logoName
	}

	if payload.LogoWhite != nil {

		logoWhiteName := "logo_white_" + time.Now().Format("20060102150405.000000000") + filepath.Ext(payload.LogoWhite.Filename)

		_, err = s.Storage.UploadFile("system/"+logoWhiteName, payload.LogoWhite)

		if err != nil {
			return nil, err
		}

		systemSetting.LogoWhite = &logoWhiteName
	}

	if payload.LogoSmall != nil {

		logoSmallName := "logo_small_" + time.Now().Format("20060102150405.000000000") + filepath.Ext(payload.LogoSmall.Filename)

		_, err = s.Storage.UploadFile("system/"+logoSmallName, payload.LogoSmall)

		if err != nil {
			return nil, err
		}

		systemSetting.LogoSmall = &logoSmallName
	}

	if payload.LogoDarkSmall != nil {

		logoDarkSmallName := "logo_dark_small_" + time.Now().Format("20060102150405.000000000") + filepath.Ext(payload.LogoDarkSmall.Filename)

		_, err = s.Storage.UploadFile("system/"+logoDarkSmallName, payload.LogoDarkSmall)

		if err != nil {
			return nil, err
		}

		systemSetting.LogoDarkSmall = &logoDarkSmallName
	}

	if payload.Icon != nil {

		iconName := "icon_" + time.Now().Format("20060102150405.000000000") + filepath.Ext(payload.Icon.Filename)

		_, err = s.Storage.UploadFile("system/"+iconName, payload.Icon)

		if err != nil {
			return nil, err
		}

		systemSetting.Icon = &iconName
	}

	if payload.AppleIcon != nil {

		appleIconName := "apple_icon_" + time.Now().Format("20060102150405.000000000") + filepath.Ext(payload.AppleIcon.Filename)

		_, err = s.Storage.UploadFile("system/"+appleIconName, payload.AppleIcon)

		if err != nil {
			return nil, err
		}

		systemSetting.AppleIcon = &appleIconName
	}

	if payload.LoginBg != nil {

		loginBgName := "login_bg_" + time.Now().Format("20060102150405.000000000") + filepath.Ext(payload.LoginBg.Filename)

		_, err = s.Storage.UploadFile("system/"+loginBgName, payload.LoginBg)

		if err != nil {
			return nil, err
		}

		systemSetting.LoginBg = &loginBgName
	}

	if payload.ForgotPasswordBg != nil {

		forgotPasswordBgName := "forgot_password_bg_" + time.Now().Format("20060102150405.000000000") + filepath.Ext(payload.ForgotPasswordBg.Filename)

		_, err = s.Storage.UploadFile("system/"+forgotPasswordBgName, payload.ForgotPasswordBg)

		if err != nil {
			return nil, err
		}

		systemSetting.ForgotPasswordBg = &forgotPasswordBgName
	}

	if payload.RegistrationBg != nil {

		registrationBgName := "registration_bg_" + time.Now().Format("20060102150405.000000000") + filepath.Ext(payload.RegistrationBg.Filename)

		_, err = s.Storage.UploadFile("system/"+registrationBgName, payload.RegistrationBg)

		if err != nil {
			return nil, err
		}

		systemSetting.RegistrationBg = &registrationBgName
	}

	if payload.VerificationBg != nil {

		verificationBgName := "verification_bg_" + time.Now().Format("20060102150405.000000000") + filepath.Ext(payload.VerificationBg.Filename)

		_, err = s.Storage.UploadFile("system/"+verificationBgName, payload.VerificationBg)

		if err != nil {
			return nil, err
		}

		systemSetting.VerificationBg = &verificationBgName
	}

	if err := repository.UpdateSystemSetting(&systemSetting); err != nil {
		return nil, err
	}

	return &systemSetting, nil
}
