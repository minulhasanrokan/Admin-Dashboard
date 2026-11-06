package repositories

import (
	"admin-dashboard/cmd/internal/dtos"
	"admin-dashboard/cmd/internal/models"
	"fmt"
)

func (r *Repository) GetActiveRightGroups() ([]models.RightGroup, error) {
	var groups []models.RightGroup

	if err := r.CacheSystem.Get("RightGroup", &groups); err == nil {
		return groups, nil
	}

	if err := r.Db.GetDb().Where("status = ?", 1).Order("sort_order ASC").Find(&groups).Error; err != nil {

		return nil, err
	}

	_ = r.CacheSystem.Set("RightGroup", groups)

	return groups, nil
}

func (r *Repository) GetActiveRightCategories() ([]models.RightCategory, error) {
	var categories []models.RightCategory

	if err := r.CacheSystem.Get("RightCategory", &categories); err == nil {
		return categories, nil
	}

	if err := r.Db.GetDb().Where("status = ?", 1).Order("sort_order ASC").Find(&categories).Error; err != nil {
		return nil, err
	}

	_ = r.CacheSystem.Set("RightCategory", categories)

	return categories, nil
}

func (r *Repository) GetActiveRightDetails() ([]models.RightDetail, error) {
	var details []models.RightDetail

	if err := r.CacheSystem.Get("RightDetail", &details); err == nil {
		return details, nil
	}

	if err := r.Db.GetDb().Where("status = ?", 1).Order("sort_order ASC").Find(&details).Error; err != nil {
		return nil, err
	}

	_ = r.CacheSystem.Set("RightDetail", details)

	return details, nil
}

func (r *Repository) GetUserRightKeys(userId uint64) ([]string, error) {

	var userRights []string

	cacheKey := fmt.Sprintf("UserRight%d", userId)

	if err := r.CacheSystem.Get(cacheKey, &userRights); err == nil {
		return userRights, nil
	}

	if err := r.Db.GetDb().Model(&models.UserRight{}).Where("user_id = ? AND status = ? AND delete_status = ?", userId, 1, 0).Pluck("right_key", &userRights).Error; err != nil {
		return nil, err
	}

	_ = r.CacheSystem.Set(cacheKey, userRights)

	return userRights, nil
}

func (r *Repository) CheckRightDetailAction(requiredAction string) (int64, error) {
	var count int64

	cacheKey := fmt.Sprintf("RightDetailAction%s", requiredAction)

	if err := r.CacheSystem.Get(cacheKey, &count); err == nil {
		return count, nil
	}

	if err := r.Db.GetDb().Model(&models.RightDetail{}).Where("action = ? AND status = ?", requiredAction, 1).Count(&count).Error; err != nil {
		return 0, err
	}

	_ = r.CacheSystem.Set(cacheKey, count)

	return count, nil
}

func (r *Repository) CheckCategoryRightDetailAction(categoryAction string) (int64, error) {
	var count int64

	cacheKey := fmt.Sprintf("RightCategoryDetailAction%s", categoryAction)

	if err := r.CacheSystem.Get(cacheKey, &count); err == nil {
		return count, nil
	}

	if err := r.Db.GetDb().Model(&models.RightDetail{}).Joins("JOIN right_categories rc ON rc.key = right_details.category_key").Where("rc.key = ?", categoryAction).Where("right_details.status = ?", 1).Where("rc.status = ?", 1).Count(&count).Error; err != nil {
		return 0, err
	}

	_ = r.CacheSystem.Set(cacheKey, count)

	return count, nil
}

func (r *Repository) CheckUserRightKey(userId uint64, requiredAction string) (int64, error) {
	var count int64

	cacheKey := fmt.Sprintf("UserRightCheck%d%s", userId, requiredAction)

	if err := r.CacheSystem.Get(cacheKey, &count); err == nil {
		return count, nil
	}

	if err := r.Db.GetDb().Model(&models.UserRight{}).Where("user_id = ? AND right_key = ?", userId, requiredAction).Where("status = ? AND delete_status = ?", 1, 0).Count(&count).Error; err != nil {
		return 0, err
	}

	_ = r.CacheSystem.Set(cacheKey, count)

	return count, nil
}

