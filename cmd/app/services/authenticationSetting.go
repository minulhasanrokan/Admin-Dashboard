package services

import (
	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/app/requests"
	"admin-dashboard/cmd/internal/models"
)

func (s *Service) AuthenticationSetting() (*models.AuthenticationSetting, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	setting, err := repository.AuthenticationSetting()

	if err != nil {
		return nil, err
	}

	return &setting, nil
}

func (s *Service) AuthenticationSettingUpdate(payload *requests.AuthenticationSettingRequest) (*models.AuthenticationSetting, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	setting, err := repository.AuthenticationSetting()
	if err != nil {
		return nil, err
	}

	setting.AllowRegistration = payload.AllowRegistration
	setting.InviteOnly = payload.InviteOnly
	setting.VerificationRequired = payload.VerificationRequired
	setting.VerificationExpired = payload.VerificationExpired
	setting.ReferralSystem = payload.ReferralSystem
	setting.LoginType = payload.LoginType
	setting.PasswordRequired = payload.PasswordRequired
	setting.OtpSystem = payload.OtpSystem
	setting.OtpType = payload.OtpType

	if err := repository.UpdateAuthenticationSetting(&setting); err != nil {
		return nil, err
	}

	return &setting, nil
}
