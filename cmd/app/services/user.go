package services

import (
	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/internal/models"
	"time"
)

func (s *Service) UpdateUserLastLoginTime(user *models.User) error {

	now := time.Now()
	user.LastLogin = &now

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	_, err := repository.UpdateUser(user.Id, *user)

	if err != nil {
		return err
	}

	return nil
}

func (s *Service) GetUserByID(id uint64) (*models.User, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	return repository.GetUserByID(id)
}

func (s *Service) UpdateUserProfile(userId uint64, name, mobile, address, details string) (*models.User, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	updates := map[string]interface{}{
		"name":    name,
		"mobile":  mobile,
		"address": address,
		"details": details,
	}
	return repository.UpdateUserMap(userId, updates)
}

func (s *Service) UpdateUserPhoto(userId uint64, photo string) (*models.User, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	updates := map[string]interface{}{
		"photo": photo,
	}
	return repository.UpdateUserMap(userId, updates)
}

