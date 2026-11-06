package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RightGroup struct {
	Id         uint64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UUID       uuid.UUID `gorm:"type:char(36);uniqueIndex;not null" json:"uuid"`
	Name       string    `gorm:"column:name;type:varchar(255)" json:"name"`
	Key        string    `gorm:"column:key;type:varchar(255)" json:"key"`
	Title      string    `gorm:"column:title;type:varchar(255)" json:"title"`
	Details    string    `gorm:"column:details;type:text" json:"details"`
	Icon       string    `gorm:"column:icon;type:varchar(255)" json:"icon"`
	SortOrder int       `gorm:"column:sort_order;type:int" json:"sort_order"`
	Status     int       `gorm:"column:status;type:int;default:1" json:"status"`
}

func (m *RightGroup) BeforeCreate(tx *gorm.DB) error {

	_ = tx

	if m.UUID == uuid.Nil {
		m.UUID = uuid.New()
	}
	return nil
}
