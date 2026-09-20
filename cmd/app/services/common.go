package services

import (
	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/app/requests"
	"admin-dashboard/cmd/common/helpers"
	"admin-dashboard/cmd/internal/models"
	"admin-dashboard/cmd/internal/session"

	"github.com/labstack/echo/v4"
)

func (s *Service) GetDynamicDropdown(req requests.GetDropdownRequest) ([]map[string]interface{}, error) {
	repository := repositories.NewRepository(&s.Db, s.CacheSystem)
	return repository.GetDynamicDropdown(req)
}

func (s *Service) CreateUserActivity(c echo.Context, processStatus string, tableName string, tableId uint64, activityText string) error {

	userSess, _, err := session.GetUserSession(c)

	if err != nil {
		return err
	}

	if userSess.RouteAction == "" {
		return nil
	}

	rightData, err := s.GetRightDetailsByKey(userSess.RouteAction)

	if err != nil {
		return err
	}

	ip := c.RealIP()
	device := c.Request().UserAgent()
	userId := userSess.UserId

	activity := &models.UserActivity{
		GroupKey:       &rightData.GroupKey,
		CategoryKey:    &rightData.CategoryKey,
		RightKey:       &rightData.RightAction,
		ActivityStatus: helpers.Str(processStatus),
		TableName:      helpers.Str(tableName),
		TableId:        helpers.UInt(tableId),
		Activity:       helpers.Str(activityText),
		IPAddress:      helpers.Str(ip),
		MacAddress:     helpers.Str(""),
		Device:         helpers.Str(device),
		Status:         1,
		CreatedBy:      &userId,
	}

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	return repository.CreateUserActivity(activity)
}
