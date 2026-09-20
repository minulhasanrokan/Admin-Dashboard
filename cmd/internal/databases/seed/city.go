package seed

import (
	"admin-dashboard/cmd/internal/models"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"

	"gorm.io/gorm"
)

func CitySeed(db *gorm.DB) error {
	var count int64

	if err := db.Model(&models.City{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing cities: %w", err)
	}

	if count > 0 {
		return nil
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to recover runtime caller metadata")
	}

	basePath := filepath.Dir(filename)
	csvPath := filepath.Join(basePath, "seed_data", "cities.csv")

	file, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("failed to open cities CSV: %w", err)
	}
	defer func() {
		_ = file.Close()
	}()

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("failed to read cities CSV header: %w", err)
	}

	var cities []models.City

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("error reading cities CSV record: %w", err)
		}

		if len(record) < 4 {
			return fmt.Errorf("malformed CSV row: expected 4 columns, got %d", len(record))
		}

		cities = append(cities, models.City{
			Name:        record[0],
			Code:        record[1],
			StateCode:   record[2],
			CountryCode: record[3],
		})
	}

	if len(cities) > 0 {
		if err := db.Create(&cities).Error; err != nil {
			return fmt.Errorf("failed to batch insert cities: %w", err)
		}
	}

	return nil
}
