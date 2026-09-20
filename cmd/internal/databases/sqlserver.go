package databases

import (
	"admin-dashboard/cmd/internal/config"
	"fmt"
	"net/url"
	"time"

	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
)

func ConnectMssql(cfg *config.Config) (DbConnection, error) {

	host := cfg.DbHost
	port := cfg.DbPort
	database := cfg.DbDatabase
	userName := cfg.DbUsername
	password := cfg.DbPassword

	dsn := fmt.Sprintf(
		"sqlserver://%s:%s@%s:%s?database=%s",
		url.QueryEscape(userName),
		url.QueryEscape(password),
		host,
		port,
		url.QueryEscape(database),
	)

	db, err := gorm.Open(
		sqlserver.Open(dsn),
		&gorm.Config{
			TranslateError: true,
		},
	)

	if err != nil {
		return nil, err
	}

	if sqlDb, err := db.DB(); err == nil {
		sqlDb.SetMaxOpenConns(50)
		sqlDb.SetMaxIdleConns(25)
		sqlDb.SetConnMaxLifetime(60 * time.Minute)
		sqlDb.SetConnMaxIdleTime(15 * time.Minute)
	}

	return &GormDb{
		DB:             db,
		NotFoundString: cfg.DbRecordNotFoundMessage,
		Driver:         cfg.DbConnection,
		Host:           host,
		Port:           port,
		Database:       database,
		Username:       userName,
		Password:       password,
	}, nil
}
