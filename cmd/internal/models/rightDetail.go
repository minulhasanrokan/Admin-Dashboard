package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
)

type RightDetail struct {
	Id          uint64    `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UUID        uuid.UUID `gorm:"type:char(36);uniqueIndex;not null" json:"uuid"`
	CategoryKey string    `gorm:"column:category_key;type:varchar(255)" json:"category_key"`
	Name        string    `gorm:"column:name;type:varchar(255)" json:"name"`
	ActionName  string    `gorm:"column:action_name;type:varchar(100)" json:"action_name"`
	Url         string    `gorm:"column:url;type:varchar(255)" json:"url"`
	Title       string    `gorm:"column:title;type:varchar(255)" json:"title"`
	Action      string    `gorm:"column:action;type:varchar(255)" json:"action"`
	Details     string    `gorm:"column:details;type:text" json:"details"`
	PopupStatus int       `gorm:"column:popup_status;type:int;default:0" json:"popup_status"`
	Width       int       `gorm:"column:width;type:int;default:1000" json:"width"`
	Icon        string    `gorm:"column:icon;type:varchar(255)" json:"icon"`
	SortOrder  int       `gorm:"column:sort_order;type:int" json:"sort_order"`
	Status      int       `gorm:"column:status;type:int;default:1" json:"status"`
	ShortName   string    `gorm:"column:short_name;type:varchar(100)" json:"short_name"`
}

func (m *RightDetail) BeforeCreate(tx *gorm.DB) error {

	_ = tx

	if m.UUID == uuid.Nil {
		m.UUID = uuid.New()
	}
	return nil
}
