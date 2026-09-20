package repositories

import (
	"admin-dashboard/cmd/internal/cacheSystem"
	"admin-dashboard/cmd/internal/databases"

	"github.com/labstack/echo/v4"
)

type Repository struct {
	Db          databases.DbConnection
	Logger      echo.Logger
	CacheSystem cacheSystem.CacheSystem
}

func NewRepository(db *databases.DbConnection, cache cacheSystem.CacheSystem) *Repository {
	return &Repository{Db: *db, CacheSystem: cache}
}
