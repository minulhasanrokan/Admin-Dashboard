package applicationData

import (
	"admin-dashboard/cmd/app/handlers"
	"admin-dashboard/cmd/app/middlewares"
	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/app/services"
	"admin-dashboard/cmd/common/mailService"
	"admin-dashboard/cmd/internal/cacheSystem"
	"admin-dashboard/cmd/internal/storage"

	"github.com/labstack/echo/v4"
)

type Application struct {
	Logger      echo.Logger
	Server      *echo.Echo
	Handler     handlers.Handler
	Middleware  middlewares.Middleware
	Repository  repositories.Repository
	Service     services.Service
	Mailer      mailService.SendMailService
	Storage     storage.Storage
	CacheSystem cacheSystem.CacheSystem
}
