package repositories

import (
	"admin-dashboard/cmd/internal/models"

	"gorm.io/gorm/clause"
)

func (r *Repository) CreateUserLoginHistory(userLoginHistory *models.UserLoginHistory) error {

	userLoginHistory.Id = 0

	err := r.Db.GetDb().Create(userLoginHistory).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateUserLoginHistory(userId uint64, data models.UserLoginHistory) ([]models.UserLoginHistory, error) {

	var updatedHistories []models.UserLoginHistory

	userIdInt := int(userId)

	err := r.Db.GetDb().Model(&updatedHistories).Clauses(clause.Returning{}).Where("user_id = ? AND log_out_status = ? AND force_status = ?", userIdInt, 0, 0).Updates(data).Error

	if err != nil {
		return nil, err
	}

	return updatedHistories, nil
}
