package databases

import (
	"admin-dashboard/cmd/internal/config"
	"fmt"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectPostgresSQL(cfg *config.Config) (DbConnection, error) {

	host := cfg.DbHost
	port := cfg.DbPort
	database := cfg.DbDatabase
	userName := cfg.DbUsername
	password := cfg.DbPassword

	dsn := fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host,
		port,
		userName,
		password,
		database,
	)

	db, err := gorm.Open(
		postgres.Open(dsn),
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
