package services

import (
	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/common/helpers"
	"admin-dashboard/cmd/internal/models"
	"time"
)

func (s *Service) UpdateUserResetPasswordToken(token *models.PasswordResetToken) error {

	now := time.Now()
	token.UsedAt = &now

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	_, err := repository.UpdateUserPasswordToken(token.Id, *token)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateUserPassword(user *models.User, password string) error {

	hashedPassword, err := helpers.HashPassword(password)

	if err != nil {

		return err
	}

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	userRole, err := repository.UserRole()

	if err != nil {

		return err
	}

	now := time.Now()

	passwordExpireDate := now.AddDate(0, 0, int(userRole.PasswordExpireDays))

	user.Password = &hashedPassword
	user.PasswordChangeStatus = helpers.UInt(1)
	user.PasswordChangeDate = &now
	user.LastPasswordChangeDate = &now
	user.PasswordExpireDate = &passwordExpireDate

	_, err = repository.UpdateUser(user.Id, *user)

	if err != nil {

		return err
	}

	return nil
}
