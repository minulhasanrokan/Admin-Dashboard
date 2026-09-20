package services

import (
	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/app/requests"
	"admin-dashboard/cmd/common/smsService"
	"admin-dashboard/cmd/internal/constants"
	"admin-dashboard/cmd/internal/models"
	"errors"
	"fmt"
)

func (s *Service) SmsSetting() (*models.SmsSetting, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	setting, err := repository.SmsSetting()

	if err != nil {
		return nil, err
	}

	return &setting, nil
}

func (s *Service) SmsSettingUpdate(payload *requests.SmsSettingRequest) (*models.SmsSetting, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	setting, err := repository.SmsSetting()
	if err != nil {
		return nil, err
	}

	setting.SmsDriver = payload.SmsDriver
	setting.Status = payload.Status
	setting.SslApiUrl = payload.SslApiUrl
	setting.SslApiToken = payload.SslApiToken
	setting.SslSid = payload.SslSid
	setting.SslCsmsIdPrefix = payload.SslCsmsIdPrefix

	if err := repository.UpdateSmsSetting(&setting); err != nil {
		return nil, err
	}

	return &setting, nil
}

func (s *Service) SendSmsDirect(mobile string, message string) (*smsService.SslSmsResponse, error) {

	setting, err := s.SmsSetting()
	if err != nil {
		return nil, errors.New("SMS settings not configured: " + err.Error())
	}

	if !setting.Status {
		return nil, errors.New("SMS Gateway is currently disabled in SMS Settings")
	}

	return smsService.SendSslCommerzSms(
		setting.SslApiUrl,
		setting.SslApiToken,
		setting.SslSid,
		setting.SslCsmsIdPrefix,
		mobile,
		message,
	)
}

func (s *Service) SendTestSms(payload *requests.SendTestSmsRequest) (*smsService.SslSmsResponse, error) {

	setting, err := s.SmsSetting()
	driver := ""
	if err == nil {
		driver = setting.SmsDriver
	}

	resp, sendErr := s.SendSmsDirect(payload.Mobile, payload.Message)

	sentStatus := 1
	errorStatus := 0
	errorMsg := ""
	if sendErr != nil {
		sentStatus = 0
		errorStatus = 1
		errorMsg = sendErr.Error()
	}

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	_ = repository.CreateSmsHistory(&models.SmsHistory{
		ToMobile:    payload.Mobile,
		Message:     payload.Message,
		SmsDriver:   driver,
		SmsType:     constants.SendTestSms,
		SentStatus:  sentStatus,
		ErrorStatus: errorStatus,
		ErrorMsg:    errorMsg,
	})

	return resp, sendErr
}

func (s *Service) SendOtpSms(mobile string, otp string) error {

	setting, err := s.SmsSetting()
	driver := ""
	if err == nil {
		driver = setting.SmsDriver
	}

	message := fmt.Sprintf("Your OTP verification code is %s. Do not share this code with anyone.", otp)

	_, sendErr := s.SendSmsDirect(mobile, message)

	sentStatus := 1
	errorStatus := 0
	errorMsg := ""
	if sendErr != nil {
		sentStatus = 0
		errorStatus = 1
		errorMsg = sendErr.Error()
	}

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	_ = repository.CreateSmsHistory(&models.SmsHistory{
		ToMobile:    mobile,
		Message:     message,
		SmsDriver:   driver,
		SmsType:     constants.SendOtpSms,
		SentStatus:  sentStatus,
		ErrorStatus: errorStatus,
		ErrorMsg:    errorMsg,
	})

	return sendErr
}
