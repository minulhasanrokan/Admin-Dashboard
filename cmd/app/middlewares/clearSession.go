package middlewares

import (
	"admin-dashboard/cmd/internal/session"

	"github.com/labstack/echo/v4"
)

func (m Middleware) ClearSession() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {

			_ = session.ClearUserSession(c)

			return next(c)
		}
	}
}
