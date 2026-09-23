package handlers

import (
	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/app/requests"
	"admin-dashboard/cmd/common/response/webResponse"
	"admin-dashboard/cmd/internal/constants"
	"admin-dashboard/cmd/internal/dtos"
	"admin-dashboard/cmd/internal/session"
	"encoding/json"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type GridRightItem struct {
	ActionName  string `json:"r_action_name"`
	SortName    string `json:"r_sort_name"`
	Title       string `json:"r_title"`
	RouteUrl    string `json:"r_route_name"`
	Icon        string `json:"r_icon"`
	PopupStatus int    `json:"popup_status"`
	Width       int    `json:"width"`
}

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

	var gridRight []GridRightItem
	for _, r := range rightDetails {
		act := strings.TrimSpace(r.RightAction)
		if strings.HasSuffix(act, ".add") || strings.EqualFold(r.RightName, "add") {
			continue
		}

		shortAction := "view"
		sortName := "View"
		icon := r.RightIcon

		if strings.HasSuffix(act, ".view") {
			shortAction = "view"
			sortName = "View"
			if icon == "" {
				icon = "ti ti-eye"
			}
		} else if strings.HasSuffix(act, ".edit") {
			shortAction = "edit"
			sortName = "Edit"
			if icon == "" {
				icon = "ti ti-edit"
			}
		} else if strings.HasSuffix(act, ".delete") {
			shortAction = "delete"
			sortName = "Del"
			if icon == "" {
				icon = "ti ti-trash"
			}
		} else if strings.HasSuffix(act, ".stc") {
			shortAction = "sent_to_checker"
			sortName = "STC"
			if icon == "" {
				icon = "ti ti-arrow-forward-up"
			}
		} else if strings.HasSuffix(act, ".verify") {
			shortAction = "verify"
			sortName = "Ver"
			if icon == "" {
				icon = "ti ti-check"
			}
		} else if strings.HasSuffix(act, ".active") {
			shortAction = "active"
			sortName = "Act"
			if icon == "" {
				icon = "ti ti-lock-open"
			}
		} else if strings.HasSuffix(act, ".deactivate") {
			shortAction = "deactive"
			sortName = "Deact"
			if icon == "" {
				icon = "ti ti-lock"
			}
		}

		popupStatus := r.RightPopupStatus
		if popupStatus == 0 {
			popupStatus = 1
		}
		width := r.RightWidth
		if width == 0 {
			width = 1000
		}

		gridRight = append(gridRight, GridRightItem{
			ActionName:  shortAction,
			SortName:    sortName,
			Title:       r.RightTitle,
			RouteUrl:    r.RightURL,
			Icon:        icon,
			PopupStatus: popupStatus,
			Width:       width,
		})
	}

	gridRightJSON, _ := json.Marshal(gridRight)
	checkerUsers, _ := h.Service.GetCheckerUserList()

	html, err := webResponse.RenderHTML(c, "reference/designation", map[string]any{
		"RightDetails":  rightDetails,
		"GridRight":     gridRight,
		"GridRightJSON": string(gridRightJSON),
		"RouteAction":   userSess.RouteAction,
		"CheckerUsers":  checkerUsers,
		"SuperAdmin":    userSess.SuperAdmin,
	})
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Designations", html)
}

func (h Handler) DesignationAddPage(c echo.Context) error {
	return h.DesignationPage(c)
}

