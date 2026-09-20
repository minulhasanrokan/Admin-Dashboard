package seed

import (
	"encoding/csv"
	"os"
	"path/filepath"
	"runtime"
	"testing"
)

func TestSeedCSVFilesIntegrity(t *testing.T) {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatalf("failed to get caller path")
	}
	basePath := filepath.Dir(filename)
	dataDir := filepath.Join(basePath, "seed_data")

	expectedCSVs := []string{
		"activity_lists.csv",
		"backup_types.csv",
		"cities.csv",
		"countries.csv",
		"regions.csv",
		"right_categories.csv",
		"right_details.csv",
		"right_groups.csv",
		"states.csv",
		"sub_regions.csv",
	}

	for _, csvFile := range expectedCSVs {
		csvPath := filepath.Join(dataDir, csvFile)
		f, err := os.Open(csvPath)
		if err != nil {
			t.Errorf("missing or unreadable CSV file: %s, err: %v", csvFile, err)
			continue
		}

		reader := csv.NewReader(f)
		header, err := reader.Read()
		if err != nil {
			t.Errorf("failed to read header in %s: %v", csvFile, err)
			_ = f.Close()
			continue
		}

		if len(header) == 0 {
			t.Errorf("empty header in %s", csvFile)
		}

		// Verify at least 1 data row
		record, err := reader.Read()
		if err != nil {
			t.Errorf("no records found in %s: %v", csvFile, err)
		} else if len(record) != len(header) {
			t.Errorf("column count mismatch in %s: header has %d, row has %d", csvFile, len(header), len(record))
		}

		_ = f.Close()
	}
}

func TestBaseSeedData(t *testing.T) {
	base := BaseSeedData()
	if base.Status != 1 {
		t.Errorf("expected Status 1, got %d", base.Status)
	}
	if base.DeleteStatus != 0 {
		t.Errorf("expected DeleteStatus 0, got %d", base.DeleteStatus)
	}
	if base.ActiveStatus == 0 {
		t.Errorf("expected ActiveStatus 1, got %d", base.ActiveStatus)
	}
}
