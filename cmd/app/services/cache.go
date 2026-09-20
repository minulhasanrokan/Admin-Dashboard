package services

import "admin-dashboard/cmd/app/repositories"

func (s *Service) ClearCache() error {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	return repository.ClearCache()
}
