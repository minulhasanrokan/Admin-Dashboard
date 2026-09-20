package session

import (
	"admin-dashboard/cmd/common/helpers"
	"admin-dashboard/cmd/internal/config"
	"admin-dashboard/cmd/internal/constants"
	"admin-dashboard/cmd/internal/dtos"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/sessions"
	"github.com/labstack/echo-contrib/session"
	"github.com/labstack/echo/v4"
)

type UserSession struct {
	UserId               uint64
	UUID                 uuid.UUID
	Email                string
	UserName             string
	UserType             uint64
	Group                *uint64
	BranchID             *uint64
	Department           *uint64
	SuperAdmin           uint64
	IsActive             bool
	Verified             bool
	PasswordChangeStatus *uint64
	PasswordExpireDate   *time.Time
	PasswordChangeDate   *time.Time
	LastLoginIP          string
	LoginTime            *time.Time
	Photo                string
	LastActiveTime       *time.Time
	RouteAction          string
	RouteActionDetails   dtos.RightInfo
	CategoryRouteAction  string
}

func SetUserSession(c echo.Context, u UserSession, cfg config.Config) error {
	sess, err := session.Get(constants.SessionName, c)
	if err != nil {
		return err
	}

	sess.Values[constants.SessionKeyUserID] = u.UserId
	sess.Values[constants.SessionKeyUUID] = u.UUID.String()
	sess.Values[constants.SessionKeyEmail] = u.Email
	sess.Values[constants.SessionKeyUserName] = u.UserName
	sess.Values[constants.SessionKeyUserType] = u.UserType
	sess.Values[constants.SessionKeySuperAdmin] = u.SuperAdmin
	sess.Values[constants.SessionKeyIsActive] = u.IsActive
	sess.Values[constants.SessionKeyVerified] = u.Verified
	sess.Values[constants.SessionKeyLastLoginIP] = u.LastLoginIP
	sess.Values[constants.SessionKeyPhoto] = u.Photo

	if u.Group != nil {
		sess.Values[constants.SessionKeyGroup] = *u.Group
	}
	if u.BranchID != nil {
		sess.Values[constants.SessionKeyBranchID] = *u.BranchID
	}
	if u.Department != nil {
		sess.Values[constants.SessionKeyDepartment] = *u.Department
	}
	if u.PasswordChangeStatus != nil {
		sess.Values[constants.SessionKeyPasswordChangeStatus] = *u.PasswordChangeStatus
	}

	if u.PasswordExpireDate != nil {
		sess.Values[constants.SessionKeyPasswordExpireDate] = u.PasswordExpireDate.Format(time.RFC3339)
	}
	if u.PasswordChangeDate != nil {
		sess.Values[constants.SessionKeyPasswordChangeDate] = u.PasswordChangeDate.Format(time.RFC3339)
	}
	if u.LoginTime != nil {
		sess.Values[constants.SessionKeyLoginTime] = u.LoginTime.Format(time.RFC3339)
	}
	if u.LastActiveTime != nil {
		sess.Values[constants.SessionKeyLastActiveTime] = u.LastActiveTime.Format(time.RFC3339)
	}

	sess.Values[constants.SessionKeyRouteAction] = u.RouteAction
	sess.Values[constants.SessionKeyRouteActionDetails] = u.RouteActionDetails
	sess.Values[constants.SessionKeyCategoryRouteAction] = u.CategoryRouteAction

	sess.Options = &sessions.Options{
		Path:     cfg.WebCookiePath,
		Domain:   cfg.WebCookieDomain,
		MaxAge:   cfg.WebCookieMaxAge,
		HttpOnly: cfg.WebCookieHTTPOnly,
		Secure:   cfg.WebCookieSecure,
		SameSite: helpers.StrToSameSite(cfg.WebCookieSameSite),
	}

	return sess.Save(c.Request(), c.Response())
}

func GetUserSession(c echo.Context) (UserSession, *sessions.Session, error) {
	var u UserSession

	sess, err := session.Get(constants.SessionName, c)

	if err != nil {
		return u, nil, err
	}

	if val, ok := sess.Values[constants.SessionKeyUserID].(uint64); ok {
		u.UserId = val
	}

	if val, ok := sess.Values[constants.SessionKeyUUID].(string); ok {
		u.UUID, _ = uuid.Parse(val)
	}

	if val, ok := sess.Values[constants.SessionKeyEmail].(string); ok {
		u.Email = val
	}

	if val, ok := sess.Values[constants.SessionKeyUserName].(string); ok {
		u.UserName = val
	}

	if val, ok := sess.Values[constants.SessionKeyUserType].(uint64); ok {
		u.UserType = val
	}

	if val, ok := sess.Values[constants.SessionKeySuperAdmin].(uint64); ok {
		u.SuperAdmin = val
	}

	if val, ok := sess.Values[constants.SessionKeyIsActive].(bool); ok {
		u.IsActive = val
	}

	if val, ok := sess.Values[constants.SessionKeyVerified].(bool); ok {
		u.Verified = val
	}

	if val, ok := sess.Values[constants.SessionKeyLastLoginIP].(string); ok {
		u.LastLoginIP = val
	}

	if val, ok := sess.Values[constants.SessionKeyPhoto].(string); ok {
		u.Photo = val
	}

	if val, ok := sess.Values[constants.SessionKeyGroup].(uint64); ok {
		u.Group = &val
	}

	if val, ok := sess.Values[constants.SessionKeyBranchID].(uint64); ok {
		u.BranchID = &val
	}

	if val, ok := sess.Values[constants.SessionKeyDepartment].(uint64); ok {
		u.Department = &val
	}

	if val, ok := sess.Values[constants.SessionKeyPasswordChangeStatus].(uint64); ok {
		u.PasswordChangeStatus = &val
	}

	if val, ok := sess.Values[constants.SessionKeyPasswordExpireDate].(string); ok {
		if parsedTime, err := time.Parse(time.RFC3339, val); err == nil {
			u.PasswordExpireDate = &parsedTime
		}
	}

	if val, ok := sess.Values[constants.SessionKeyPasswordChangeDate].(string); ok {
		if parsedTime, err := time.Parse(time.RFC3339, val); err == nil {
			u.PasswordChangeDate = &parsedTime
		}
	}

	if val, ok := sess.Values[constants.SessionKeyLoginTime].(string); ok {
		if parsedTime, err := time.Parse(time.RFC3339, val); err == nil {
			u.LoginTime = &parsedTime
		}
	}

	if val, ok := sess.Values[constants.SessionKeyLastActiveTime].(string); ok {
		if parsedTime, err := time.Parse(time.RFC3339, val); err == nil {
			u.LastActiveTime = &parsedTime
		}
	}

	if val, ok := sess.Values[constants.SessionKeyRouteAction].(string); ok {
		u.RouteAction = val
	}

	if val, ok := sess.Values[constants.SessionKeyRouteActionDetails].(dtos.RightInfo); ok {
		u.RouteActionDetails = val
	}

	if val, ok := sess.Values[constants.SessionKeyCategoryRouteAction].(string); ok {
		u.CategoryRouteAction = val
	}

	return u, sess, nil
}

func (u UserSession) IsEmpty() bool {
	return u.UserId == 0
}

func ClearUserSession(c echo.Context) error {
	sess, err := session.Get(constants.SessionName, c)
	if err != nil {
		return err
	}
	sess.Values = make(map[interface{}]interface{})
	sess.Options.MaxAge = -1
	return sess.Save(c.Request(), c.Response())
}
