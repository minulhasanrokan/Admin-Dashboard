package repositories

import (
	"admin-dashboard/cmd/internal/models"

	"gorm.io/gorm/clause"
)

func (r *Repository) GetUserByEmail(email string) (*models.User, error) {

	var user models.User

	err := r.Db.GetDb().Where("email = ?", email).First(&user).Error

	if err != nil {

		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetUserByMobile(mobile string) (*models.User, error) {

	var user models.User

	err := r.Db.GetDb().Where("mobile = ?", mobile).First(&user).Error

	if err != nil {

		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetUserByUserName(username string) (*models.User, error) {

	var user models.User

	err := r.Db.GetDb().Where("user_name = ?", username).First(&user).Error

	if err != nil {

		return nil, err
	}

	return &user, nil
}

func (r *Repository) GetUserByUuId(uuid string) (*models.User, error) {

	var user models.User

	err := r.Db.GetDb().Where("uuid = ?", uuid).First(&user).Error

	if err != nil {

		return nil, err
	}

	return &user, nil
}

func (r *Repository) UpdateUser(id uint64, data models.User) (*models.User, error) {

	var user models.User

	err := r.Db.GetDb().Model(&user).Clauses(clause.Returning{}).Where("id=?", id).Updates(data).Error

	if err != nil {

		return &models.User{}, err
	}

	return &user, nil
}

func (r *Repository) GetUserByID(id uint64) (*models.User, error) {
	var user models.User
	err := r.Db.GetDb().Where("id = ?", id).First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) UpdateUserMap(id uint64, data map[string]interface{}) (*models.User, error) {
	var user models.User
	err := r.Db.GetDb().Model(&user).Clauses(clause.Returning{}).Where("id = ?", id).Updates(data).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

