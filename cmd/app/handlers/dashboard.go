package handlers

import (
	"admin-dashboard/cmd/internal/session"

	"github.com/labstack/echo/v4"
)

func (h Handler) Dashboard(c echo.Context) error {

	sessionData, _, err := session.GetUserSession(c)

	if err != nil {
		return h.RenderWithErrorCheck(c, nil, err, "500-error")
	}

	if sessionData.IsEmpty() {

		return h.RenderWithErrorCheck(c, nil, err, "500-error")
	}

	isSuperAdmin := sessionData.SuperAdmin == 1
	menuTree, err := h.Service.GetUserDashboardRights(sessionData.UserId, isSuperAdmin)

	if err != nil {

		return h.RenderWithErrorCheck(c, nil, nil, "500-error")
	}

	data := map[string]interface{}{
		"SessionData": sessionData,
		"MenuTree":    menuTree,
	}

	return h.RenderWithErrorCheck(c, data, err, "dashboard")
}
