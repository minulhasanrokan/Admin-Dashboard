package repositories

import (
	"admin-dashboard/cmd/internal/models"
)

func (r *Repository) TermsPrivacy() (models.PrivacyPolicy, error) {

	var privacyPolicy models.PrivacyPolicy

	err := r.CacheSystem.Get("PrivacyPolicy", &privacyPolicy)

	if err == nil {

		return privacyPolicy, nil
	}

	err = r.Db.GetDb().Model(&models.PrivacyPolicy{}).First(&privacyPolicy).Error

	if err != nil {

		return models.PrivacyPolicy{}, err
	}

	_ = r.CacheSystem.Set("PrivacyPolicy", privacyPolicy)

	return privacyPolicy, nil
}
