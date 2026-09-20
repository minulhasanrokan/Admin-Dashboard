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

func CountrySeed(db *gorm.DB) error {
	var count int64

	if err := db.Model(&models.Country{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing countries: %w", err)
	}

	if count > 0 {
		return nil
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to recover runtime caller metadata")
	}

	basePath := filepath.Dir(filename)
	csvPath := filepath.Join(basePath, "seed_data", "countries.csv")

	file, err := os.Open(csvPath)

	if err != nil {
		return fmt.Errorf("failed to open countries CSV: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("failed to read countries CSV header: %w", err)
	}

	var countries []models.Country

	for {

		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("error reading country CSV record: %w", err)
		}

		if len(record) < 7 {
			return fmt.Errorf("malformed CSV row: expected 7 columns, got %d", len(record))
		}

		countries = append(countries, models.Country{
			Name:          &record[0],
			Alpha2:        &record[1],
			Alpha3:        &record[2],
			Code:          &record[3],
			Iso31662:      &record[4],
			RegionCode:    &record[5],
			SubRegionCode: &record[6],
		})
	}

	if len(countries) > 0 {
		if err := db.Create(&countries).Error; err != nil {
			return fmt.Errorf("failed to batch insert countries: %w", err)
		}
	}

	return nil
}
