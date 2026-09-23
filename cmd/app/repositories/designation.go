package repositories

import (
	"admin-dashboard/cmd/internal/models"
	"strings"
	"time"

	"github.com/google/uuid"
)

func (r *Repository) GetDesignations(params BackupFilterParams) ([]models.Designation, int64, error) {
	var list []models.Designation
	var total int64

	if r.Db == nil || r.Db.GetDb() == nil {
		return list, 0, nil
	}

	db := r.Db.GetDb().Model(&models.Designation{}).Where("delete_status = ?", 0)

	if params.Search != "" {
		s := "%" + strings.TrimSpace(params.Search) + "%"
		db = db.Where("name LIKE ? OR title LIKE ? OR code LIKE ? OR details LIKE ?", s, s, s, s)
	}

	validColumns := map[string]string{
		"id":         "id",
		"name":       "name",
		"title":      "title",
		"code":       "code",
		"status":     "status",
		"created_at": "created_at",
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

	orderCol := "id"
	orderDir := "DESC"

	if params.OrderColumn != "" {
		if c, ok := validColumns[strings.ToLower(params.OrderColumn)]; ok {
			orderCol = c
		}
	}
	if strings.ToUpper(params.OrderDir) == "ASC" {
		orderDir = "ASC"
	}

	db = db.Order(orderCol + " " + orderDir)

	if params.Limit > 0 {
		offset := (params.Page - 1) * params.Limit
		if offset < 0 {
			offset = 0
		}
		db = db.Offset(offset).Limit(params.Limit)
	}

	if err := db.Find(&list).Error; err != nil {
		return nil, 0, err
	}

	return list, total, nil
}

func (r *Repository) GetDesignationById(id uint64) (*models.Designation, error) {
	if r.Db == nil || r.Db.GetDb() == nil {
		return nil, nil
	}
	var d models.Designation
	if err := r.Db.GetDb().Where("id = ? AND delete_status = ?", id, 0).First(&d).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *Repository) CreateDesignation(d *models.Designation) error {
	if r.Db == nil || r.Db.GetDb() == nil {
		return nil
	}
	if d.UUID == uuid.Nil {
		d.UUID = uuid.New()
	}
	return r.Db.GetDb().Create(d).Error
}

func (r *Repository) UpdateDesignation(d *models.Designation) error {
	if r.Db == nil || r.Db.GetDb() == nil {
		return nil
	}
	return r.Db.GetDb().Save(d).Error
}

func (r *Repository) DeleteDesignation(id uint64) error {
	if r.Db == nil || r.Db.GetDb() == nil {
		return nil
	}
	now := time.Now()
	return r.Db.GetDb().Model(&models.Designation{}).Where("id = ?", id).Updates(map[string]interface{}{
		"delete_status": 1,
		"deleted_at":    now,
	}).Error
}

func (r *Repository) UpdateDesignationStatus(id uint64, status int) error {
	if r.Db == nil || r.Db.GetDb() == nil {
		return nil
	}
	return r.Db.GetDb().Model(&models.Designation{}).Where("id = ?", id).Update("status", status).Error
}

func (r *Repository) UpdateDesignationChecker(id uint64, checker uint64) error {
	if r.Db == nil || r.Db.GetDb() == nil {
		return nil
	}
	return r.Db.GetDb().Model(&models.Designation{}).Where("id = ?", id).Update("checker", checker).Error
}
