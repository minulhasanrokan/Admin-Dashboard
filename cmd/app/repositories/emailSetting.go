package repositories

import "admin-dashboard/cmd/internal/models"

func (r *Repository) EmailSetting() (models.EmailSetting, error) {

	var emailSetting models.EmailSetting

	err := r.CacheSystem.Get("EmailSetting", &emailSetting)

	if err == nil {

		return emailSetting, nil
	}

	err = r.Db.GetDb().Model(&models.EmailSetting{}).First(&emailSetting).Error

	if err != nil {

		return models.EmailSetting{}, err
	}

	_ = r.CacheSystem.Set("EmailSetting", emailSetting)

	return emailSetting, nil
}

func (r *Repository) UpdateEmailSetting(emailSetting *models.EmailSetting) error {

	if err := r.Db.GetDb().Save(emailSetting).Error; err != nil {
		return err
	}

	return r.CacheSystem.Set("EmailSetting", emailSetting)
}
