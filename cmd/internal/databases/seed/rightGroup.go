package seed

import (
	"admin-dashboard/cmd/internal/models"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strconv"

	"gorm.io/gorm"
)

func RightGroupSeed(db *gorm.DB) error {

	var count int64

	if err := db.Model(&models.RightGroup{}).Count(&count).Error; err != nil {

		return fmt.Errorf("failed to check existing right groups: %w", err)
	}

	if count > 0 {
		return nil
	}

	_, filename, _, ok := runtime.Caller(0)

	if !ok {

		return fmt.Errorf("failed to recover runtime caller metadata")
	}

	basePath := filepath.Dir(filename)
	csvPath := filepath.Join(basePath, "seed_data", "right_groups.csv")

	file, err := os.Open(csvPath)

	if err != nil {
		
		return fmt.Errorf("failed to open right_groups CSV: %w", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("failed to read right_groups CSV header: %w", err)
	}

	var groups []models.RightGroup

	for {

		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("error reading right_groups CSV: %w", err)
		}

		if len(record) < 7 {
			return fmt.Errorf("invalid CSV row, expected 7 columns")
		}

		shortOrder, err := strconv.Atoi(record[5])
		if err != nil {
			return fmt.Errorf("invalid short_order: %w", err)
		}

		status, err := strconv.Atoi(record[6])
		if err != nil {
			return fmt.Errorf("invalid status: %w", err)
		}

		groups = append(groups, models.RightGroup{
			Name:       record[0],
			Key:        record[1],
			Title:      record[2],
			Details:    record[3],
			Icon:       record[4],
			ShortOrder: shortOrder,
			Status:     status,
		})
	}

	if len(groups) > 0 {
		if err := db.Create(&groups).Error; err != nil {
			return fmt.Errorf("failed to insert right groups: %w", err)
		}
	}

	return nil
}
