package repositories

import (
	"admin-dashboard/cmd/internal/models"
)

func (r *Repository) SmsSetting() (models.SmsSetting, error) {
	var data models.SmsSetting

	if err := r.CacheSystem.Get("SmsSetting", &data); err == nil {
		return data, nil
	}

	if err := r.Db.GetDb().Model(&models.SmsSetting{}).First(&data).Error; err != nil {
		return models.SmsSetting{}, err
	}

	_ = r.CacheSystem.Set("SmsSetting", data)

	return data, nil
}

func (r *Repository) UpdateSmsSetting(setting *models.SmsSetting) error {

	if err := r.Db.GetDb().Save(setting).Error; err != nil {
		return err
	}

	return r.CacheSystem.Set("SmsSetting", *setting)
}
