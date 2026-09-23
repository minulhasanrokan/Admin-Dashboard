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

func (h Handler) DesignationPage(c echo.Context) error {
	userSess, _, err := session.GetUserSession(c)
	if err != nil || userSess.IsEmpty() || !userSess.IsActive {
		_ = session.ClearUserSession(c)
		return webResponse.SendUnauthorizedResponse(c, "Session expired", nil)
	}

	categoryKey := "reference.designation"
	if userSess.CategoryRouteAction != "" {
		categoryKey = userSess.CategoryRouteAction
	}

	var rightDetails []dtos.RightInfo
	if userSess.SuperAdmin == 1 {
		rightDetails, err = h.Service.GetRightDetailsByCategoryKeyKey(categoryKey)
	} else {
		rightDetails, err = h.Service.GetRightDetailsByCategoryKeyUserId(userSess.UserId, categoryKey)
	}
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	html, err := webResponse.RenderHTML(c, "reference/designation", map[string]any{
		"RightDetails": rightDetails,
		"RouteAction":  userSess.RouteAction,
	})
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Designations", html)
}

func (h Handler) DesignationList(c echo.Context) error {
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

	orderColumn := c.FormValue("order_column")
	if orderColumn == "" {
		orderColumn = c.QueryParam("order_column")
	}

	orderDir := c.FormValue("order_dir")
	if orderDir == "" {
		orderDir = c.QueryParam("order_dir")
	}

	filtersStr := c.FormValue("filters")
	if filtersStr == "" {
		filtersStr = c.QueryParam("filters")
	}

	var filters []repositories.ColumnFilter
	if filtersStr != "" {
		_ = json.Unmarshal([]byte(filtersStr), &filters)
	}

	params := repositories.BackupFilterParams{
		Page:        page,
		Limit:       limit,
		Search:      search,
		OrderColumn: orderColumn,
		OrderDir:    orderDir,
		Filters:     filters,
	}

	if params.Page < 1 {
		params.Page = 1
	}
	if params.Limit < 1 {
		params.Limit = constants.DefaultPageLimit
	}

	res, err := h.Service.GetDesignationList(params)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Designations retrieved successfully", res)
}

func (h Handler) DesignationAdd(c echo.Context) error {
	name := c.FormValue("name")
	title := c.FormValue("title")
	code := c.FormValue("code")
	details := c.FormValue("details")

	if name == "" {
		return webResponse.SendBadRequestResponse(c, "Designation name is required", nil)
	}

	err := h.Service.CreateDesignation(name, title, code, details)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Designation created successfully", nil)
}

func (h Handler) DesignationEdit(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return webResponse.SendBadRequestResponse(c, "Invalid designation ID", nil)
	}

	name := c.FormValue("name")
	title := c.FormValue("title")
	code := c.FormValue("code")
	details := c.FormValue("details")

	err = h.Service.UpdateDesignation(id, name, title, code, details)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Designation updated successfully", nil)
}

func (h Handler) DesignationDelete(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return webResponse.SendBadRequestResponse(c, "Invalid designation ID", nil)
	}

	err = h.Service.DeleteDesignation(id)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Designation deleted successfully", nil)
}

func (h Handler) DesignationActive(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return webResponse.SendBadRequestResponse(c, "Invalid designation ID", nil)
	}

	err = h.Service.SetDesignationStatus(id, 1)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Designation activated successfully", nil)
}

func (h Handler) DesignationDeactivate(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return webResponse.SendBadRequestResponse(c, "Invalid designation ID", nil)
	}

	err = h.Service.SetDesignationStatus(id, 0)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Designation deactivated successfully", nil)
}

func (h Handler) DesignationSTC(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return webResponse.SendBadRequestResponse(c, "Invalid designation ID", nil)
	}

	err = h.Service.SetDesignationChecker(id, 1)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Designation sent to checker", nil)
}

func (h Handler) DesignationVerify(c echo.Context) error {
	idStr := c.Param("id")
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return webResponse.SendBadRequestResponse(c, "Invalid designation ID", nil)
	}

	err = h.Service.SetDesignationChecker(id, 2)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Designation verified successfully", nil)
}