func (r *Repository) CheckUserCategoryRightKey(userId uint64, categoryAction string) (int64, error) {
	var count int64

	cacheKey := fmt.Sprintf("UserRightCategoryCheck%d%s", userId, categoryAction)

	if err := r.CacheSystem.Get(cacheKey, &count); err == nil {
		return count, nil
	}

	if err := r.Db.GetDb().
		Model(&models.UserRight{}).
		Joins("JOIN right_details rd ON rd.action = user_rights.right_key").
		Joins("JOIN right_categories rc ON rc.key = rd.category_key").
		Where("user_rights.user_id = ?", userId).
		Where("rc.key = ?", categoryAction).
		Where("user_rights.status = ?", 1).
		Where("user_rights.delete_status = ?", 0).
		Where("rd.status = ?", 1).
		Where("rc.status = ?", 1).
		Count(&count).Error; err != nil {
		return 0, err
	}

	_ = r.CacheSystem.Set(cacheKey, count)

	return count, nil
}

func (r *Repository) GetRightDetailsByKey(action string) (dtos.RightInfo, error) {
	var info dtos.RightInfo

	cacheKey := fmt.Sprintf("RightInfo%s", action)

	if err := r.CacheSystem.Get(cacheKey, &info); err == nil {
		return info, nil
	}

	if err := r.Db.GetDb().
		Table("right_details rd").
		Select(`
			rg.name         AS group_name,
			rg.key          AS group_key,
			rg.title        AS group_title,
			rg.details      AS group_details,
			rg.icon         AS group_icon,
			rg.sort_order  AS group_sort_order,
			rc.name         AS category_name,
			rc.key          AS category_key,
			rc.url          AS category_url,
			rc.title        AS category_title,
			rc.action       AS category_action,
			rc.details      AS category_details,
			rc.icon         AS category_icon,
			rc.sort_order  AS category_sort_order,
			rd.name         AS right_name,
			rd.action_name  AS right_action_name,
			rd.url          AS right_url,
			rd.title        AS right_title,
			rd.action       AS right_action,
			rd.details      AS right_details,
			rd.popup_status AS right_popup_status,
			rd.width        AS right_width,
			rd.icon         AS right_icon,
			rd.sort_order  AS right_sort_order ,
			rd.short_name  	AS right_short_name`).
		Joins("INNER JOIN right_categories rc ON rc.key = rd.category_key").
		Joins("INNER JOIN right_groups rg ON rg.key = rc.group_key").
		Where("rd.action = ?", action).
		Where("rd.status = ?", 1).
		Scan(&info).Error; err != nil {
		return dtos.RightInfo{}, err
	}

	_ = r.CacheSystem.Set(cacheKey, info)

	return info, nil
}

func (r *Repository) GetRightAllDetailsByKey(categoryKey string) ([]dtos.RightInfo, error) {
	var info []dtos.RightInfo

	cacheKey := fmt.Sprintf("RightAllInfo%s", categoryKey)

	if err := r.CacheSystem.Get(cacheKey, &info); err == nil {
		return info, nil
	}

	if err := r.Db.GetDb().
		Table("right_details rd").
		Select(`
			rg.name         AS group_name,
			rg.key          AS group_key,
			rg.title        AS group_title,
			rg.details      AS group_details,
			rg.icon         AS group_icon,
			rg.sort_order  AS group_sort_order,
			rc.name         AS category_name,
			rc.key          AS category_key,
			rc.url          AS category_url,
			rc.title        AS category_title,
			rc.action       AS category_action,
			rc.details      AS category_details,
			rc.icon         AS category_icon,
			rc.sort_order  AS category_sort_order,
			rd.name         AS right_name,
			rd.action_name  AS right_action_name,
			rd.url          AS right_url,
			rd.title        AS right_title,
			rd.action       AS right_action,
			rd.details      AS right_details,
			rd.popup_status AS right_popup_status,
			rd.width        AS right_width,
			rd.icon         AS right_icon,
			rd.sort_order  AS right_sort_order ,
			rd.short_name  	AS right_short_name`).
		Joins("INNER JOIN right_categories rc ON rc.key = rd.category_key").
		Joins("INNER JOIN right_groups rg ON rg.key = rc.group_key").
		Where("rd.category_key = ?", categoryKey).
		Where("rd.status = ?", 1).
		Where("rc.status = ?", 1).
		Where("rg.status = ?", 1).
		Order("rg.sort_order ASC, rc.sort_order ASC, rd.sort_order ASC").
		Find(&info).Error; err != nil {
		return nil, err
	}

	_ = r.CacheSystem.Set(cacheKey, info)

	return info, nil
}

