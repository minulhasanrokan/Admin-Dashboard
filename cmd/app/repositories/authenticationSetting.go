package repositories

import (
	"admin-dashboard/cmd/internal/models"
)

func (r *Repository) AuthenticationSetting() (models.AuthenticationSetting, error) {
	var data models.AuthenticationSetting

	if err := r.CacheSystem.Get("AuthenticationSetting", &data); err == nil {
		return data, nil
	}

	if err := r.Db.GetDb().Model(&models.AuthenticationSetting{}).First(&data).Error; err != nil {
		return models.AuthenticationSetting{}, err
	}

	_ = r.CacheSystem.Set("AuthenticationSetting", data)

	return data, nil
}

func (r *Repository) UpdateAuthenticationSetting(setting *models.AuthenticationSetting) error {

	if err := r.Db.GetDb().Save(setting).Error; err != nil {
		return err
	}

	return r.CacheSystem.Set("AuthenticationSetting", *setting)
}
