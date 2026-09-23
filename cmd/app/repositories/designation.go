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

func (r *Repository) GetDesignationByUUID(uuidStr string) (*models.Designation, error) {
	if r.Db == nil || r.Db.GetDb() == nil {
		return nil, nil
	}
	var d models.Designation
	if err := r.Db.GetDb().Where("uuid = ? AND delete_status = ?", uuidStr, 0).First(&d).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *Repository) GetDesignationByIdOrUUID(identifier string) (*models.Designation, error) {
	if r.Db == nil || r.Db.GetDb() == nil {
		return nil, nil
	}
	var d models.Designation
	if err := r.Db.GetDb().Where("(uuid = ? OR id = ?) AND delete_status = ?", identifier, identifier, 0).First(&d).Error; err != nil {
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
	if err := r.Db.GetDb().Create(d).Error; err != nil {
		return err
	}
	// Log history
	hist := models.DesignationHistory{
		Designation: *d,
		MainId:      d.Id,
	}
	hist.Id = 0
	_ = r.Db.GetDb().Create(&hist)
	return nil
}

func (r *Repository) UpdateDesignation(d *models.Designation) error {
	if r.Db == nil || r.Db.GetDb() == nil {
		return nil
	}
	if err := r.Db.GetDb().Save(d).Error; err != nil {
		return err
	}
	// Log history
	hist := models.DesignationHistory{
		Designation: *d,
		MainId:      d.Id,
	}
	hist.Id = 0
	_ = r.Db.GetDb().Create(&hist)
	return nil
}

func (r *Repository) DeleteDesignation(identifier string) error {
	if r.Db == nil || r.Db.GetDb() == nil {
		return nil
	}
	d, err := r.GetDesignationByIdOrUUID(identifier)
	if err != nil || d == nil {
		return err
	}
	now := time.Now()
	if err := r.Db.GetDb().Model(&models.Designation{}).Where("id = ?", d.Id).Updates(map[string]interface{}{
		"delete_status": 1,
		"deleted_at":    now,
	}).Error; err != nil {
		return err
	}
	d.DeleteStatus = 1
	d.DeletedAt = now
	hist := models.DesignationHistory{
		Designation: *d,
		MainId:      d.Id,
	}
	hist.Id = 0
	_ = r.Db.GetDb().Create(&hist)
	return nil
}

func (r *Repository) UpdateDesignationStatus(identifier string, status int) error {
	if r.Db == nil || r.Db.GetDb() == nil {
		return nil
	}
	d, err := r.GetDesignationByIdOrUUID(identifier)
	if err != nil || d == nil {
		return err
	}
	if err := r.Db.GetDb().Model(&models.Designation{}).Where("id = ?", d.Id).Update("status", status).Error; err != nil {
		return err
	}
	d.Status = status
	hist := models.DesignationHistory{
		Designation: *d,
		MainId:      d.Id,
	}
	hist.Id = 0
	_ = r.Db.GetDb().Create(&hist)
	return nil
}

func (r *Repository) UpdateDesignationChecker(identifier string, checker uint64) error {
	if r.Db == nil || r.Db.GetDb() == nil {
		return nil
	}
	d, err := r.GetDesignationByIdOrUUID(identifier)
	if err != nil || d == nil {
		return err
	}
	if err := r.Db.GetDb().Model(&models.Designation{}).Where("id = ?", d.Id).Update("checker", checker).Error; err != nil {
		return err
	}
	d.Checker = checker
	hist := models.DesignationHistory{
		Designation: *d,
		MainId:      d.Id,
	}
	hist.Id = 0
	_ = r.Db.GetDb().Create(&hist)
	return nil
}

func (r *Repository) GetDesignationHistories(designationId uint64) ([]models.DesignationHistory, error) {
	var histories []models.DesignationHistory
	if r.Db == nil || r.Db.GetDb() == nil {
		return histories, nil
	}
	db := r.Db.GetDb().Model(&models.DesignationHistory{})
	if designationId > 0 {
		db = db.Where("main_id = ?", designationId)
	}
	err := db.Order("id DESC").Limit(100).Find(&histories).Error
	return histories, err
}
