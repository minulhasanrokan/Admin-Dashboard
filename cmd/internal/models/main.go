package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MainBaseModel struct {
	Id            uint64      `gorm:"primaryKey;autoIncrement;type:int unsigned" json:"id"`
	UUID          uuid.UUID `gorm:"type:char(36);uniqueIndex;not null" json:"uuid"`
	Status        int       `gorm:"type:tinyint;default:1" json:"status"`
	DeleteStatus  int       `gorm:"type:tinyint;default:0" json:"delete_status"`
	ProcessStatus uint64      `gorm:"type:int unsigned;default:0" json:"process_status"`
	ActualStatus  uint64      `gorm:"type:int unsigned;default:0" json:"actual_status"`
	EditStatus    int       `gorm:"type:tinyint;default:0" json:"edit_status"`
	Message       string    `gorm:"type:varchar(255);default:''" json:"message"`
	ActiveStatus  uint64      `gorm:"type:tinyint;default:0" json:"active_status"`
	Checker       uint64      `gorm:"type:int unsigned;default:0" json:"checker"`
	NotifyType    uint64      `gorm:"type:int unsigned;default:0" json:"notify_type"`
	Remarks       string    `gorm:"type:text" json:"remarks"`
	CreatedAt     time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt     time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	DeletedAt     time.Time `gorm:"column:deleted_at" json:"deleted_at"`
}

func (m *MainBaseModel) BeforeCreate(tx *gorm.DB) error {

	_ = tx

	if m.UUID == uuid.Nil {
		m.UUID = uuid.New()
	}
	return nil
}
