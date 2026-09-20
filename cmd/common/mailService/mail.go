package mailService

import (
	"admin-dashboard/cmd/internal/databases"
	"admin-dashboard/cmd/internal/mailer"

	"github.com/labstack/echo/v4"
)

type SendMailService struct {
	Db     databases.DbConnection
	Mailer mailer.MailConnection
	Logger echo.Logger
}
