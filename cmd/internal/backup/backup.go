package backup

import (
	"fmt"

	"gorm.io/gorm"
)

type Config struct {
	Driver   string
	Host     string
	Port     string
	Username string
	Password string
	Database string
	Db       *gorm.DB
}

type DatabaseBackup interface {
	Backup(outputPath string) error
}

func NewDatabaseBackup(cfg Config) (DatabaseBackup, error) {

	switch cfg.Driver {

	case "mysql", "mariadb":
		return &MySQLBackup{
			Host:     cfg.Host,
			Port:     cfg.Port,
			Username: cfg.Username,
			Password: cfg.Password,
			Database: cfg.Database,
			Db:       cfg.Db,
		}, nil

	case "postgres", "postgresql":
		return &PostgreSQLBackup{
			Host:     cfg.Host,
			Port:     cfg.Port,
			Username: cfg.Username,
			Password: cfg.Password,
			Database: cfg.Database,
			Db:       cfg.Db,
		}, nil

	case "mssql", "sqlserver":
		return &SQLServerBackup{
			Host:     cfg.Host,
			Port:     cfg.Port,
			Username: cfg.Username,
			Password: cfg.Password,
			Database: cfg.Database,
			Db:       cfg.Db,
		}, nil

	default:
		return nil, fmt.Errorf("unsupported database driver: %s", cfg.Driver)
	}
}
