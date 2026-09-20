package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type UserActivity struct {
	Id             uint64      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	UUID           uuid.UUID `gorm:"type:char(36);uniqueIndex;not null" json:"uuid"`
	GroupKey       *string   `gorm:"column:group_key;type:varchar(255)" json:"group_key"`
	CategoryKey    *string   `gorm:"column:category_key;type:varchar(255)" json:"category_key"`
	RightKey       *string   `gorm:"column:right_key;type:varchar(255)" json:"right_key"`
	ActivityStatus *string   `gorm:"column:activity_status;type:varchar(255)" json:"activity_status"`
	TableName      *string   `gorm:"column:table_name;type:varchar(255)" json:"table_name"`
	TableId        *uint64     `gorm:"column:table_id" json:"table_id"`
	Activity       *string   `gorm:"column:activity;type:varchar(255)" json:"activity"`
	IPAddress      *string   `gorm:"column:ip_address;type:varchar(100)" json:"ip_address"`
	MacAddress     *string   `gorm:"column:mac_address;type:varchar(100)" json:"mac_address"`
	Device         *string   `gorm:"column:device;type:varchar(255)" json:"device"`
	Status         int       `gorm:"column:status;default:1" json:"status"`
	CreatedBy      *uint64     `gorm:"column:add_by" json:"created_by"`
	CreatedAt      time.Time `gorm:"autoCreateTime;column:created_at" json:"created_at"`
}

func (m *UserActivity) BeforeCreate(tx *gorm.DB) error {
	_ = tx

	if m.UUID == uuid.Nil {
		m.UUID = uuid.New()
	}

	return nil
}
