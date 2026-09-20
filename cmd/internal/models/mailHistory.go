package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type MailHistory struct {
	Id          uint64      `gorm:"column:id;primaryKey;autoIncrement" json:"id"`
	UUID        uuid.UUID `gorm:"type:char(36);uniqueIndex;not null" json:"uuid"`
	FromAddress string    `gorm:"column:from_address;type:varchar(255);not null" json:"from_address"`
	ToAddress   string    `gorm:"column:to_address;type:text;not null" json:"to_address"`
	BccAddress  string    `gorm:"column:bcc_address;type:text;not null" json:"bcc_address"`
	Attachment  string    `gorm:"column:attachment;type:text;not null" json:"attachment"`
	CcAddress   string    `gorm:"column:cc_address;type:text" json:"cc_address"`
	MailType    string    `gorm:"column:mail_type;type:varchar(255);not null" json:"mail_type"`
	TableName   string    `gorm:"column:table_name;type:varchar(255);not null" json:"table_name"`
	TableId     uint64      `gorm:"column:table_id;type:int" json:"table_id"`
	Subject     string    `gorm:"column:subject;type:varchar(255);not null" json:"subject"`
	Body        string    `gorm:"column:body;type:text;not null" json:"body"`
	SentStatus  int       `gorm:"column:sent_status;type:int;default:0" json:"sent_status"`
	TryStatus   int       `gorm:"column:try_status;type:int;default:0" json:"try_status"`
	ErrorStatus int       `gorm:"column:error_status;type:int;default:0" json:"error_status"`
	ErrorMsg    string    `gorm:"column:error_msg;type:text;not null" json:"error_msg"`
	CreatedAt   time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt   time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

func (m *MailHistory) BeforeCreate(tx *gorm.DB) error {
	_ = tx

	if m.UUID == uuid.Nil {
		m.UUID = uuid.New()
	}

	return nil
}
