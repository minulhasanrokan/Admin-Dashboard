package middlewares

import (
	"admin-dashboard/cmd/internal/constants"
	"admin-dashboard/cmd/internal/session"
	"net/http"
	"time"

	"github.com/labstack/echo/v4"
)

func (m Middleware) LoginCheck() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			userSess, sess, err := session.GetUserSession(c)

			if err != nil || userSess.IsEmpty() || !userSess.IsActive {
				_ = session.ClearUserSession(c)
				return next(c)
			}

			if userSess.LastActiveTime == nil {
				_ = session.ClearUserSession(c)
				return next(c)
			}

			userRole, err := m.Service.UserRole()
			if err != nil {
				_ = session.ClearUserSession(c)
				return next(c)
			}

			idleTimeout := time.Duration(userRole.SessionTime) * time.Minute

			if time.Since(*userSess.LastActiveTime) > idleTimeout {
				_ = session.ClearUserSession(c)
				return next(c)
			}

			now := time.Now()
			userSess.LastActiveTime = &now

			sess.Values[constants.SessionKeyLastActiveTime] = userSess.LastActiveTime.Format(time.RFC3339)

			_ = sess.Save(c.Request(), c.Response())

			c.Set(constants.SessionName, userSess)

			if err := m.Service.UpdateUserLogTime(userSess.UserId, 0, 0); err != nil {
				_ = session.ClearUserSession(c)
				return next(c)
			}

			if err != nil {
				_ = session.ClearUserSession(c)
				return c.Redirect(http.StatusUnauthorized, "/login")
			}

			return c.Redirect(http.StatusSeeOther, "/dashboard")
		}
	}
}
