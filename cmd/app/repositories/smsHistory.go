package repositories

import "admin-dashboard/cmd/internal/models"

func (r *Repository) CreateSmsHistory(smsHistory *models.SmsHistory) error {

	if err := r.Db.GetDb().Create(smsHistory).Error; err != nil {

		return err
	}

	return nil
}

func (r *Repository) GetPendingSmsHistories(limit int, maxTries int) ([]models.SmsHistory, error) {

	var smsHistories []models.SmsHistory

	err := r.Db.GetDb().
		Model(&models.SmsHistory{}).
		Where("sent_status = ? AND try_status < ?", 0, maxTries).
		Order("id asc").
		Limit(limit).
		Find(&smsHistories).Error

	if err != nil {
		return nil, err
	}

	return smsHistories, nil
}

func (r *Repository) UpdateSmsHistory(smsHistory *models.SmsHistory) error {

	return r.Db.GetDb().Save(smsHistory).Error
}