func (h Handler) DesignationList(c echo.Context) error {
	drawStr := c.FormValue("draw")
	if drawStr == "" {
		drawStr = c.QueryParam("draw")
	}

	startStr := c.FormValue("start")
	if startStr == "" {
		startStr = c.QueryParam("start")
	}

	lengthStr := c.FormValue("length")
	if lengthStr == "" {
		lengthStr = c.QueryParam("length")
	}

	var page, limit int
	if startStr != "" && lengthStr != "" {
		start, _ := strconv.Atoi(startStr)
		length, _ := strconv.Atoi(lengthStr)
		if length <= 0 {
			length = 20
		}
		limit = length
		page = (start / length) + 1
	} else {
		pageStr := c.FormValue("page")
		if pageStr == "" {
			pageStr = c.QueryParam("page")
		}
		page, _ = strconv.Atoi(pageStr)

		limitStr := c.FormValue("limit")
		if limitStr == "" {
			limitStr = c.QueryParam("limit")
		}
		limit, _ = strconv.Atoi(limitStr)
	}

	search := c.FormValue("search[value]")
	if search == "" {
		search = c.FormValue("search")
	}
	if search == "" {
		search = c.QueryParam("search")
	}

	orderColumn := c.FormValue("order_column")
	if orderColumn == "" {
		orderColumn = c.QueryParam("order_column")
	}

	orderDir := c.FormValue("order_dir")
	if orderDir == "" {
		orderDir = c.FormValue("order[0][dir]")
	}
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

	if drawStr != "" {
		draw, _ := strconv.Atoi(drawStr)
		total := res["total"]
		data := res["data"]
		return c.JSON(200, map[string]interface{}{
			"draw":            draw,
			"recordsTotal":    total,
			"recordsFiltered": total,
			"data":            data,
			"items":           data,
			"success":         true,
		})
	}

	return webResponse.SendSuccessResponse(c, "Designations retrieved successfully", res)
}

func (h Handler) DesignationAdd(c echo.Context) error {
	req := new(requests.CreateDesignationRequest)
	_ = c.Bind(req)

	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return webResponse.SendBadRequestResponse(c, "Designation name is required", nil)
	}
	if len(name) > 150 {
		return webResponse.SendBadRequestResponse(c, "Designation name must not exceed 150 characters", nil)
	}

	code := strings.TrimSpace(req.GetCode())
	if code == "" {
		return webResponse.SendBadRequestResponse(c, "Designation code is required", nil)
	}
	if len(code) > 50 {
		return webResponse.SendBadRequestResponse(c, "Designation code must not exceed 50 characters", nil)
	}

	title := strings.TrimSpace(req.GetTitle())
	if title == "" {
		title = name
	}
	if len(title) > 150 {
		return webResponse.SendBadRequestResponse(c, "Designation title must not exceed 150 characters", nil)
	}

	details := strings.TrimSpace(req.GetDetails())

	// Check duplicate name
	isDupName, err := h.Service.CheckDuplicateDesignationName(name, 0)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}
	if isDupName {
		return webResponse.SendBadRequestResponse(c, "Duplicate Name Found", nil)
	}

	// Check duplicate code
	isDupCode, err := h.Service.CheckDuplicateDesignationCode(code, 0)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}
	if isDupCode {
		return webResponse.SendBadRequestResponse(c, "Duplicate Code Found", nil)
	}

	d, err := h.Service.CreateDesignation(name, title, code, details)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	_ = h.Service.CreateUserActivity(c, constants.ActivityAdd, "designations", d.Id, "Add Designation Information")

	return webResponse.SendSuccessResponse(c, "Saved Successfully", d)
}

func (h Handler) DesignationDetailsPage(c echo.Context) error {
	var uUId requests.UUIDRequest
	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	d, err := h.Service.GetDesignationByIdOrUUID(uUId)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}
	if d == nil {
		return webResponse.SendNotFoundResponse(c, "Designation not found", nil)
	}
	return webResponse.SendSuccessResponse(c, "Designation details", d)
}

func (h Handler) DesignationEditPage(c echo.Context) error {
	var uUId requests.UUIDRequest
	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	d, err := h.Service.GetDesignationByIdOrUUID(uUId)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}
	if d == nil {
		return webResponse.SendNotFoundResponse(c, "Designation not found", nil)
	}

	return webResponse.SendSuccessResponse(c, "Designation data", d)
}

