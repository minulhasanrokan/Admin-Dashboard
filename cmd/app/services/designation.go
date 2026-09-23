package services

import (
	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/app/requests"
	"admin-dashboard/cmd/internal/models"
	"errors"
	"strings"
)

type DesignationRow struct {
	Id                uint64 `json:"id"`
	UUID              string `json:"uuid"`
	Name              string `json:"name"`
	DesignationName   string `json:"designation_name"`
	Title             string `json:"title"`
	DesignationTitle  string `json:"designation_title"`
	Code              string `json:"code"`
	DesignationCode   string `json:"designation_code"`
	Details           string `json:"details"`
	DesignationDetails string `json:"designation_details"`
	Status            int    `json:"status"`
	StatusText        string `json:"status_text"`
	ActualStatus      uint64 `json:"actual_status"`
	OpStatusText      string `json:"op_status_text"`
	ProcessStatus     uint64 `json:"process_status"`
	ProcessStatusText string `json:"process_status_text"`
	Message           string `json:"message"`
	Checker           uint64 `json:"checker"`
	CheckerText       string `json:"checker_text"`
	NotifyType        uint64 `json:"notify_type"`
	Remarks           string `json:"remarks"`
	CreatedAt         string `json:"created_at"`
	FormattedDate     string `json:"formatted_date"`
}

func (s *Service) GetDesignationList(params repositories.BackupFilterParams) (map[string]interface{}, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	list, total, err := repository.GetDesignations(params)
	if err != nil {
		return nil, err
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

		statusText := "Deactive"
		if d.Status == 1 {
			statusText = "Active"
		}

		opStatusText := "New Add"
		switch d.ActualStatus {
		case 1:
			opStatusText = "New Add"
		case 2:
			opStatusText = "Add-Edit"
		case 3:
			opStatusText = "Edit"
		case 4:
			opStatusText = "Delete"
		case 5:
			opStatusText = "Deactive"
		case 6:
			opStatusText = "Active"
		default:
			if d.ActualStatus > 0 {
				opStatusText = "Status " + string(rune('0'+d.ActualStatus))
			}
		}

		processStatusText := "Draft"
		switch d.ProcessStatus {
		case 0:
			processStatusText = "Draft"
		case 1:
			processStatusText = "Sent to Checker"
		case 2:
			processStatusText = "Returned"
		case 3:
			processStatusText = "Approved"
		}

		checkerText := "Pending"
		if d.Checker > 0 {
			checkerUser, _ := repository.GetUserByID(d.Checker)
			if checkerUser != nil && checkerUser.Name != nil {
				checkerText = *checkerUser.Name
			} else {
				checkerText = "Checker Assigned"
			}
		}

		dateStr := d.CreatedAt.Format("02 Jan, 2006 15:04:05")

		rows = append(rows, DesignationRow{
			Id:                 d.Id,
			UUID:               d.UUID.String(),
			Name:               name,
			DesignationName:    name,
			Title:              title,
			DesignationTitle:   title,
			Code:               code,
			DesignationCode:    code,
			Details:            details,
			DesignationDetails: details,
			Status:             d.Status,
			StatusText:         statusText,
			ActualStatus:       d.ActualStatus,
			OpStatusText:       opStatusText,
			ProcessStatus:      d.ProcessStatus,
			ProcessStatusText:  processStatusText,
			Message:            d.Message,
			Checker:            d.Checker,
			CheckerText:        checkerText,
			NotifyType:         d.NotifyType,
			Remarks:            d.Remarks,
			CreatedAt:          dateStr,
			FormattedDate:      dateStr,
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

func (s *Service) CheckDuplicateDesignationName(name string, excludeId uint64) (bool, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	return repository.CheckDuplicateDesignationName(name, excludeId)
}

func (s *Service) CheckDuplicateDesignationCode(code string, excludeId uint64) (bool, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	return repository.CheckDuplicateDesignationCode(code, excludeId)
}

func (s *Service) CreateDesignation(name, title, code, details string) (*models.Designation, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return nil, errors.New("designation name is required")
	}
	title = strings.TrimSpace(title)
	if title == "" {
		title = name
	}
	code = strings.TrimSpace(code)
	if code == "" {
		return nil, errors.New("designation code is required")
	}
	details = strings.TrimSpace(details)

	d := &models.Designation{
		Name:    &name,
		Title:   &title,
		Code:    &code,
		Details: &details,
	}

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	if err := repository.CreateDesignation(d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *Service) UpdateDesignation(identifier requests.UUIDRequest, name, title, code, details string) (*models.Designation, string, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	return repository.UpdateDesignationData(identifier, name, title, code, details)
}

func (s *Service) DeleteDesignationWithReason(identifier requests.UUIDRequest, reason string) (*models.Designation, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	return repository.DeleteDesignationWithReason(identifier, reason)
}

func (s *Service) SentToCheckerDesignation(identifier requests.UUIDRequest, checker uint64, notifyType uint64, remarks string) (*models.Designation, error) {
	if notifyType == 0 {
		authSetting, _ := s.AuthenticationSetting()
		if authSetting != nil && strings.ToLower(authSetting.OtpType) == "sms" {
			notifyType = 2
		} else {
			notifyType = 1
		}
	}
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	return repository.SentToCheckerDesignation(identifier, checker, notifyType, remarks)
}

func (s *Service) VerifyDesignation(identifier requests.UUIDRequest, verifyType int, returnReason string) (*models.Designation, string, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	return repository.VerifyDesignation(identifier, verifyType, returnReason)
}

func (s *Service) DeactiveDesignation(identifier requests.UUIDRequest, reason string) (*models.Designation, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	return repository.DeactiveDesignation(identifier, reason)
}

func (s *Service) ActiveDesignation(identifier requests.UUIDRequest, reason string) (*models.Designation, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	return repository.ActiveDesignation(identifier, reason)
}

func (s *Service) GetDesignationByIdOrUUID(identifier requests.UUIDRequest) (*models.Designation, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	return repository.GetDesignationByIdOrUUID(identifier)
}

func (s *Service) GetDesignationHistory(identifier requests.UUIDRequest) ([]models.DesignationHistory, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	d, err := repository.GetDesignationByIdOrUUID(identifier)
	if err != nil || d == nil {
		return repository.GetDesignationHistories(0)
	}
	return repository.GetDesignationHistories(d.Id)
}

func (s *Service) GetCheckerUserList() ([]map[string]interface{}, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	users, err := repository.GetCheckerUserList()
	if err != nil {
		return nil, err
	}
	result := make([]map[string]interface{}, 0, len(users))
	for _, u := range users {
		name := ""
		if u.Name != nil {
			name = *u.Name
		}
		result = append(result, map[string]interface{}{
			"id":   u.Id,
			"name": name,
		})
	}
	return result, nil
}
