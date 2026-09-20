package repositories

import (
	"admin-dashboard/cmd/internal/models"
)

func (r *Repository) SystemSetting() (models.SystemSetting, error) {
	var setting models.SystemSetting

	err := r.CacheSystem.Get("SystemSetting", &setting)

	if err == nil {
		return setting, nil
	}

	err = r.Db.GetDb().Model(&models.SystemSetting{}).First(&setting).Error

	if err != nil {
		return models.SystemSetting{}, err
	}

	_ = r.CacheSystem.Set("SystemSetting", setting)

	return setting, nil
}

func (r *Repository) UpdateSystemSetting(setting *models.SystemSetting) error {

	if err := r.Db.GetDb().Save(setting).Error; err != nil {
		return err
	}

	return r.CacheSystem.Set("SystemSetting", setting)
}
