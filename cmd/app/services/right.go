package services

import (
	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/internal/dtos"
	"fmt"
)

func (s *Service) HasPermissionSuperAdmin(requiredAction string) (bool, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	count, err := repository.CheckRightDetailAction(requiredAction)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *Service) HasCategoryPermissionSuperAdmin(categoryAction string) (bool, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	count, err := repository.CheckCategoryRightDetailAction(categoryAction)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *Service) HasPermission(userID uint64, requiredAction string) (bool, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	count, err := repository.CheckUserRightKey(userID, requiredAction)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *Service) HasCategoryPermission(userID uint64, requiredAction string) (bool, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	count, err := repository.CheckUserCategoryRightKey(userID, requiredAction)

	if err != nil {
		return false, err
	}

	return count > 0, nil
}

func (s *Service) GetUserDashboardRights(userId uint64, isSuperAdmin bool) ([]dtos.RightGroupTreeDTO, error) {

	var resultTree []dtos.RightGroupTreeDTO

	cacheKey := fmt.Sprintf("UserRightDashboard%d%t", userId, isSuperAdmin)

	if err := s.CacheSystem.Get(cacheKey, &resultTree); err == nil {
		return resultTree, nil
	}

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	allGroups, err := repository.GetActiveRightGroups()

	if err != nil {
		return nil, err
	}

	allCategories, err := repository.GetActiveRightCategories()

	if err != nil {
		return nil, err
	}

	allDetails, err := repository.GetActiveRightDetails()

	if err != nil {
		return nil, err
	}

	allowedKeys := make(map[string]bool)

	if !isSuperAdmin {
		userRights, err := repository.GetUserRightKeys(userId)
		if err != nil {
			return nil, err
		}

		for _, key := range userRights {
			allowedKeys[key] = true
		}
	}

	allowedCategories := make(map[string]bool)

	for _, d := range allDetails {
		if isSuperAdmin || allowedKeys[d.Action] {
			allowedCategories[d.CategoryKey] = true
		}
	}

	groupToCategories := make(map[string][]dtos.RightCategoryDTO)

	for _, c := range allCategories {

		if !allowedCategories[c.Key] {
			continue
		}

		groupToCategories[c.GroupKey] = append(groupToCategories[c.GroupKey], dtos.RightCategoryDTO{
			Id:         c.Id,
			Name:       c.Name,
			Key:        c.Key,
			Url:        c.Url,
			Title:      c.Title,
			Action:     c.Action,
			Icon:       c.Icon,
			ShortOrder: c.ShortOrder,
		})
	}

	for _, g := range allGroups {
		categories, ok := groupToCategories[g.Key]
		if !ok || len(categories) == 0 {
			continue
		}

		resultTree = append(resultTree, dtos.RightGroupTreeDTO{
			Id:         g.Id,
			Name:       g.Name,
			Key:        g.Key,
			Title:      g.Title,
			Icon:       g.Icon,
			ShortOrder: g.ShortOrder,
			Categories: categories,
		})
	}

	_ = s.CacheSystem.Set(cacheKey, resultTree)

	return resultTree, nil
}

func (s *Service) GetRightDetailsByKey(action string) (dtos.RightInfo, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	return repository.GetRightDetailsByKey(action)
}

func (s *Service) GetRightDetailsByCategoryKeyKey(action string) ([]dtos.RightInfo, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	return repository.GetRightDetailsByCategoryKey(action)
}

func (s *Service) GetRightAllDetailsByKey(action string) ([]dtos.RightInfo, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	return repository.GetRightAllDetailsByKey(action)
}

func (s *Service) GetRightDetailsByCategoryKeyUserId(userid uint64, action string) ([]dtos.RightInfo, error) {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	return repository.GetRightDetailsByCategoryKeyUserId(userid, action)
}
