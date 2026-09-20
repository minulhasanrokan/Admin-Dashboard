package services

import (
	"admin-dashboard/cmd/internal/cacheSystem"
	"admin-dashboard/cmd/internal/config"
	"admin-dashboard/cmd/internal/databases"
	"admin-dashboard/cmd/internal/storage"

	"github.com/labstack/echo/v4"
)

type Service struct {
	Db          databases.DbConnection
	Logger      echo.Logger
	Storage     storage.Storage
	CacheSystem cacheSystem.CacheSystem
	Config      config.Config
}
