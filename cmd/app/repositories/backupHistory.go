package repositories

import (
	"admin-dashboard/cmd/common/helpers"
	"admin-dashboard/cmd/internal/constants"
	"admin-dashboard/cmd/internal/models"
	"errors"
	"strconv"
	"strings"
	"time"
)

type ColumnFilter struct {
	Field    string `json:"field"`
	Operator string `json:"operator"`
	Value    string `json:"value"`
}

type BackupFilterParams struct {
	Page        int
	Limit       int
	Search      string
	SearchField string
	Operator    string
	Filters     []ColumnFilter
	OrderColumn string
	OrderDir    string
}

func (r *Repository) BackupHistory() (models.BackupHistory, error) {
	var data models.BackupHistory

	if err := r.CacheSystem.Get("BackupHistory", &data); err == nil {
		return data, nil
	}

	if err := r.Db.GetDb().Model(&models.BackupHistory{}).First(&data).Error; err != nil {
		return models.BackupHistory{}, err
	}

	_ = r.CacheSystem.Set("BackupHistory", data)

	return data, nil
}

func (r *Repository) CreateBackupHistory(backupHistory *models.BackupHistory) (*models.BackupHistory, error) {

	if backupHistory == nil {
		return nil, errors.New("backup history is empty")
	}

	if err := r.Db.GetDb().Create(backupHistory).Error; err != nil {
		return nil, err
	}

	return backupHistory, nil
}

