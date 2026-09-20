package middlewares

import (
	"admin-dashboard/cmd/common/response/webResponse"
	"admin-dashboard/cmd/internal/constants"
	"admin-dashboard/cmd/internal/session"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (m Middleware) RequireSession() echo.MiddlewareFunc {

	return func(next echo.HandlerFunc) echo.HandlerFunc {

		return func(c echo.Context) error {

			userSess, sess, err := session.GetUserSession(c)

			if err != nil || userSess.IsEmpty() || !userSess.IsActive {
				_ = session.ClearUserSession(c)
				return c.Redirect(http.StatusSeeOther, "/login")
			}

			userRole, err := m.Service.UserRole()

			if err != nil {
				_ = session.ClearUserSession(c)
				return c.Redirect(http.StatusSeeOther, "/login")
			}

			if userSess.LastActiveTime == nil {
				_ = session.ClearUserSession(c)
				return c.Redirect(http.StatusSeeOther, "/login")
			}

			idleTimeout := time.Duration(userRole.SessionTime) * time.Minute

			if time.Since(*userSess.LastActiveTime) > idleTimeout {
				_ = session.ClearUserSession(c)
				return c.Redirect(http.StatusSeeOther, "/login?reason=timeout")
			}

			now := time.Now()
			userSess.LastActiveTime = &now

			sess.Values[constants.SessionKeyLastActiveTime] = userSess.LastActiveTime.Format(time.RFC3339)

			_ = sess.Save(c.Request(), c.Response())

			c.Set(constants.SessionName, userSess)

			err = m.Service.UpdateUserLogTime(userSess.UserId, 0, 0)

			if err != nil {
				_ = session.ClearUserSession(c)
				return c.Redirect(http.StatusSeeOther, "/login")
			}

			return next(c)
		}
	}
}

func (m Middleware) RequireLoginSession() echo.MiddlewareFunc {

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

			now := time.Now()
			userSess.LastActiveTime = &now

			sess.Values[constants.SessionKeyLastActiveTime] = userSess.LastActiveTime.Format(time.RFC3339)

			_ = sess.Save(c.Request(), c.Response())

			c.Set(constants.SessionName, userSess)

			err = m.Service.UpdateUserLogTime(userSess.UserId, 0, 0)

			if err != nil {
				_ = session.ClearUserSession(c)
				return webResponse.SendUnauthorizedResponse(c, "Session expired", nil)
			}

			return next(c)
		}
	}
}
