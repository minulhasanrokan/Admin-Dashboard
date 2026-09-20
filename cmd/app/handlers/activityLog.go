package handlers

import (
	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/common/response/webResponse"
	"admin-dashboard/cmd/internal/constants"
	"admin-dashboard/cmd/internal/dtos"
	"admin-dashboard/cmd/internal/session"
	"encoding/json"
	"strconv"

	"github.com/labstack/echo/v4"
)

func (h Handler) ActivityLogPage(c echo.Context) error {
	userSess, _, err := session.GetUserSession(c)
	if err != nil || userSess.IsEmpty() || !userSess.IsActive || userSess.CategoryRouteAction == "" {
		_ = session.ClearUserSession(c)
		return webResponse.SendUnauthorizedResponse(c, "Session expired", nil)
	}

	var rightInfo dtos.RightInfo
	rightInfo, err = h.Service.GetRightDetailsByKey(userSess.RouteAction)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	var rightDetails []dtos.RightInfo
	rightDetails, err = h.Service.GetRightDetailsByCategoryKeyKey(rightInfo.CategoryKey)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	html, err := webResponse.RenderHTML(c, "settings/activity_logs", map[string]any{
		"RightDetails": rightDetails,
		"RouteAction":  userSess.RouteAction,
	})
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Activity Logs", html)
}

func (h Handler) ActivityLogList(c echo.Context) error {
	pageStr := c.FormValue("page")
	if pageStr == "" {
		pageStr = c.QueryParam("page")
	}
	page, _ := strconv.Atoi(pageStr)

	limitStr := c.FormValue("limit")
	if limitStr == "" {
		limitStr = c.QueryParam("limit")
	}
	limit, _ := strconv.Atoi(limitStr)

	search := c.FormValue("search")
	if search == "" {
		search = c.QueryParam("search")
	}

	searchField := c.FormValue("search_field")
	if searchField == "" {
		searchField = c.QueryParam("search_field")
	}

	operator := c.FormValue("operator")
	if operator == "" {
		operator = c.QueryParam("operator")
	}

	orderColumn := c.FormValue("order_column")
	if orderColumn == "" {
		orderColumn = c.QueryParam("order_column")
	}

	orderDir := c.FormValue("order_dir")
	if orderDir == "" {
		orderDir = c.QueryParam("order_dir")
	}

	var filters []repositories.ColumnFilter
	filtersJSON := c.FormValue("filters")
	if filtersJSON == "" {
		filtersJSON = c.QueryParam("filters")
	}
	if filtersJSON != "" {
		_ = json.Unmarshal([]byte(filtersJSON), &filters)
	}

	if limit <= 0 {
		limit = constants.DefaultPageLimit
	}

	params := repositories.BackupFilterParams{
		Page:        page,
		Limit:       limit,
		Search:      search,
		SearchField: searchField,
		Operator:    operator,
		Filters:     filters,
		OrderColumn: orderColumn,
		OrderDir:    orderDir,
	}

	result, err := h.Service.GetActivityLogList(params)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Activity logs retrieved successfully", result)
}

func (h Handler) LoginHistoryPage(c echo.Context) error {
	userSess, _, err := session.GetUserSession(c)
	if err != nil || userSess.IsEmpty() || !userSess.IsActive || userSess.CategoryRouteAction == "" {
		_ = session.ClearUserSession(c)
		return webResponse.SendUnauthorizedResponse(c, "Session expired", nil)
	}

	var rightInfo dtos.RightInfo
	rightInfo, err = h.Service.GetRightDetailsByKey(userSess.RouteAction)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	var rightDetails []dtos.RightInfo
	rightDetails, err = h.Service.GetRightDetailsByCategoryKeyKey(rightInfo.CategoryKey)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	html, err := webResponse.RenderHTML(c, "settings/login_history", map[string]any{
		"RightDetails": rightDetails,
		"RouteAction":  userSess.RouteAction,
	})
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Login History", html)
}

func (h Handler) LoginHistoryList(c echo.Context) error {
	pageStr := c.FormValue("page")
	if pageStr == "" {
		pageStr = c.QueryParam("page")
	}
	page, _ := strconv.Atoi(pageStr)

	limitStr := c.FormValue("limit")
	if limitStr == "" {
		limitStr = c.QueryParam("limit")
	}
	limit, _ := strconv.Atoi(limitStr)

	search := c.FormValue("search")
	if search == "" {
		search = c.QueryParam("search")
	}

	searchField := c.FormValue("search_field")
	if searchField == "" {
		searchField = c.QueryParam("search_field")
	}

	operator := c.FormValue("operator")
	if operator == "" {
		operator = c.QueryParam("operator")
	}

	orderColumn := c.FormValue("order_column")
	if orderColumn == "" {
		orderColumn = c.QueryParam("order_column")
	}

	orderDir := c.FormValue("order_dir")
	if orderDir == "" {
		orderDir = c.QueryParam("order_dir")
	}

	var filters []repositories.ColumnFilter
	filtersJSON := c.FormValue("filters")
	if filtersJSON == "" {
		filtersJSON = c.QueryParam("filters")
	}
	if filtersJSON != "" {
		_ = json.Unmarshal([]byte(filtersJSON), &filters)
	}

	if limit <= 0 {
		limit = constants.DefaultPageLimit
	}

	params := repositories.BackupFilterParams{
		Page:        page,
		Limit:       limit,
		Search:      search,
		SearchField: searchField,
		Operator:    operator,
		Filters:     filters,
		OrderColumn: orderColumn,
		OrderDir:    orderDir,
	}

	result, err := h.Service.GetLoginHistoryList(params)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Login history retrieved successfully", result)
}