func (r *Repository) GetBackupHistories(params BackupFilterParams) ([]models.BackupHistory, int64, error) {

	var list []models.BackupHistory
	var total int64

	db := r.Db.GetDb().Model(&models.BackupHistory{})

	validColumns := map[string]string{
		"file_name":           "file_name",
		"backup_type":         "backup_type",
		"backup_status":       "backup_status",
		"download_status":     "download_status",
		"file_size":           "file_size",
		"file_size_formatted": "file_size",
		"created_at":          "created_at",
		"backup_time":         "created_at",
		"downloaded_by":       "downloaded_by",
		"download_date":       "downloaded_at",
		"downloaded_at":       "downloaded_at",
		"error_log":           "error_log",
	}

	for _, f := range params.Filters {
		val := strings.TrimSpace(f.Value)
		if val == "" {
			continue
		}
		col, ok := validColumns[strings.ToLower(f.Field)]
		if !ok {
			continue
		}
		op := strings.ToLower(strings.TrimSpace(f.Operator))

		if col == "backup_type" {
			valLower := strings.ToLower(val)
			if strings.Contains(valLower, "file") || valLower == "0002" {
				db = db.Where("backup_type IN ('0002', 'files', 'Files')")
			} else if strings.Contains(valLower, "data") || valLower == "0001" {
				db = db.Where("backup_type IN ('0001', 'database', 'Database')")
			} else {
				db = db.Where("backup_type LIKE ?", "%"+val+"%")
			}
			continue
		}

		if col == "file_size" {
			sizeBytes, ok := parseFileSizeToBytes(val)
			if ok {
				switch op {
				case "less_than", "lt":
					db = db.Where("file_size < ?", sizeBytes)
				case "less_than_or_equal", "lte":
					db = db.Where("file_size <= ?", sizeBytes)
				case "greater_than", "gt":
					db = db.Where("file_size > ?", sizeBytes)
				case "greater_than_or_equal", "gte":
					db = db.Where("file_size >= ?", sizeBytes)
				case "not_equals", "not_equal":
					db = db.Where("file_size != ?", sizeBytes)
				case "equals", "equal":
					if sizeBytes < 1024 {
						db = db.Where("file_size = ?", sizeBytes)
					} else {
						minBytes := uint64(float64(sizeBytes) * 0.85)
						maxBytes := uint64(float64(sizeBytes) * 1.15)
						db = db.Where("file_size >= ? AND file_size <= ?", minBytes, maxBytes)
					}
				case "contains":
					minBytes := uint64(float64(sizeBytes) * 0.80)
					maxBytes := uint64(float64(sizeBytes) * 1.20)
					db = db.Where("file_size >= ? AND file_size <= ?", minBytes, maxBytes)
				default:
					minBytes := uint64(float64(sizeBytes) * 0.85)
					maxBytes := uint64(float64(sizeBytes) * 1.15)
					db = db.Where("file_size >= ? AND file_size <= ?", minBytes, maxBytes)
				}
			} else {
				db = db.Where("file_size LIKE ?", "%"+val+"%")
			}
			continue
		}

		if col == "backup_status" {
			valLower := strings.ToLower(val)
			if strings.Contains(valLower, "comp") || strings.Contains(valLower, "succ") || valLower == "3" {
				db = db.Where("backup_status = ?", 3)
			} else if strings.Contains(valLower, "run") || valLower == "2" {
				db = db.Where("backup_status = ?", 2)
			} else if strings.Contains(valLower, "fail") || valLower == "4" {
				db = db.Where("backup_status = ?", 4)
			} else if strings.Contains(valLower, "pend") || valLower == "1" {
				db = db.Where("backup_status = ?", 1)
			} else {
				db = db.Where("backup_status = ?", val)
			}
			continue
		}

		if col == "download_status" {
			valLower := strings.ToLower(val)
			if strings.Contains(valLower, "not") || valLower == "0" {
				db = db.Where("download_status = 0 OR download_status IS NULL")
			} else if strings.Contains(valLower, "down") || valLower == "1" {
				db = db.Where("download_status = 1")
			} else {
				db = db.Where("download_status = ?", val)
			}
			continue
		}

		if col == "downloaded_by" {
			var userIDs []uint64
			_ = r.Db.GetDb().Model(&models.User{}).Where("name LIKE ? OR user_name LIKE ?", "%"+val+"%", "%"+val+"%").Pluck("id", &userIDs).Error
			if len(userIDs) > 0 {
				db = db.Where("downloaded_by IN (?)", userIDs)
			} else {
				db = db.Where("1 = 0")
			}
			continue
		}

		if col == "created_at" || col == "downloaded_at" {
			parsedDate, err := time.Parse("2006-01-02", val)
			if err != nil {
				parsedDate, err = time.Parse("02/01/2006", val)
			}
			if err != nil {
				parsedDate, err = helpers.ParseInputDate(val)
			}
			if err != nil {
				parsedDate, err = helpers.ParseInputDateTime(val)
			}
			if err != nil {
				parsedDate, err = time.Parse(constants.ViewDateFormat, val)
			}
			if err != nil {
				parsedDate, err = time.Parse(constants.ViewDateTimeFormat, val)
			}

			if err == nil {
				startOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 0, 0, 0, 0, time.Local)
				endOfDay := time.Date(parsedDate.Year(), parsedDate.Month(), parsedDate.Day(), 23, 59, 59, 999999999, time.Local)
				switch op {
				case "equals", "equal":
					db = db.Where(col+" >= ? AND "+col+" <= ?", startOfDay, endOfDay)
				case "not_equals", "not_equal":
					db = db.Where(col+" < ? OR "+col+" > ?", startOfDay, endOfDay)
				case "less_than", "lt":
					db = db.Where(col+" < ?", startOfDay)
				case "less_than_or_equal", "lte":
					db = db.Where(col+" <= ?", endOfDay)
				case "greater_than", "gt":
					db = db.Where(col+" > ?", endOfDay)
				case "greater_than_or_equal", "gte":
					db = db.Where(col+" >= ?", startOfDay)
				default:
					db = db.Where(col+" >= ? AND "+col+" <= ?", startOfDay, endOfDay)
				}
			} else {
				db = db.Where("DATE_FORMAT("+col+", '%d %b %Y') LIKE ? OR DATE_FORMAT("+col+", '%Y-%m-%d') LIKE ? OR "+col+" LIKE ?", "%"+val+"%", "%"+val+"%", "%"+val+"%")
			}
			continue
		}

		switch op {
		case "not_contains":
			db = db.Where(col+" NOT LIKE ?", "%"+val+"%")
		case "starts_with":
			db = db.Where(col+" LIKE ?", val+"%")
		case "ends_with":
			db = db.Where(col+" LIKE ?", "%"+val)
		case "equals", "equal":
			db = db.Where(col+" = ?", val)
		case "not_equals", "not_equal":
			db = db.Where(col+" != ?", val)
		case "less_than", "lt":
			db = db.Where(col+" < ?", val)
		case "less_than_or_equal", "lte":
			db = db.Where(col+" <= ?", val)
		case "greater_than", "gt":
			db = db.Where(col+" > ?", val)
		case "greater_than_or_equal", "gte":
			db = db.Where(col+" >= ?", val)
		case "empty", "null":
			db = db.Where(col + " IS NULL OR " + col + " = ''")
		case "not_empty", "not_null":
			db = db.Where(col + " IS NOT NULL AND " + col + " != ''")
		default:
			db = db.Where(col+" LIKE ?", "%"+val+"%")
		}
	}

	search := strings.TrimSpace(params.Search)
	if search != "" {
		var matchedUserIDs []uint64
		_ = r.Db.GetDb().Model(&models.User{}).Where("name LIKE ? OR user_name LIKE ?", "%"+search+"%", "%"+search+"%").Pluck("id", &matchedUserIDs).Error

		searchConditions := []string{
			"file_name LIKE ?",
			"backup_type LIKE ?",
			"error_log LIKE ?",
			"DATE_FORMAT(created_at, '%d %b %Y %h:%i %p') LIKE ?",
			"DATE_FORMAT(downloaded_at, '%d %b %Y %h:%i %p') LIKE ?",
		}
		searchArgs := []interface{}{
			"%" + search + "%",
			"%" + search + "%",
			"%" + search + "%",
			"%" + search + "%",
			"%" + search + "%",
		}

		if sizeBytes, ok := parseFileSizeToBytes(search); ok {
			minB := uint64(float64(sizeBytes) * 0.80)
			maxB := uint64(float64(sizeBytes) * 1.20)
			searchConditions = append(searchConditions, "(file_size >= ? AND file_size <= ?)")
			searchArgs = append(searchArgs, minB, maxB)
		}

		if len(matchedUserIDs) > 0 {
			searchConditions = append(searchConditions, "downloaded_by IN (?)")
			searchArgs = append(searchArgs, matchedUserIDs)
		}

		queryStr := strings.Join(searchConditions, " OR ")
		db = db.Where(queryStr, searchArgs...)
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	page := params.Page
	if page < 1 {
		page = 1
	}
	limit := params.Limit
	if limit < 1 || limit > constants.MaxPageLimit {
		limit = constants.DefaultPageLimit
	}

	offset := (page - 1) * limit

	orderColumn := "id"
	validOrderCols := map[string]string{
		"id":              "id",
		"backup_type":     "backup_type",
		"file_name":       "file_name",
		"file_size":       "file_size",
		"backup_status":   "backup_status",
		"download_status": "download_status",
		"downloaded_by":   "downloaded_by",
		"created_at":      "created_at",
		"backup_time":     "created_at",
		"download_date":   "downloaded_at",
		"downloaded_at":   "downloaded_at",
	}
	if col, ok := validOrderCols[strings.ToLower(params.OrderColumn)]; ok {
		orderColumn = col
	}

	orderDir := "DESC"
	if strings.ToUpper(params.OrderDir) == "ASC" {
		orderDir = "ASC"
	}

	orderClause := orderColumn + " " + orderDir

	if err := db.Order(orderClause).Offset(offset).Limit(limit).Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *Repository) GetBackupHistoryByID(id uint64) (*models.BackupHistory, error) {

	var data models.BackupHistory

	if err := r.Db.GetDb().Model(&models.BackupHistory{}).Where("id = ?", id).First(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *Repository) GetBackupHistoryByUUID(uuidStr string) (*models.BackupHistory, error) {

	var data models.BackupHistory

	if err := r.Db.GetDb().Model(&models.BackupHistory{}).Where("uuid = ?", uuidStr).First(&data).Error; err != nil {
		return nil, err
	}

	return &data, nil
}

func (r *Repository) DeleteBackupHistory(id uint64) error {

	return r.Db.GetDb().Where("id = ?", id).Delete(&models.BackupHistory{}).Error
}

func (r *Repository) UpdateBackupStatusFailed(historyID uint64, errorLog string) error {
	return r.Db.GetDb().Model(&models.BackupHistory{}).
		Where("id = ?", historyID).
		Updates(map[string]interface{}{
			"backup_status": constants.BackupStatusFailed,
			"error_log":     errorLog,
		}).Error
}

func (r *Repository) UpdateBackupStatusSuccess(historyID uint64, fileName string, fileSize uint64) error {
	updateData := map[string]interface{}{
		"file_name":     fileName,
		"backup_status": constants.BackupStatusSuccess,
		"error_log":     "",
	}
	if fileSize > 0 {
		updateData["file_size"] = fileSize
	}

	return r.Db.GetDb().Model(&models.BackupHistory{}).
		Where("id = ?", historyID).
		Updates(updateData).Error
}

func (r *Repository) MarkBackupDownloaded(historyID uint64, uuidStr string, downloadedBy uint64, downloadedAt time.Time) error {
	updateData := map[string]interface{}{
		"download_status": 1,
		"downloaded_at":   downloadedAt,
	}
	if downloadedBy > 0 {
		updateData["downloaded_by"] = downloadedBy
	}

	return r.Db.GetDb().Model(&models.BackupHistory{}).
		Where("id = ? OR uuid = ?", historyID, uuidStr).
		Updates(updateData).Error
}

func parseFileSizeToBytes(input string) (uint64, bool) {
	input = strings.TrimSpace(strings.ToUpper(input))
	if input == "" {
		return 0, false
	}
	input = strings.ReplaceAll(input, ",", "")

	var multiplier uint64 = 1
	valStr := input

	if strings.HasSuffix(input, "TB") || strings.HasSuffix(input, "T") {
		multiplier = 1024 * 1024 * 1024 * 1024
		valStr = strings.TrimSpace(strings.TrimSuffix(strings.TrimSuffix(input, "TB"), "T"))
	} else if strings.HasSuffix(input, "GB") || strings.HasSuffix(input, "G") {
		multiplier = 1024 * 1024 * 1024
		valStr = strings.TrimSpace(strings.TrimSuffix(strings.TrimSuffix(input, "GB"), "G"))
	} else if strings.HasSuffix(input, "MB") || strings.HasSuffix(input, "M") {
		multiplier = 1024 * 1024
		valStr = strings.TrimSpace(strings.TrimSuffix(strings.TrimSuffix(input, "MB"), "M"))
	} else if strings.HasSuffix(input, "KB") || strings.HasSuffix(input, "K") {
		multiplier = 1024
		valStr = strings.TrimSpace(strings.TrimSuffix(strings.TrimSuffix(input, "KB"), "K"))
	} else if strings.HasSuffix(input, "BYTES") || strings.HasSuffix(input, "BYTE") || strings.HasSuffix(input, "B") {
		multiplier = 1
		valStr = strings.TrimSpace(strings.TrimSuffix(strings.TrimSuffix(strings.TrimSuffix(input, "BYTES"), "BYTE"), "B"))
	} else {
		if f, err := strconv.ParseFloat(input, 64); err == nil && f > 0 {
			if f <= 1024 {
				multiplier = 1024 * 1024
			} else {
				multiplier = 1
			}
		}
	}

	if f, err := strconv.ParseFloat(valStr, 64); err == nil && f >= 0 {
		return uint64(f * float64(multiplier)), true
	}
	return 0, false
}