func (h Handler) DesignationEdit(c echo.Context) error {
	var uUId requests.UUIDRequest
	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil {
		return webResponse.SendBadRequestResponse(c, "Invalid designation identifier", nil)
	}

	existing, err := h.Service.GetDesignationByIdOrUUID(uUId)
	if err != nil || existing == nil {
		return webResponse.SendNotFoundResponse(c, "Designation not found", nil)
	}

	req := new(requests.UpdateDesignationRequest)
	_ = c.Bind(req)

	name := strings.TrimSpace(req.GetName())
	if name == "" {
		return webResponse.SendBadRequestResponse(c, "Designation name is required", nil)
	}
	if len(name) > 150 {
		return webResponse.SendBadRequestResponse(c, "Designation name must not exceed 150 characters", nil)
	}

	code := strings.TrimSpace(req.GetCode())
	if code == "" {
		return webResponse.SendBadRequestResponse(c, "Designation code is required", nil)
	}
	if len(code) > 50 {
		return webResponse.SendBadRequestResponse(c, "Designation code must not exceed 50 characters", nil)
	}

	title := strings.TrimSpace(req.GetTitle())
	if title == "" {
		title = name
	}
	if len(title) > 150 {
		return webResponse.SendBadRequestResponse(c, "Designation title must not exceed 150 characters", nil)
	}

	details := strings.TrimSpace(req.GetDetails())

	// Check duplicate name excluding current record ID
	isDupName, err := h.Service.CheckDuplicateDesignationName(name, existing.Id)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}
	if isDupName {
		return webResponse.SendBadRequestResponse(c, "Duplicate Name Found", nil)
	}

	// Check duplicate code excluding current record ID
	isDupCode, err := h.Service.CheckDuplicateDesignationCode(code, existing.Id)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}
	if isDupCode {
		return webResponse.SendBadRequestResponse(c, "Duplicate Code Found", nil)
	}

	updated, des, err := h.Service.UpdateDesignation(uUId, name, title, code, details)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	_ = h.Service.CreateUserActivity(c, constants.ActivityEdit, "designations", updated.Id, des)

	return webResponse.SendSuccessResponse(c, "Updated Successfully", updated)
}

func (h Handler) DesignationDeletePage(c echo.Context) error {
	return h.DesignationDetailsPage(c)
}

func (h Handler) DesignationDelete(c echo.Context) error {
	var uUId requests.UUIDRequest
	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil {
		return webResponse.SendBadRequestResponse(c, "Invalid designation identifier", nil)
	}

	reason := strings.TrimSpace(c.FormValue("delete_reason"))
	if reason == "" {
		reason = strings.TrimSpace(c.FormValue("reason"))
	}
	if reason == "" {
		return webResponse.SendBadRequestResponse(c, "Delete reason is required", nil)
	}
	if len(reason) > 250 {
		return webResponse.SendBadRequestResponse(c, "Delete reason must not exceed 250 characters", nil)
	}

	d, err := h.Service.DeleteDesignationWithReason(uUId, reason)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	_ = h.Service.CreateUserActivity(c, constants.ActivityDelete, "designations", d.Id, "Delete Designation Information")

	return webResponse.SendSuccessResponse(c, "Deleted Successfully", nil)
}

func (h Handler) DesignationStcPage(c echo.Context) error {
	return h.DesignationDetailsPage(c)
}

func (h Handler) DesignationStc(c echo.Context) error {
	var uUId requests.UUIDRequest
	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil {
		return webResponse.SendBadRequestResponse(c, "Invalid designation identifier", nil)
	}

	checkerStr := c.FormValue("checker")
	checker, _ := strconv.ParseUint(checkerStr, 10, 64)
	if checker == 0 {
		return webResponse.SendBadRequestResponse(c, "Checker user is required", nil)
	}

	notifyTypeStr := c.FormValue("notify_type")
	notifyType, _ := strconv.ParseUint(notifyTypeStr, 10, 64)
	if notifyType == 0 {
		authSetting, _ := h.Service.AuthenticationSetting()
		if authSetting != nil && strings.ToLower(authSetting.OtpType) == "sms" {
			notifyType = 2
		} else {
			notifyType = 1
		}
	}

	remarks := strings.TrimSpace(c.FormValue("remarks"))
	if len(remarks) > 250 {
		return webResponse.SendBadRequestResponse(c, "Remarks must not exceed 250 characters", nil)
	}

	d, err := h.Service.SentToCheckerDesignation(uUId, checker, notifyType, remarks)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	_ = h.Service.CreateUserActivity(c, "1", "designations", d.Id, "Sent To Checker Designation Information")

	return webResponse.SendSuccessResponse(c, "Sent To Checker Successfully", nil)
}

