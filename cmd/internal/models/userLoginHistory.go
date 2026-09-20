package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserLoginHistory struct {
	Id             uint64       `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UUID           uuid.UUID  `gorm:"type:char(36);uniqueIndex;not null" json:"uuid"`
	UserId         *int       `gorm:"column:user_id;type:int" json:"user_id"`
	IpAddress      *string    `gorm:"column:ip_address;type:varchar(45)" json:"ip_address"`
	UserAgent      *string    `gorm:"column:user_agent;type:text" json:"user_agent"`
	LogInTime      *time.Time `gorm:"column:log_in_time;type:datetime" json:"log_in_time"`
	LogOutTime     *time.Time `gorm:"column:log_out_time;type:datetime" json:"log_out_time"`
	LastActiveTime *time.Time `gorm:"column:last_active_time;type:datetime" json:"last_active_time"`
	LogOutStatus   *int       `gorm:"column:log_out_status;type:int;default:0" json:"log_out_status"`
	ForceStatus    *int       `gorm:"column:force_status;type:int;default:0" json:"force_status"`
	DeleteStatus   *int       `gorm:"column:delete_status;type:int;default:0" json:"delete_status"`
	Status         *int       `gorm:"column:status;type:int;default:1" json:"status"`
	CreatedAt      time.Time  `gorm:"column:created_at;type:timestamp" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"column:updated_at;type:timestamp" json:"updated_at"`
}

func (m *UserLoginHistory) BeforeCreate(tx *gorm.DB) error {

	_ = tx

	if m.UUID == uuid.Nil {
		m.UUID = uuid.New()
	}
	return nil
}
