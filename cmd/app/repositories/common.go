package repositories

import (
	"admin-dashboard/cmd/app/requests"
	"fmt"
	"strings"

	"gorm.io/gorm"
)

func (r *Repository) GetDynamicDropdown(req requests.GetDropdownRequest) ([]map[string]interface{}, error) {
	var result []map[string]interface{}

	cacheKey := fmt.Sprintf("Dropdown%s%s%s%s%s%s%s", req.TableName, req.FieldName, req.Value, req.OtherField, req.OtherValue, req.NotEqualField, req.NotEqualValue)

	if err := r.CacheSystem.Get(cacheKey, &result); err == nil {
		return result, nil
	}

	tableName := strings.TrimSpace(req.TableName)
	fieldName := strings.TrimSpace(req.FieldName)

	db := r.Db.GetDb().Table(tableName)

	if fieldName != "" && fieldName != "0" {
		fields := strings.Split(fieldName, ",")
		selectFields := make([]string, 0, len(fields))
		for _, f := range fields {
			selectFields = append(selectFields, strings.TrimSpace(f))
		}
		db = db.Select(selectFields)
	} else {
		db = db.Select("*")
	}

	if req.OtherField != "" && req.OtherField != "0" {
		otherFields := strings.Split(req.OtherField, "****")
		otherValues := strings.Split(req.OtherValue, "****")

		if len(otherFields) == 1 {
			values := strings.Split(req.OtherValue, ",")
			if strings.Contains(strings.ToUpper(req.OtherField), "ASSIGN") {
				db = db.Where(func(tx *gorm.DB) *gorm.DB {
					for _, v := range values {
						if v != "" {
							tx = tx.Or("FIND_IN_SET(?, "+req.OtherField+")", v)
						}
					}
					return tx
				})
			} else {
				db = db.Where(req.OtherField+" IN ?", values)
			}
		} else {
			for i, field := range otherFields {
				if i < len(otherValues) {
					db = db.Where(field+" = ?", otherValues[i])
				}
			}
		}
	}

	if req.NotEqualField != "" && req.NotEqualField != "0" {
		notFields := strings.Split(req.NotEqualField, "****")
		notValues := strings.Split(req.NotEqualValue, "****")

		if len(notFields) == 1 {
			values := strings.Split(req.NotEqualValue, ",")
			db = db.Where(req.NotEqualField+" NOT IN ?", values)
		} else {
			for i, field := range notFields {
				if i < len(notValues) {
					db = db.Where(field+" <> ?", notValues[i])
				}
			}
		}
	}

	if req.Value != "" && req.Value != "0" {
		values := strings.Split(req.Value, ",")
		db = db.Where("(delete_status = ? AND status = ?) OR id IN ?", 0, 1, values)
	} else {
		db = db.Where("delete_status = ? AND status = ?", 0, 1)
	}

	if err := db.Find(&result).Error; err != nil {
		return nil, err
	}

	_ = r.CacheSystem.Set(cacheKey, result)

	return result, nil
}
