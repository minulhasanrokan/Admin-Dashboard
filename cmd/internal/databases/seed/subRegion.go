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

func SubRegionSeed(db *gorm.DB) error {
	var count int64

	if err := db.Model(&models.SubRegion{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing sub-regions: %w", err)
	}

	if count > 0 {
		return nil
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to recover runtime caller metadata")
	}

	basePath := filepath.Dir(filename)
	csvPath := filepath.Join(basePath, "seed_data", "sub_regions.csv")

	file, err := os.Open(csvPath)

	if err != nil {
		return fmt.Errorf("failed to open sub-regions CSV: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("failed to read sub-regions CSV header: %w", err)
	}

	var subRegions []models.SubRegion

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("error reading sub-region CSV record: %w", err)
		}

		if len(record) < 3 {
			return fmt.Errorf("malformed CSV row: expected 3 columns, got %d", len(record))
		}

		name := record[0]
		code := record[1]
		parentRegionCode := record[2]

		subRegions = append(subRegions, models.SubRegion{
			Name:       name,
			Code:       code,
			RegionCode: parentRegionCode,
		})
	}

	if len(subRegions) > 0 {
		if err := db.Create(&subRegions).Error; err != nil {
			return fmt.Errorf("failed to batch insert sub-regions: %w", err)
		}
	}

	return nil
}
