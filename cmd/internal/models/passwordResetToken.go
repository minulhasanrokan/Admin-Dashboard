package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type PasswordResetToken struct {
	Id        uint64       `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	UUID      uuid.UUID  `gorm:"type:char(36);uniqueIndex;not null" json:"uuid"`
	UserId    uint64       `gorm:"column:user_id;index;not null" json:"user_id"`
	TokenHash string     `gorm:"column:token_hash;uniqueIndex;type:varchar(64);not null" json:"token_hash"`
	OtpToken  string     `gorm:"column:otp_token;uniqueIndex;type:varchar(64);not null" json:"otp_token"`
	ExpiresAt time.Time  `gorm:"column:expires_at;not null" json:"expires_at"`
	UsedAt    *time.Time `gorm:"column:used_at" json:"used_at"`
	CreatedAt time.Time  `gorm:"autoCreateTime;column:created_at" json:"created_at"`
}

func (m *PasswordResetToken) BeforeCreate(tx *gorm.DB) error {

	_ = tx

	if m.UUID == uuid.Nil {
		m.UUID = uuid.New()
	}
	return nil
}
