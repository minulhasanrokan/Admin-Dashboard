package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type ActivityList struct {
	Id           uint64      `gorm:"primaryKey;autoIncrement;column:id" json:"id"`
	UUID         uuid.UUID `gorm:"type:char(36);uniqueIndex;not null" json:"uuid"`
	Name         string    `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Code         string    `gorm:"column:code;type:varchar(20);not null" json:"code"`
	DeleteStatus int       `gorm:"column:delete_status;default:0" json:"delete_status"`
	Status       int       `gorm:"column:status;default:1" json:"status"`
}

func (m *ActivityList) BeforeCreate(tx *gorm.DB) error {
	_ = tx

	if m.UUID == uuid.Nil {
		m.UUID = uuid.New()
	}

	return nil
}