func (r *Repository) GetRightDetailsByCategoryKey(categoryKey string) ([]dtos.RightInfo, error) {
	var info []dtos.RightInfo

	cacheKey := fmt.Sprintf("RightDetailsByCategory%s", categoryKey)

	if err := r.CacheSystem.Get(cacheKey, &info); err == nil {
		return info, nil
	}

	if err := r.Db.GetDb().
		Table("right_details rd").
		Select(`
			rg.name         AS group_name,
			rg.key          AS group_key,
			rg.title        AS group_title,
			rg.details      AS group_details,
			rg.icon         AS group_icon,
			rg.sort_order  AS group_sort_order,
			rc.name         AS category_name,
			rc.key          AS category_key,
			rc.url          AS category_url,
			rc.title        AS category_title,
			rc.action       AS category_action,
			rc.details      AS category_details,
			rc.icon         AS category_icon,
			rc.sort_order  AS category_sort_order,
			rd.name         AS right_name,
			rd.action_name  AS right_action_name,
			rd.url          AS right_url,
			rd.title        AS right_title,
			rd.action       AS right_action,
			rd.details      AS right_details,
			rd.popup_status AS right_popup_status,
			rd.width        AS right_width,
			rd.icon         AS right_icon,
			rd.sort_order  AS right_sort_order ,
			rd.short_name  	AS right_short_name`).
		Joins("INNER JOIN right_categories rc ON rc.key = rd.category_key").
		Joins("INNER JOIN right_groups rg ON rg.key = rc.group_key").
		Where("rd.category_key = ?", categoryKey).
		Where("rd.status = ?", 1).
		Where("rc.status = ?", 1).
		Where("rg.status = ?", 1).
		Order("rg.sort_order ASC, rc.sort_order ASC, rd.sort_order ASC").
		Find(&info).Error; err != nil {
		return nil, err
	}

	_ = r.CacheSystem.Set(cacheKey, info)

	return info, nil
}

func (r *Repository) GetRightDetailsByCategoryKeyUserId(userId uint64, categoryKey string) ([]dtos.RightInfo, error) {
	var info []dtos.RightInfo

	cacheKey := fmt.Sprintf("UserRightDetails%d%s", userId, categoryKey)

	if err := r.CacheSystem.Get(cacheKey, &info); err == nil {
		return info, nil
	}

	if err := r.Db.GetDb().
		Table("user_rights ur").
		Select(`
			rg.name         AS group_name,
			rg.key          AS group_key,
			rg.title        AS group_title,
			rg.details      AS group_details,
			rg.icon         AS group_icon,
			rg.sort_order  AS group_sort_order,
			rc.name         AS category_name,
			rc.key          AS category_key,
			rc.url          AS category_url,
			rc.title        AS category_title,
			rc.action       AS category_action,
			rc.details      AS category_details,
			rc.icon         AS category_icon,
			rc.sort_order  AS category_sort_order,
			rd.name         AS right_name,
			rd.action_name  AS right_action_name,
			rd.url          AS right_url,
			rd.title        AS right_title,
			rd.action       AS right_action,
			rd.details      AS right_details,
			rd.popup_status AS right_popup_status,
			rd.width        AS right_width,
			rd.icon         AS right_icon,
			rd.sort_order  AS right_sort_order ,
			rd.short_name  	AS right_short_name`).
		Joins("INNER JOIN right_details rd ON rd.action = ur.right_key").
		Joins("INNER JOIN right_categories rc ON rc.key = rd.category_key").
		Joins("INNER JOIN right_groups rg ON rg.key = rc.group_key").
		Where("ur.user_id = ?", userId).
		Where("rd.category_key = ?", categoryKey).
		Where("ur.status = ?", 1).
		Where("ur.delete_status = ?", 0).
		Where("rd.status = ?", 1).
		Where("rc.status = ?", 1).
		Where("rg.status = ?", 1).
		Order("rg.sort_order ASC, rc.sort_order ASC, rd.sort_order ASC").
		Find(&info).Error; err != nil {
		return nil, err
	}

	_ = r.CacheSystem.Set(cacheKey, info)

	return info, nil
}
