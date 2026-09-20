package repositories

import (
	"admin-dashboard/cmd/internal/models"
	"time"

	"gorm.io/gorm/clause"
)

func (r *Repository) GetActivePasswordTokenUserId(userId uint64) (*models.PasswordResetToken, error) {

	var token models.PasswordResetToken

	err := r.Db.GetDb().Order("id DESC").Where("user_id = ? AND used_at IS NULL AND expires_at > ?", userId, time.Now()).First(&token).Error

	if err != nil {

		return nil, err
	}

	return &token, nil
}

func (r *Repository) UpdateUserPasswordToken(id uint64, data models.PasswordResetToken) (*models.PasswordResetToken, error) {

	var token models.PasswordResetToken

	err := r.Db.GetDb().Model(&token).Clauses(clause.Returning{}).Where("id=?", id).Updates(data).Error

	if err != nil {

		return &models.PasswordResetToken{}, err
	}

	return &token, nil
}
