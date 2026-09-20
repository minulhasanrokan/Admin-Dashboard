package repositories

import (
	"admin-dashboard/cmd/internal/models"
	"strings"
)

func (r *Repository) GetActivityLogs(params BackupFilterParams) ([]models.UserActivity, int64, error) {
	var list []models.UserActivity
	var total int64

	if r.Db == nil || r.Db.GetDb() == nil {
		return list, 0, nil
	}

	db := r.Db.GetDb().Model(&models.UserActivity{})

	if params.Search != "" {
		s := "%" + strings.TrimSpace(params.Search) + "%"
		db = db.Where("activity LIKE ? OR table_name LIKE ? OR ip_address LIKE ? OR device LIKE ?", s, s, s, s)
	}

	validColumns := map[string]string{
		"activity":        "activity",
		"table_name":      "table_name",
		"activity_status": "activity_status",
		"ip_address":      "ip_address",
		"device":          "device",
		"created_at":      "created_at",
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
		switch op {
		case "equals", "equal", "eq":
			db = db.Where(col+" = ?", val)
		case "not_equals", "not_equal", "ne":
			db = db.Where(col+" != ?", val)
		default:
			db = db.Where(col+" LIKE ?", "%"+val+"%")
		}
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderCol := "created_at"
	orderDir := "DESC"

	if params.OrderColumn != "" {
		if c, ok := validColumns[strings.ToLower(params.OrderColumn)]; ok {
			orderCol = c
		}
	}
	if strings.ToLower(params.OrderDir) == "asc" {
		orderDir = "ASC"
	}
	db = db.Order(orderCol + " " + orderDir)

	if params.Limit > 0 {
		page := params.Page
		if page < 1 {
			page = 1
		}
		offset := (page - 1) * params.Limit
		db = db.Offset(offset).Limit(params.Limit)
	}

	if err := db.Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *Repository) GetLoginHistories(params BackupFilterParams) ([]models.UserLoginHistory, int64, error) {
	var list []models.UserLoginHistory
	var total int64

	if r.Db == nil || r.Db.GetDb() == nil {
		return list, 0, nil
	}

	db := r.Db.GetDb().Model(&models.UserLoginHistory{})

	if params.Search != "" {
		s := "%" + strings.TrimSpace(params.Search) + "%"
		db = db.Where("ip_address LIKE ? OR user_agent LIKE ?", s, s)
	}

	validColumns := map[string]string{
		"ip_address":     "ip_address",
		"user_agent":     "user_agent",
		"log_out_status": "log_out_status",
		"created_at":     "created_at",
		"log_in_time":    "log_in_time",
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
		switch op {
		case "equals", "equal", "eq":
			db = db.Where(col+" = ?", val)
		case "not_equals", "not_equal", "ne":
			db = db.Where(col+" != ?", val)
		default:
			db = db.Where(col+" LIKE ?", "%"+val+"%")
		}
	}

	if err := db.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	orderCol := "created_at"
	orderDir := "DESC"

	if params.OrderColumn != "" {
		if c, ok := validColumns[strings.ToLower(params.OrderColumn)]; ok {
			orderCol = c
		}
	}
	if strings.ToLower(params.OrderDir) == "asc" {
		orderDir = "ASC"
	}
	db = db.Order(orderCol + " " + orderDir)

	if params.Limit > 0 {
		page := params.Page
		if page < 1 {
			page = 1
		}
		offset := (page - 1) * params.Limit
		db = db.Offset(offset).Limit(params.Limit)
	}

	if err := db.Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}
