package repositories

import (
	"admin-dashboard/cmd/internal/models"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm/clause"
)

func (r *Repository) CreateUser(user *models.User) error {

	user.Id = 0

	err := r.Db.GetDb().Create(user).Error

	if err != nil {
		return err
	}

	return nil
}

func (r *Repository) UpdateUserToken(id uint64, data models.EmailVerificationToken) (*models.EmailVerificationToken, error) {

	var token models.EmailVerificationToken

	err := r.Db.GetDb().Model(&token).Clauses(clause.Returning{}).Where("id=?", id).Updates(data).Error

	if err != nil {

		return &models.EmailVerificationToken{}, err
	}

	return &token, nil
}

func (r *Repository) GenerateUniqueUserName(Name string) (string, error) {

	base := strings.ToLower(strings.ReplaceAll(Name, " ", ""))

	userName := base

	suffix := 0

	for {
		existingUser, err := r.GetUserByUserName(userName)

		if err != nil && err.Error() != r.Db.NotFoundMessage() {
			return "", err
		}

		if existingUser == nil {
			break
		}
		suffix++
		userName = fmt.Sprintf("%s%d", base, suffix)
	}

	return userName, nil
}

func (r *Repository) GetActiveVerifyTokenUserId(userId uint64) (*models.EmailVerificationToken, error) {

	var token models.EmailVerificationToken

	err := r.Db.GetDb().Order("id DESC").Where("user_id = ? AND used_at IS NULL AND expires_at > ?", userId, time.Now()).First(&token).Error

	if err != nil {

		return nil, err
	}

	return &token, nil
}
