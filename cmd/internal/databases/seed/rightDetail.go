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

func RightDetailSeed(db *gorm.DB) error {

	var count int64

	if err := db.Model(&models.RightDetail{}).Count(&count).Error; err != nil {
		return fmt.Errorf("failed to check existing right details: %w", err)
	}

	if count > 0 {
		return nil
	}

	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		return fmt.Errorf("failed to recover runtime caller metadata")
	}

	basePath := filepath.Dir(filename)
	csvPath := filepath.Join(basePath, "seed_data", "right_details.csv")

	file, err := os.Open(csvPath)
	if err != nil {
		return fmt.Errorf("failed to open right_details CSV: %w", err)
	}

	defer func(file *os.File) {
		err := file.Close()
		if err != nil {

		}
	}(file)

	reader := csv.NewReader(file)

	if _, err := reader.Read(); err != nil {
		return fmt.Errorf("failed to read right_details CSV header: %w", err)
	}

	var details []models.RightDetail

	for {

		record, err := reader.Read()

		if err == io.EOF {
			break
		}

		if err != nil {
			return fmt.Errorf("error reading right_details CSV: %w", err)
		}

		if len(record) < 11 {
			return fmt.Errorf("malformed CSV row: expected 11 columns, got %d", len(record))
		}

		popupStatus, err := strconv.Atoi(record[7])
		if err != nil {
			return fmt.Errorf("invalid popup_status: %w", err)
		}

		width, err := strconv.Atoi(record[8])
		if err != nil {
			return fmt.Errorf("invalid width: %w", err)
		}

		sortOrder, err := strconv.Atoi(record[11])
		if err != nil {
			return fmt.Errorf("invalid sort_order: %w", err)
		}

		status, err := strconv.Atoi(record[12])
		if err != nil {
			return fmt.Errorf("invalid status: %w", err)
		}

		details = append(details, models.RightDetail{
			CategoryKey: record[0],
			Name:        record[1],
			ActionName:  record[2],
			Url:         record[3],
			Title:       record[4],
			Action:      record[5],
			Details:     record[6],
			PopupStatus: popupStatus,
			Width:       width,
			Icon:        record[9],
			ShortName:   record[10],
			SortOrder:  sortOrder,
			Status:      status,
		})
	}

	if len(details) > 0 {
		if err := db.Create(&details).Error; err != nil {
			return fmt.Errorf("failed to insert right details: %w", err)
		}
	}

	return nil
}
