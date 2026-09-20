package services

import (
	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/internal/models"
)

func (s *Service) PrivacyPolicy() (*models.PrivacyPolicy, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	data, err := repository.TermsPrivacy()

	if err != nil {

		return nil, err
	}

	return &data, nil
}
