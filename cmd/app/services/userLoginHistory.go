package services

import (
	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/common/helpers"
	"admin-dashboard/cmd/internal/models"
	"time"

	"github.com/labstack/echo/v4"
)

func (s *Service) InsertUserLogHistory(c echo.Context, userId uint64) error {

	var userLoginHistory = models.UserLoginHistory{}

	userLoginHistory.LogOutStatus = helpers.Int(1)
	userLoginHistory.ForceStatus = helpers.Int(1)

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	_, err := repository.UpdateUserLoginHistory(userId, userLoginHistory)

	if err != nil {
		return err
	}

	ipAddress := c.RealIP()
	userAgent := c.Request().UserAgent()
	now := time.Now()

	userLoginHistory.LogOutStatus = helpers.Int(0)
	userLoginHistory.ForceStatus = helpers.Int(0)
	userLoginHistory.UserId = helpers.Int(userId)
	userLoginHistory.IpAddress = &ipAddress
	userLoginHistory.UserAgent = &userAgent
	userLoginHistory.LogInTime = &now
	userLoginHistory.LastActiveTime = &now

	err = repository.CreateUserLoginHistory(&userLoginHistory)

	if err != nil {
		return err
	}
	return nil
}

func (s *Service) UpdateUserLogTime(userLogID uint64, logOutStatus int, forceStatus int) error {

	repository := repositories.NewRepository(&s.Db, s.CacheSystem)

	now := time.Now()

	userLoginHistory := models.UserLoginHistory{
		LastActiveTime: &now,
		LogOutTime:     &now,
		LogOutStatus:   &logOutStatus,
		ForceStatus:    &forceStatus,
	}

	_, err := repository.UpdateUserLoginHistory(userLogID, userLoginHistory)

	if err != nil {
		return err
	}

	return nil
}
