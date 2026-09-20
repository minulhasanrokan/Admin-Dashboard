package repositories

import "admin-dashboard/cmd/internal/models"

func (r *Repository) CreateMailHistory(mailHistory *models.MailHistory) error {

	if err := r.Db.GetDb().Create(mailHistory).Error; err != nil {
		return err
	}

	return nil
}

func (r *Repository) GetPendingMailHistories(limit int, maxTries int) ([]models.MailHistory, error) {

	var mailHistories []models.MailHistory

	err := r.Db.GetDb().
		Model(&models.MailHistory{}).
		Where("sent_status = ? AND try_status < ?", 0, maxTries).
		Order("id asc").
		Limit(limit).
		Find(&mailHistories).Error

	if err != nil {
		return nil, err
	}

	return mailHistories, nil
}

func (r *Repository) UpdateMailHistory(mailHistory *models.MailHistory) error {

	return r.Db.GetDb().Save(mailHistory).Error
}
