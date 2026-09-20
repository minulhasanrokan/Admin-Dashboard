package services

import (
	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/internal/models"
)

func (s *Service) BackupHistory() (*models.BackupHistory, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	data, err := repository.BackupHistory()

	if err != nil {

		return nil, err
	}

	return &data, nil
}
