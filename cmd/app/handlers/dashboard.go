package handlers

import (
	"admin-dashboard/cmd/internal/constants"
	"admin-dashboard/cmd/internal/session"
	"encoding/json"

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

	dashData, _ := h.Service.GetDashboardData()

	ActivityList, err := json.Marshal(constants.ActivityList)

	if err != nil {
		return h.RenderWithErrorCheck(c, nil, nil, "500-error")
	}

	ProcessStatusList, err := json.Marshal(constants.ProcessStatusList)

	if err != nil {
		return h.RenderWithErrorCheck(c, nil, nil, "500-error")
	}

	DataActiveStatusList, err := json.Marshal(constants.DataActiveStatusList)

	if err != nil {
		return h.RenderWithErrorCheck(c, nil, nil, "500-error")
	}

	data := map[string]interface{}{
		"SessionData":          sessionData,
		"MenuTree":             menuTree,
		"Stats":                dashData.Stats,
		"RecentActivities":     dashData.RecentActivities,
		"RecentLogins":         dashData.RecentLogins,
		"ActivityList":         ActivityList,
		"ProcessStatusList":    ProcessStatusList,
		"DataActiveStatusList": DataActiveStatusList,
	}

	return h.RenderWithErrorCheck(c, data, err, "dashboard")
}
