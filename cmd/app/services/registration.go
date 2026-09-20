package services

import (
	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/app/requests"
	"admin-dashboard/cmd/common/helpers"
	"admin-dashboard/cmd/internal/constants"
	"admin-dashboard/cmd/internal/models"
	"time"
)

func (s *Service) RegisterUser(payload *requests.UserRegisterRequest, signUpIp string) (*models.User, error) {

	hashedPassword, err := helpers.HashPassword(payload.Password)

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	userName, err := repository.GenerateUniqueUserName(payload.Name)

	if err != nil {
		return nil, err
	}

	userRoe, err := repository.UserRole()

	if err != nil {
		return nil, err
	}

	userRole, err := repository.UserRole()

	if err != nil {

		return nil, err
	}

	now := time.Now()

	passwordExpireDate := now.AddDate(0, 0, int(userRole.PasswordExpireDays))

	createdUser := models.User{
		Name:                   &payload.Name,
		Email:                  &payload.Email,
		Mobile:                 &payload.Mobile,
		Password:               &hashedPassword,
		PasswordChangeStatus:   helpers.UInt(1),
		PasswordChangeDate:     &now,
		PasswordExpireDate:     &passwordExpireDate,
		TermsPrivacy:           payload.TermsPrivacy,
		UserName:               &userName,
		Verified:               false,
		UserType:               constants.UserTypeNormal,
		SuperAdminStatus:       constants.AdminStatusUser,
		LastPasswordChangeDate: &now,
		SignupIP:               signUpIp,
		IsActive:               true,
		SessionTime:            &userRoe.SessionTime,
		SessionCheckTime:       &userRoe.SessionCheckTime,
		Photo:                  helpers.Str("avatar.jpg"),
	}

	err = repository.CreateUser(&createdUser)

	if err != nil {

		return nil, err
	}

	return &createdUser, nil
}

func (s *Service) UpdateUserVerifyToken(token *models.EmailVerificationToken) error {

	now := time.Now()
	token.UsedAt = &now

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	_, err := repository.UpdateUserToken(token.Id, *token)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateUserVerifyStatus(user *models.User) error {

	user.ActiveStatus = 1
	user.Verified = true
	user.IsActive = true

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	_, err := repository.UpdateUser(user.Id, *user)

	if err != nil {
		return err
	}

	return nil
}
