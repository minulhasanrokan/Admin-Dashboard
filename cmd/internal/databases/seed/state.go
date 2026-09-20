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

func StateSeed(db *gorm.DB) error {
	var count int64

	if err := db.Model(&models.State{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing states: %w", err)
	}

	if count > 0 {
		return nil
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to recover runtime caller metadata")
	}

	basePath := filepath.Dir(filename)
	csvPath := filepath.Join(basePath, "seed_data", "states.csv")

	file, err := os.Open(csvPath)

	if err != nil {
		return fmt.Errorf("failed to open states CSV: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("failed to read states CSV header: %w", err)
	}

	var states []models.State

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("error reading states CSV record: %w", err)
		}

		if len(record) < 3 {
			return fmt.Errorf("malformed CSV row: expected 3 columns, got %d", len(record))
		}

		name := record[0]
		code := record[1]
		countryCode := record[2]

		states = append(states, models.State{
			Name:        name,
			Code:        code,
			CountryCode: countryCode,
		})
	}

	if len(states) > 0 {
		if err := db.Create(&states).Error; err != nil {
			return fmt.Errorf("failed to batch insert states: %w", err)
		}
	}

	return nil
}
