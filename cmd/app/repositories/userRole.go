package repositories

import (
	"admin-dashboard/cmd/internal/models"
)

func (r *Repository) UserRole() (models.UserRole, error) {
	var data models.UserRole

	if err := r.CacheSystem.Get("UserRole", &data); err == nil {
		return data, nil
	}

	if err := r.Db.GetDb().Model(&models.UserRole{}).First(&data).Error; err != nil {
		return models.UserRole{}, err
	}

	_ = r.CacheSystem.Set("UserRole", data)

	return data, nil
}

func (r *Repository) UpdateUserRoleSetting(userRole *models.UserRole) error {

	if err := r.Db.GetDb().Save(userRole).Error; err != nil {
		return err
	}

	return r.CacheSystem.Set("UserRole", userRole)
}
