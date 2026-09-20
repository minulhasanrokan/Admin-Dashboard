package middlewares

import (
	"admin-dashboard/cmd/common/response/webResponse"
	"admin-dashboard/cmd/internal/constants"
	"admin-dashboard/cmd/internal/session"
	"time"

	"github.com/labstack/echo/v4"
)

func (m Middleware) CheckPermission(requiredAction string) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {

			userSess, sess, err := session.GetUserSession(c)

			if err != nil || userSess.IsEmpty() || !userSess.IsActive {
				_ = session.ClearUserSession(c)
				return webResponse.SendUnauthorizedResponse(c, "Session expired", nil)
			}

			userRole, err := m.Service.UserRole()

			if err != nil {
				_ = session.ClearUserSession(c)
				return webResponse.SendUnauthorizedResponse(c, "Session expired", nil)
			}

			if userSess.LastActiveTime == nil {
				_ = session.ClearUserSession(c)
				return webResponse.SendUnauthorizedResponse(c, "Session expired", nil)
			}

			idleTimeout := time.Duration(userRole.SessionTime) * time.Minute

			if time.Since(*userSess.LastActiveTime) > idleTimeout {
				_ = session.ClearUserSession(c)
				return webResponse.SendUnauthorizedResponse(c, "Session expired", nil)
			}

			err = m.Service.UpdateUserLogTime(userSess.UserId, 0, 0)

			if err != nil {
				_ = session.ClearUserSession(c)
				return webResponse.SendUnauthorizedResponse(c, "Session expired", nil)
			}

			rightInfo, _ := m.Service.GetRightDetailsByKey(requiredAction)

			now := time.Now()
			userSess.LastActiveTime = &now
			userSess.RouteAction = requiredAction
			userSess.RouteActionDetails = rightInfo

			sess.Values[constants.SessionKeyRouteAction] = userSess.RouteAction
			sess.Values[constants.SessionKeyRouteActionDetails] = userSess.RouteActionDetails
			sess.Values[constants.SessionKeyLastActiveTime] = userSess.LastActiveTime.Format(time.RFC3339)

			_ = sess.Save(c.Request(), c.Response())

			c.Set(constants.SessionName, userSess)

			if userSess.SuperAdmin == 1 {

				hasAccess, err := m.Service.HasPermissionSuperAdmin(requiredAction)

				if err != nil || !hasAccess {
					_ = session.ClearUserSession(c)
					return webResponse.SendUnauthorizedResponse(c, "Super admin does not have permission for this action", nil)
				}

				return next(c)
			}

			hasAccess, err := m.Service.HasPermission(userSess.UserId, requiredAction)

			if err != nil || !hasAccess {
				_ = session.ClearUserSession(c)
				return webResponse.SendUnauthorizedResponse(c, "You do not have permission to access this resource", nil)
			}

			err = m.Service.UpdateUserLogTime(userSess.UserId, 0, 0)

			if err != nil {
				_ = session.ClearUserSession(c)
				return webResponse.SendInternalServerErrorResponse(c, err)
			}

			return next(c)
		}
	}
}

func (m Middleware) CheckCategoryPermission(categoryAction string) echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {

			userSess, sess, err := session.GetUserSession(c)

			if err != nil || userSess.IsEmpty() || !userSess.IsActive {
				_ = session.ClearUserSession(c)
				return webResponse.SendUnauthorizedResponse(c, "Session expired", nil)
			}

			userRole, err := m.Service.UserRole()

			if err != nil {
				_ = session.ClearUserSession(c)
				return webResponse.SendUnauthorizedResponse(c, "Session expired", nil)
			}

			if userSess.LastActiveTime == nil {
				_ = session.ClearUserSession(c)
				return webResponse.SendUnauthorizedResponse(c, "Session expired", nil)
			}

			idleTimeout := time.Duration(userRole.SessionTime) * time.Minute

			if time.Since(*userSess.LastActiveTime) > idleTimeout {
				_ = session.ClearUserSession(c)
				return webResponse.SendUnauthorizedResponse(c, "Session expired", nil)
			}

			err = m.Service.UpdateUserLogTime(userSess.UserId, 0, 0)

			if err != nil {
				_ = session.ClearUserSession(c)
				return webResponse.SendUnauthorizedResponse(c, "Session expired", nil)
			}

			now := time.Now()
			userSess.LastActiveTime = &now
			userSess.CategoryRouteAction = categoryAction

			sess.Values[constants.SessionKeyCategoryRouteAction] = userSess.CategoryRouteAction
			sess.Values[constants.SessionKeyLastActiveTime] = userSess.LastActiveTime.Format(time.RFC3339)

			_ = sess.Save(c.Request(), c.Response())

			c.Set(constants.SessionName, userSess)

			err = m.Service.UpdateUserLogTime(userSess.UserId, 0, 0)

			if userSess.SuperAdmin == 1 {

				hasAccess, err := m.Service.HasCategoryPermissionSuperAdmin(categoryAction)

				if err != nil || !hasAccess {
					_ = session.ClearUserSession(c)
					return webResponse.SendUnauthorizedResponse(c, "Super admin does not have permission for this action", nil)
				}

				return next(c)
			}

			hasAccess, err := m.Service.HasCategoryPermission(userSess.UserId, categoryAction)

			if err != nil || !hasAccess {
				_ = session.ClearUserSession(c)
				return webResponse.SendUnauthorizedResponse(c, "You do not have permission to access this resource", nil)
			}

			if err != nil {
				_ = session.ClearUserSession(c)
				return webResponse.SendInternalServerErrorResponse(c, err)
			}

			return next(c)
		}
	}
}
