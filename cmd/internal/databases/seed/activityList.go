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

func ActivityList(db *gorm.DB) error {
	var count int64

	if err := db.Model(&models.ActivityList{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing activity lists: %w", err)
	}

	if count > 0 {
		return nil
	}

	_, filename, _, ok := runtime.Caller(0)

	if !ok {
		return fmt.Errorf("failed to recover runtime caller metadata")
	}

	basePath := filepath.Dir(filename)
	csvPath := filepath.Join(basePath, "seed_data", "activity_lists.csv")

	file, err := os.Open(csvPath)

	if err != nil {
		return fmt.Errorf("failed to open activity lists CSV: %w", err)
	}

	defer func() {
		_ = file.Close()
	}()

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("failed to read activity lists CSV header: %w", err)
	}

	var activityLists []models.ActivityList

	for {
		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("error reading activity lists CSV record: %w", err)
		}

		if len(record) < 2 {
			return fmt.Errorf("malformed CSV row: expected 2 columns, got %d", len(record))
		}

		name := record[0]
		code := record[1]

		activityLists = append(activityLists, models.ActivityList{
			Name: name,
			Code: code,
		})
	}

	if len(activityLists) > 0 {
		if err := db.Create(&activityLists).Error; err != nil {
			return fmt.Errorf("failed to batch insert activity lists: %w", err)
		}
	}

	return nil
}
