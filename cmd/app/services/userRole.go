package services

import (
	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/app/requests"
	"admin-dashboard/cmd/internal/models"
)

func (s *Service) UserRole() (*models.UserRole, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	data, err := repository.UserRole()

	if err != nil {

		return nil, err
	}

	return &data, nil
}

func (s *Service) UserRoleSettingUpdate(payload *requests.UserRoleRequest) (*models.UserRole, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	userRole, err := repository.UserRole()

	if err != nil {
		return nil, err
	}

	userRole.PasswordMinLength = payload.PasswordMinLength
	userRole.PasswordMaxLength = payload.PasswordMaxLength
	userRole.SessionTime = payload.SessionTime
	userRole.SessionCheckTime = payload.SessionCheckTime
	userRole.PasswordExpireDays = payload.PasswordExpireDays

	if err := repository.UpdateUserRoleSetting(&userRole); err != nil {
		return nil, err
	}

	return &userRole, nil
}
