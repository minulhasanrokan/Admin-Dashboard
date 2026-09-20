package repositories

import "admin-dashboard/cmd/internal/models"

func (r *Repository) CreateUserActivity(activity *models.UserActivity) error {

	if err := r.Db.GetDb().Create(activity).Error; err != nil {

		return err
	}

	return nil
}
