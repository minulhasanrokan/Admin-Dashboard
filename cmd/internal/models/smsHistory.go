package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type SmsHistory struct {
	Id          uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UUID        uuid.UUID `gorm:"type:char(36);uniqueIndex;not null" json:"uuid"`
	ToMobile    string    `gorm:"column:to_mobile;type:varchar(50);not null" json:"to_mobile"`
	Message     string    `gorm:"column:message;type:text;not null" json:"message"`
	SmsDriver   string    `gorm:"column:sms_driver;type:varchar(50);not null" json:"sms_driver"`
	SmsType     string    `gorm:"column:sms_type;type:varchar(100);not null" json:"sms_type"`
	SentStatus  int       `gorm:"column:sent_status;type:int;default:0" json:"sent_status"`
	TryStatus   int       `gorm:"column:try_status;type:int;default:0" json:"try_status"`
	ErrorStatus int       `gorm:"column:error_status;type:int;default:0" json:"error_status"`
	ErrorMsg    string    `gorm:"column:error_msg;type:text;not null" json:"error_msg"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (m *SmsHistory) BeforeCreate(tx *gorm.DB) error {
	_ = tx

	if m.UUID == uuid.Nil {
		m.UUID = uuid.New()
	}

	return nil
}
