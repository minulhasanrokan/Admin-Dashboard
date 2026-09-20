package models

import (
	"time"
)

type Session struct {
	Id           string    `gorm:"primaryKey;type:varchar(128);column:id" json:"id"`
	UserId       *uint64   `gorm:"column:user_id;index" json:"user_id"`
	IpAddress    string    `gorm:"column:ip_address;type:varchar(45)" json:"ip_address"`
	UserAgent    string    `gorm:"column:user_agent;type:text" json:"user_agent"`
	Payload      string    `gorm:"column:payload;type:longtext;not null" json:"payload"`
	LastActivity int64     `gorm:"column:last_activity;index;not null" json:"last_activity"`
	ExpiresAt    time.Time `gorm:"column:expires_at;index;not null" json:"expires_at"`
	CreatedAt    time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime;column:updated_at" json:"updated_at"`
}

func (Session) TableName() string {
	return "sessions"
}
