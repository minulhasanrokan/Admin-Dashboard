package databases

import (
	"admin-dashboard/cmd/internal/config"
	"errors"

	"gorm.io/gorm"
)

type DbConnection interface {
	GetDb() *gorm.DB
	CloseDb() error
	NotFoundMessage() string

	GetDriver() string
	GetHost() string
	GetPort() string
	GetDatabase() string
	GetUsername() string
	GetPassword() string
}

type GormDb struct {
	DB             *gorm.DB
	NotFoundString string
	Driver         string
	Host           string
	Port           string
	Database       string
	Username       string
	Password       string
}

func (g *GormDb) GetDb() *gorm.DB {
	return g.DB
}

func (g *GormDb) CloseDb() error {

	sqlDb, err := g.DB.DB()
	if err != nil {
		return err
	}

	return sqlDb.Close()
}

func (g *GormDb) NotFoundMessage() string {
	return g.NotFoundString
}

func (g *GormDb) GetDriver() string {
	return g.Driver
}

func (g *GormDb) GetHost() string {
	return g.Host
}

func (g *GormDb) GetPort() string {
	return g.Port
}

func (g *GormDb) GetDatabase() string {
	return g.Database
}

func (g *GormDb) GetUsername() string {
	return g.Username
}

func (g *GormDb) GetPassword() string {
	return g.Password
}

func ConnectDatabase(cfg *config.Config) (DbConnection, error) {

	if cfg.DbConnection == "" {
		return nil, errors.New("DB_DRIVER environment variable not set")
	}

	switch cfg.DbConnection {

	case "mysql":
		return ConnectMsql(cfg)
	case "mssql", "sqlserver":
		return ConnectMssql(cfg)
	case "postgres", "postgresql":
		return ConnectPostgresSQL(cfg)
	default:
		return nil, errors.New(
			cfg.DbConnection + " is not supported",
		)
	}
}
