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

func RightCategorySeed(db *gorm.DB) error {

	var count int64

	if err := db.Model(&models.RightCategory{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing right categories: %w", err)
	}

	if count > 0 {
		return nil
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to recover runtime caller metadata")
	}

	basePath := filepath.Dir(filename)
	csvPath := filepath.Join(basePath, "seed_data", "right_categories.csv")

	file, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("failed to open right_categories CSV: %w", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("failed to read right_categories CSV header: %w", err)
	}

	var categories []models.RightCategory

	for {

		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("error reading right_categories CSV: %w", err)
		}

		if len(record) < 10 {
			return fmt.Errorf("malformed CSV row: expected 10 columns, got %d", len(record))
		}

		sortOrder, err := strconv.Atoi(record[8])
		if err != nil {
			return fmt.Errorf("invalid sort_order: %w", err)
		}

		status, err := strconv.Atoi(record[9])
		if err != nil {
			return fmt.Errorf("invalid status: %w", err)
		}

		categories = append(categories, models.RightCategory{
			GroupKey:   record[0],
			Name:       record[1],
			Key:        record[2],
			Url:        record[3],
			Title:      record[4],
			Action:     record[5],
			Details:    record[6],
			Icon:       record[7],
			SortOrder: sortOrder,
			Status:     status,
		})
	}

	if len(categories) > 0 {
		if err := db.Create(&categories).Error; err != nil {
			return fmt.Errorf("failed to insert right categories: %w", err)
		}
	}

	return nil
}
