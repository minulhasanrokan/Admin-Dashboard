package services

import (
	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/internal/models"
	"errors"
	"strings"
)

func (s *Service) GetDesignationList(params repositories.BackupFilterParams) (map[string]interface{}, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	list, total, err := repository.GetDesignations(params)
	if err != nil {
		return nil, err
	}

	type DesignationRow struct {
		Id            uint64 `json:"id"`
		UUID          string `json:"uuid"`
		Name          string `json:"name"`
		Title         string `json:"title"`
		Code          string `json:"code"`
		Details       string `json:"details"`
		Status        int    `json:"status"`
		StatusText    string `json:"status_text"`
		Checker       uint64 `json:"checker"`
		CheckerText   string `json:"checker_text"`
		CreatedAt     string `json:"created_at"`
		FormattedDate string `json:"formatted_date"`
	}

	rows := make([]DesignationRow, 0, len(list))
	for _, d := range list {
		name := ""
		if d.Name != nil {
			name = *d.Name
		}
		title := ""
		if d.Title != nil {
			title = *d.Title
		}
		code := ""
		if d.Code != nil {
			code = *d.Code
		}
		details := ""
		if d.Details != nil {
			details = *d.Details
		}

		statusText := "Inactive"
		if d.Status == 1 {
			statusText = "Active"
		}

		checkerText := "Pending"
		if d.Checker == 1 {
			checkerText = "Sent to Checker"
		} else if d.Checker == 2 {
			checkerText = "Verified"
		}

		dateStr := d.CreatedAt.Format("02 Jan, 2006 15:04:05")

		rows = append(rows, DesignationRow{
			Id:            d.Id,
			UUID:          d.UUID.String(),
			Name:          name,
			Title:         title,
			Code:          code,
			Details:       details,
			Status:        d.Status,
			StatusText:    statusText,
			Checker:       d.Checker,
			CheckerText:   checkerText,
			CreatedAt:     dateStr,
			FormattedDate: dateStr,
		})
	}

	result := map[string]interface{}{
		"items": rows,
		"data":  rows,
		"total": total,
		"page":  params.Page,
		"limit": params.Limit,
	}

	return result, nil
}

func (s *Service) GetDesignationById(id uint64) (*models.Designation, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	return repository.GetDesignationById(id)
}

func (s *Service) GetDesignationByIdOrUUID(identifier string) (*models.Designation, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	return repository.GetDesignationByIdOrUUID(identifier)
}

func (s *Service) CreateDesignation(name, title, code, details string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return errors.New("name is required")
	}
	title = strings.TrimSpace(title)
	if title == "" {
		title = name
	}
	code = strings.TrimSpace(code)

	d := &models.Designation{
		Name:    &name,
		Title:   &title,
		Code:    &code,
		Details: &details,
	}
	d.Status = 1

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	return repository.CreateDesignation(d)
}

func (s *Service) UpdateDesignation(identifier string, name, title, code, details string) error {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	existing, err := repository.GetDesignationByIdOrUUID(identifier)
	if err != nil {
		return err
	}
	if existing == nil {
		return errors.New("designation not found")
	}

	name = strings.TrimSpace(name)
	if name != "" {
		existing.Name = &name
	}
	title = strings.TrimSpace(title)
	if title != "" {
		existing.Title = &title
	}
	code = strings.TrimSpace(code)
	existing.Code = &code
	details = strings.TrimSpace(details)
	existing.Details = &details

	return repository.UpdateDesignation(existing)
}

func (s *Service) DeleteDesignation(identifier string) error {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	return repository.DeleteDesignation(identifier)
}

func (s *Service) SetDesignationStatus(identifier string, status int) error {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	return repository.UpdateDesignationStatus(identifier, status)
}

func (s *Service) SetDesignationChecker(identifier string, checker uint64) error {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	return repository.UpdateDesignationChecker(identifier, checker)
}

func (s *Service) GetDesignationHistory(identifier string) ([]models.DesignationHistory, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	d, err := repository.GetDesignationByIdOrUUID(identifier)
	if err != nil || d == nil {
		return repository.GetDesignationHistories(0)
	}
	return repository.GetDesignationHistories(d.Id)
}

