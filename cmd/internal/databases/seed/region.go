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

func RegionSeed(db *gorm.DB) error {

	var count int64

	if err := db.Model(&models.Region{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing regions: %w", err)
	}

	if count > 0 {
		return nil
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to recover runtime caller metadata")
	}

	basePath := filepath.Dir(filename)
	csvPath := filepath.Join(basePath, "seed_data", "regions.csv")

	file, err := os.Open(csvPath)

	if err != nil {
		return fmt.Errorf("failed to open regions CSV: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {

		return fmt.Errorf("failed to read regions CSV header: %w", err)
	}

	var regions []models.Region

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("error reading region CSV record: %w", err)
		}

		if len(record) < 2 {
			return fmt.Errorf("malformed CSV row: expected at least 2 columns, got %d", len(record))
		}

		regions = append(regions, models.Region{
			Name: record[0],
			Code: record[1],
		})
	}

	if len(regions) > 0 {
		if err := db.Create(&regions).Error; err != nil {
			return fmt.Errorf("failed to batch insert regions: %w", err)
		}
	}

	return nil
}
