package repositories

import (
	"admin-dashboard/cmd/app/requests"
	"admin-dashboard/cmd/internal/models"
	"errors"
	"strings"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
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
		db = db.Where("name LIKE ? OR title LIKE ? OR code LIKE ? OR details LIKE ? OR message LIKE ?", s, s, s, s, s)
	}

	validColumns := map[string]string{
		"id":             "id",
		"name":           "name",
		"title":          "title",
		"code":           "code",
		"status":         "status",
		"actual_status":  "actual_status",
		"process_status": "process_status",
		"message":        "message",
		"created_at":     "created_at",
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

func (r *Repository) GetDesignationByIdOrUUID(identifier requests.UUIDRequest) (*models.Designation, error) {
	if r.Db == nil || r.Db.GetDb() == nil {
		return nil, nil
	}
	var d models.Designation
	val := strings.TrimSpace(identifier.UUID)
	if val == "" {
		return nil, errors.New("empty identifier")
	}

	if err := r.Db.GetDb().Where("(uuid = ? OR id = ?) AND delete_status = ?", val, val, 0).First(&d).Error; err != nil {
		return nil, err
	}
	return &d, nil
}

func (r *Repository) CheckDuplicateDesignationName(name string, excludeId uint64) (bool, error) {
	if r.Db == nil || r.Db.GetDb() == nil {
		return false, nil
	}
	trimmed := strings.TrimSpace(name)
	if trimmed == "" {
		return false, nil
	}

	var count int64
	q := r.Db.GetDb().Model(&models.Designation{}).Where("delete_status = 0 AND LOWER(TRIM(name)) = LOWER(?)", trimmed)
	if excludeId > 0 {
		q = q.Where("id != ?", excludeId)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) CheckDuplicateDesignationCode(code string, excludeId uint64) (bool, error) {
	if r.Db == nil || r.Db.GetDb() == nil {
		return false, nil
	}
	trimmed := strings.TrimSpace(code)
	if trimmed == "" {
		return false, nil
	}

	var count int64
	q := r.Db.GetDb().Model(&models.Designation{}).Where("delete_status = 0 AND LOWER(TRIM(code)) = LOWER(?)", trimmed)
	if excludeId > 0 {
		q = q.Where("id != ?", excludeId)
	}
	if err := q.Count(&count).Error; err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *Repository) CreateDesignation(d *models.Designation) error {
	if r.Db == nil || r.Db.GetDb() == nil {
		return errors.New("database connection unavailable")
	}

	return r.Db.GetDb().Transaction(func(tx *gorm.DB) error {
		if d.UUID == uuid.Nil {
			d.UUID = uuid.New()
		}
		d.ProcessStatus = 0
		d.ActualStatus = 1
		d.Status = 0
		d.DeleteStatus = 0

		if err := tx.Create(d).Error; err != nil {
			return err
		}

		hist := models.DesignationHistory{
			Designation: *d,
			MainId:      d.Id,
		}
		hist.Id = 0
		return tx.Create(&hist).Error
	})
}

func (r *Repository) UpdateDesignationData(identifier requests.UUIDRequest, name, title, code, details string) (*models.Designation, string, error) {
	if r.Db == nil || r.Db.GetDb() == nil {
		return nil, "", errors.New("database connection unavailable")
	}

	existing, err := r.GetDesignationByIdOrUUID(identifier)
	if err != nil || existing == nil {
		return nil, "", errors.New("designation not found")
	}

	var des string
	err = r.Db.GetDb().Transaction(func(tx *gorm.DB) error {
		// Business rule matching Laravel DesignationController:
		// if (actual_status == 1 || actual_status == 2) && process_status != 3 => actual_status = 2 ("Add-Edit")
		// else => edit_status = 1, actual_status = 3 ("Edit")
		if (existing.ActualStatus == 1 || existing.ActualStatus == 2) && existing.ProcessStatus != 3 {
			existing.ActualStatus = 2
			des = "Add-Edit Designation Information"
		} else {
			existing.EditStatus = 1
			existing.ActualStatus = 3
			des = "Edit Designation Information"
		}
		existing.ProcessStatus = 0

		trimmedName := strings.TrimSpace(name)
		if trimmedName != "" {
			existing.Name = &trimmedName
		}
		trimmedTitle := strings.TrimSpace(title)
		if trimmedTitle != "" {
			existing.Title = &trimmedTitle
		}
		trimmedCode := strings.TrimSpace(code)
		if trimmedCode != "" {
			existing.Code = &trimmedCode
		}
		trimmedDetails := strings.TrimSpace(details)
		existing.Details = &trimmedDetails

		if err := tx.Save(existing).Error; err != nil {
			return err
		}

		hist := models.DesignationHistory{
			Designation: *existing,
			MainId:      existing.Id,
		}
		hist.Id = 0
		return tx.Create(&hist).Error
	})

	return existing, des, err
}

func (r *Repository) DeleteDesignationWithReason(identifier requests.UUIDRequest, reason string) (*models.Designation, error) {
	if r.Db == nil || r.Db.GetDb() == nil {
		return nil, errors.New("database connection unavailable")
	}

	existing, err := r.GetDesignationByIdOrUUID(identifier)
	if err != nil || existing == nil {
		return nil, errors.New("designation not found")
	}

	err = r.Db.GetDb().Transaction(func(tx *gorm.DB) error {
		// Business rule matching Laravel:
		// If (actual_status == 1 || actual_status == 2) && process_status != 3 => delete_status = 1 immediately
		// Else => actual_status = 4, process_status = 0, message = delete_reason (waiting for checker)
		if (existing.ActualStatus == 1 || existing.ActualStatus == 2) && existing.ProcessStatus != 3 {
			existing.DeleteStatus = 1
		}
		existing.ActualStatus = 4
		existing.ProcessStatus = 0
		existing.Message = strings.TrimSpace(reason)
		existing.DeletedAt = time.Now()

		if err := tx.Save(existing).Error; err != nil {
			return err
		}

		hist := models.DesignationHistory{
			Designation: *existing,
			MainId:      existing.Id,
		}
		hist.Id = 0
		return tx.Create(&hist).Error
	})

	return existing, err
}

func (r *Repository) SentToCheckerDesignation(identifier requests.UUIDRequest, checker uint64, notifyType uint64, remarks string) (*models.Designation, error) {
	if r.Db == nil || r.Db.GetDb() == nil {
		return nil, errors.New("database connection unavailable")
	}

	existing, err := r.GetDesignationByIdOrUUID(identifier)
	if err != nil || existing == nil {
		return nil, errors.New("designation not found")
	}

	err = r.Db.GetDb().Transaction(func(tx *gorm.DB) error {
		existing.ProcessStatus = 1
		existing.Checker = checker
		existing.NotifyType = notifyType
		existing.Remarks = strings.TrimSpace(remarks)

		if err := tx.Save(existing).Error; err != nil {
			return err
		}

		hist := models.DesignationHistory{
			Designation: *existing,
			MainId:      existing.Id,
		}
		hist.Id = 0
		return tx.Create(&hist).Error
	})

	return existing, err
}

func (r *Repository) VerifyDesignation(identifier requests.UUIDRequest, verifyType int, returnReason string) (*models.Designation, string, error) {
	if r.Db == nil || r.Db.GetDb() == nil {
		return nil, "", errors.New("database connection unavailable")
	}

	existing, err := r.GetDesignationByIdOrUUID(identifier)
	if err != nil || existing == nil {
		return nil, "", errors.New("designation not found")
	}

	var des string
	err = r.Db.GetDb().Transaction(func(tx *gorm.DB) error {
		if verifyType == 1 {
			// Approve
			existing.ProcessStatus = 3
			des = "Approve Designation Information"

			if existing.ActualStatus == 1 || existing.ActualStatus == 2 {
				existing.Status = 1
			} else if existing.ActualStatus == 4 {
				existing.DeleteStatus = 1
				existing.DeletedAt = time.Now()
			} else if existing.ActualStatus == 5 {
				existing.Status = 0
			} else if existing.ActualStatus == 6 {
				existing.Status = 1
			}
		} else {
			// Return / Reject
			existing.ProcessStatus = 2
			existing.Message = strings.TrimSpace(returnReason)
			des = "Return Designation Information"
		}

		if err := tx.Save(existing).Error; err != nil {
			return err
		}

		hist := models.DesignationHistory{
			Designation: *existing,
			MainId:      existing.Id,
		}
		hist.Id = 0
		return tx.Create(&hist).Error
	})

	return existing, des, err
}

func (r *Repository) DeactiveDesignation(identifier requests.UUIDRequest, reason string) (*models.Designation, error) {
	if r.Db == nil || r.Db.GetDb() == nil {
		return nil, errors.New("database connection unavailable")
	}

	existing, err := r.GetDesignationByIdOrUUID(identifier)
	if err != nil || existing == nil {
		return nil, errors.New("designation not found")
	}

	err = r.Db.GetDb().Transaction(func(tx *gorm.DB) error {
		existing.ActualStatus = 5
		existing.Status = 0
		existing.ProcessStatus = 0
		existing.Message = strings.TrimSpace(reason)

		if err := tx.Save(existing).Error; err != nil {
			return err
		}

		hist := models.DesignationHistory{
			Designation: *existing,
			MainId:      existing.Id,
		}
		hist.Id = 0
		return tx.Create(&hist).Error
	})

	return existing, err
}

func (r *Repository) ActiveDesignation(identifier requests.UUIDRequest, reason string) (*models.Designation, error) {
	if r.Db == nil || r.Db.GetDb() == nil {
		return nil, errors.New("database connection unavailable")
	}

	existing, err := r.GetDesignationByIdOrUUID(identifier)
	if err != nil || existing == nil {
		return nil, errors.New("designation not found")
	}

	err = r.Db.GetDb().Transaction(func(tx *gorm.DB) error {
		existing.ActualStatus = 6
		existing.Status = 1
		existing.ProcessStatus = 0
		existing.Message = strings.TrimSpace(reason)

		if err := tx.Save(existing).Error; err != nil {
			return err
		}

		hist := models.DesignationHistory{
			Designation: *existing,
			MainId:      existing.Id,
		}
		hist.Id = 0
		return tx.Create(&hist).Error
	})

	return existing, err
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

func (r *Repository) GetCheckerUserList() ([]models.User, error) {
	var users []models.User
	if r.Db == nil || r.Db.GetDb() == nil {
		return users, nil
	}
	err := r.Db.GetDb().Model(&models.User{}).
		Where("is_active = ?", true).
		Order("name ASC").
		Find(&users).Error
	return users, err
}
