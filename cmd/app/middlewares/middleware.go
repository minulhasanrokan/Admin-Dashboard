package middlewares

import (
	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/app/services"
	"admin-dashboard/cmd/internal/cacheSystem"
	"admin-dashboard/cmd/internal/config"
	"admin-dashboard/cmd/internal/databases"

	"github.com/labstack/echo/v4"
)

type Middleware struct {
	Logger      echo.Logger
	Db          databases.DbConnection
	Config      config.Config
	Repository  repositories.Repository
	Service     services.Service
	CacheSystem cacheSystem.CacheSystem
}
