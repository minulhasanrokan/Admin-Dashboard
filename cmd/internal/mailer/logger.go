package mailer

import (
	"admin-dashboard/cmd/internal/databases"
	"admin-dashboard/cmd/internal/models"
	"time"
)

type MailLog struct {
	FromAddress    string
	ToAddress      string
	BccAddress     string
	AttachmentPath string
	CcAddress      string
	TableName      string
	TableId        uint64
	Subject        string
	Body           string
	SentStatus     int
	TryStatus      int
	ErrorStatus    int
	ErrorMsg       string
	MailType       string
	CreatedAt      *time.Time
}

type MailLogger interface {
	InsertMailLog(log MailLog) error
}

type MailServiceLogger struct {
	Db databases.DbConnection
}

func NewMailServiceLogger(db databases.DbConnection) *MailServiceLogger {
	return &MailServiceLogger{
		Db: db,
	}
}

func (m *MailServiceLogger) InsertMailLog(log MailLog) error {
	mailHistory := &models.MailHistory{
		FromAddress: log.FromAddress,
		ToAddress:   log.ToAddress,
		BccAddress:  log.BccAddress,
		Attachment:  log.AttachmentPath,
		CcAddress:   log.CcAddress,
		TableName:   log.TableName,
		TableId:     log.TableId,
		Subject:     log.Subject,
		Body:        log.Body,
		SentStatus:  log.SentStatus,
		TryStatus:   log.TryStatus,
		ErrorStatus: log.ErrorStatus,
		ErrorMsg:    log.ErrorMsg,
	}

	return m.Db.GetDb().Create(mailHistory).Error
}
