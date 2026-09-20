package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type BackupHistory struct {
	Id           uint64           `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UUID         uuid.UUID      `gorm:"column:uuid;type:char(36);uniqueIndex;not null" json:"uuid"`
	FileName     string         `gorm:"column:file_name;type:varchar(255);not null" json:"file_name"`
	FilePath     string         `gorm:"column:file_path;type:varchar(500);not null" json:"file_path"`
	FileSize     uint64         `gorm:"column:file_size;type:bigint unsigned;default:0" json:"file_size"`
	BackupType   string         `gorm:"column:backup_type;type:varchar(100);not null" json:"backup_type"`
	Status       int            `gorm:"column:status;type:int;not null;default:1" json:"status"`
	BackupStatus   int            `gorm:"column:backup_status;type:int;not null;default:1" json:"backup_status"`
	DownloadStatus int            `gorm:"column:download_status;type:int;not null;default:0" json:"download_status"`
	ErrorLog       string         `gorm:"column:error_log;type:text" json:"error_log"`
	CreatedBy      *uint64        `gorm:"column:created_by" json:"created_by"`
	DownloadedBy   *uint64        `gorm:"column:downloaded_by" json:"downloaded_by"`
	DownloadedAt   *time.Time     `gorm:"column:downloaded_at" json:"downloaded_at"`
	CreatedAt      time.Time      `gorm:"column:created_at;autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time      `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
	DeletedAt      gorm.DeletedAt `gorm:"column:deleted_at;index" json:"deleted_at"`
}

func (m *BackupHistory) BeforeCreate(tx *gorm.DB) error {
	_ = tx

	if m.UUID == uuid.Nil {
		m.UUID = uuid.New()
	}

	return nil
}