func (h Handler) DesignationVerifyPage(c echo.Context) error {
	return h.DesignationDetailsPage(c)
}

func (h Handler) DesignationVerify(c echo.Context) error {
	var uUId requests.UUIDRequest
	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil {
		return webResponse.SendBadRequestResponse(c, "Invalid designation identifier", nil)
	}

	verifyTypeStr := c.FormValue("verify_type")
	verifyType, _ := strconv.Atoi(verifyTypeStr)

	returnReason := strings.TrimSpace(c.FormValue("return_reason"))
	if verifyType == 0 && returnReason == "" {
		return webResponse.SendBadRequestResponse(c, "Return reason is required", nil)
	}
	if len(returnReason) > 250 {
		return webResponse.SendBadRequestResponse(c, "Return reason must not exceed 250 characters", nil)
	}

	d, des, err := h.Service.VerifyDesignation(uUId, verifyType, returnReason)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	_ = h.Service.CreateUserActivity(c, strconv.Itoa(int(d.ActualStatus)), "designations", d.Id, des)

	if verifyType == 1 {
		return webResponse.SendSuccessResponse(c, "Approved Successfully", nil)
	}
	return webResponse.SendSuccessResponse(c, "Returned Successfully", nil)
}

func (h Handler) DesignationActivePage(c echo.Context) error {
	return h.DesignationDetailsPage(c)
}

func (h Handler) DesignationActive(c echo.Context) error {
	var uUId requests.UUIDRequest
	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil {
		return webResponse.SendBadRequestResponse(c, "Invalid designation identifier", nil)
	}

	reason := strings.TrimSpace(c.FormValue("active_reason"))
	if reason == "" {
		reason = strings.TrimSpace(c.FormValue("reason"))
	}
	if reason == "" {
		return webResponse.SendBadRequestResponse(c, "Active reason is required", nil)
	}
	if len(reason) > 250 {
		return webResponse.SendBadRequestResponse(c, "Active reason must not exceed 250 characters", nil)
	}

	d, err := h.Service.ActiveDesignation(uUId, reason)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	_ = h.Service.CreateUserActivity(c, constants.ActivityActivate, "designations", d.Id, "Active Designation Information")

	return webResponse.SendSuccessResponse(c, "Activated Successfully", nil)
}

func (h Handler) DesignationDeactivatePage(c echo.Context) error {
	return h.DesignationDetailsPage(c)
}

func (h Handler) DesignationDeactivate(c echo.Context) error {
	var uUId requests.UUIDRequest
	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil {
		return webResponse.SendBadRequestResponse(c, "Invalid designation identifier", nil)
	}

	reason := strings.TrimSpace(c.FormValue("deactive_reason"))
	if reason == "" {
		reason = strings.TrimSpace(c.FormValue("reason"))
	}
	if reason == "" {
		return webResponse.SendBadRequestResponse(c, "Deactive reason is required", nil)
	}
	if len(reason) > 250 {
		return webResponse.SendBadRequestResponse(c, "Deactive reason must not exceed 250 characters", nil)
	}

	d, err := h.Service.DeactiveDesignation(uUId, reason)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	_ = h.Service.CreateUserActivity(c, constants.ActivityDeactivate, "designations", d.Id, "Deactive Designation Information")

	return webResponse.SendSuccessResponse(c, "Deactivated Successfully", nil)
}

func (h Handler) DesignationHistory(c echo.Context) error {
	var uUId requests.UUIDRequest
	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil {
		return webResponse.SendBadRequestResponse(c, "Invalid designation identifier", nil)
	}

	histories, err := h.Service.GetDesignationHistory(uUId)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Designation history", histories)
}
