package handlers

import (
	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/app/requests"
	"admin-dashboard/cmd/common/response/webResponse"
	"admin-dashboard/cmd/internal/constants"
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

func (h Handler) DesignationPage(c echo.Context) error {

	gridRightJSON, err := h.UserGridRight(c)

	if err != nil {

		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	html, err := webResponse.RenderHTML(c, "reference/designation/view", gridRightJSON)

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Designations", html)
}

func (h Handler) DesignationAddPage(c echo.Context) error {

	gridRightJSON, err := h.UserGridRight(c)

	if err != nil {

		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	html, err := webResponse.RenderHTML(c, "reference/designation/add", gridRightJSON)

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Designations", html)
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

	payload := new(requests.CreateDesignationRequest)

	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	name := strings.TrimSpace(payload.GetName())
	code := strings.TrimSpace(payload.GetCode())
	title := strings.TrimSpace(payload.GetTitle())
	details := strings.TrimSpace(payload.GetDetails())

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
	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil || uUId.UUID == "" {
		uUId.UUID = c.Param("uuid")
	}

	d, err := h.Service.GetDesignationByIdOrUUID(uUId)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}
	if d == nil {
		return webResponse.SendNotFoundResponse(c, "Designation not found", nil)
	}

	historyData, _ := h.Service.GetLastApprovedDesignationHistory(d.Id)

	nameVal := ""
	if d.Name != nil {
		nameVal = *d.Name
	}
	codeVal := ""
	if d.Code != nil {
		codeVal = *d.Code
	}
	titleVal := ""
	if d.Title != nil {
		titleVal = *d.Title
	}
	detailsVal := ""
	if d.Details != nil {
		detailsVal = *d.Details
	}

	opStatusText := "New Add"
	switch d.ActualStatus {
	case 1:
		opStatusText = "New Add"
	case 2:
		opStatusText = "Add-Edit"
	case 3:
		opStatusText = "Edit"
	case 4:
		opStatusText = "Delete"
	case 5:
		opStatusText = "Deactivate"
	case 6:
		opStatusText = "Active"
	}

	processStatusText := "Draft"
	switch d.ProcessStatus {
	case 0:
		processStatusText = "Draft"
	case 1:
		processStatusText = "Sent to Checker"
	case 2:
		processStatusText = "Returned"
	case 3:
		processStatusText = "Approved"
	}

	checkerText := "None"
	if d.Checker > 0 {
		checkerUser, _ := h.Service.GetUserByID(d.Checker)
		if checkerUser != nil && checkerUser.Name != nil {
			checkerText = *checkerUser.Name
		} else {
			checkerText = fmt.Sprintf("User ID %d", d.Checker)
		}
	}

	res := map[string]interface{}{
		"SingleData":          d,
		"HistoryData":         historyData,
		"id":                  d.Id,
		"uuid":                d.UUID,
		"name":                nameVal,
		"code":                codeVal,
		"title":               titleVal,
		"details":             detailsVal,
		"status":              d.Status,
		"actual_status":       d.ActualStatus,
		"op_status_text":      opStatusText,
		"process_status":      d.ProcessStatus,
		"process_status_text": processStatusText,
		"checker":             d.Checker,
		"checker_name":        checkerText,
		"notify_type":         d.NotifyType,
		"message":             d.Message,
		"remarks":             d.Remarks,
		"created_at":          d.CreatedAt.Format("02 Jan, 2006 15:04:05"),
		"updated_at":          d.UpdatedAt.Format("02 Jan, 2006 15:04:05"),
	}

	html, err := webResponse.RenderHTML(c, "reference/designation/single_view", res)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Designation details", html)
}

func (h Handler) DesignationEditPage(c echo.Context) error {
	var uUId requests.UUIDRequest
	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil || uUId.UUID == "" {
		uUId.UUID = c.Param("uuid")
	}

	d, err := h.Service.GetDesignationByIdOrUUID(uUId)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}
	if d == nil {
		return webResponse.SendNotFoundResponse(c, "Designation not found", nil)
	}

	html, err := webResponse.RenderHTML(c, "reference/designation/edit", map[string]interface{}{
		"SingleData":  d,
		"Designation": d,
	})
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Designation data", html)
}

func (h Handler) DesignationEdit(c echo.Context) error {

	var uUId requests.UUIDRequest

	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil || uUId.UUID == "" {
		uUId.UUID = c.Param("uuid")
	}
	if uUId.UUID == "" {
		return webResponse.SendBadRequestResponse(c, "Invalid designation identifier", nil)
	}

	existing, err := h.Service.GetDesignationByIdOrUUID(uUId)

	if err != nil || existing == nil {
		return webResponse.SendNotFoundResponse(c, "Designation not found or already deleted", nil)
	}

	payload := new(requests.UpdateDesignationRequest)

	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	name := strings.TrimSpace(payload.GetName())
	code := strings.TrimSpace(payload.GetCode())
	title := strings.TrimSpace(payload.GetTitle())
	details := strings.TrimSpace(payload.GetDetails())

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
	var uUId requests.UUIDRequest
	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil || uUId.UUID == "" {
		uUId.UUID = c.Param("uuid")
	}

	d, err := h.Service.GetDesignationByIdOrUUID(uUId)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}
	if d == nil {
		return webResponse.SendNotFoundResponse(c, "Designation not found", nil)
	}

	html, err := webResponse.RenderHTML(c, "reference/designation/delete", map[string]interface{}{
		"SingleData": d,
		"uuid":       d.UUID,
	})
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Designation delete form", html)
}

func (h Handler) DesignationDelete(c echo.Context) error {

	var uUId requests.UUIDRequest

	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil || uUId.UUID == "" {
		uUId.UUID = c.Param("uuid")
	}
	if uUId.UUID == "" {
		return webResponse.SendBadRequestResponse(c, "Invalid designation identifier", nil)
	}

	payload := new(requests.DeleteDesignationRequest)

	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	reason := strings.TrimSpace(payload.DeleteReason)

	d, err := h.Service.DeleteDesignationWithReason(uUId, reason)

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return webResponse.SendNotFoundResponse(c, "Designation not found or already deleted", nil)
		}
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	_ = h.Service.CreateUserActivity(c, constants.ActivityDelete, "designations", d.Id, "Delete Designation Information")

	return webResponse.SendSuccessResponse(c, "Deleted Successfully", nil)
}

func (h Handler) DesignationStcPage(c echo.Context) error {
	var uUId requests.UUIDRequest
	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil || uUId.UUID == "" {
		uUId.UUID = c.Param("uuid")
	}

	d, err := h.Service.GetDesignationByIdOrUUID(uUId)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}
	if d == nil {
		return webResponse.SendNotFoundResponse(c, "Designation not found", nil)
	}

	historyData, _ := h.Service.GetLastApprovedDesignationHistory(d.Id)

	checkers, err := h.Service.GetCheckerUserList()
	if err != nil || len(checkers) == 0 {
		checkers, _ = h.Service.GetDynamicDropdown(requests.GetDropdownRequest{
			TableName: "users",
			FieldName: "id,name",
		})
	}

	html, err := webResponse.RenderHTML(c, "reference/designation/sent_to_checker", map[string]interface{}{
		"SingleData":   d,
		"HistoryData":  historyData,
		"CheckerUsers": checkers,
		"uuid":         d.UUID,
	})
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Sent to checker form", html)
}

func (h Handler) DesignationStc(c echo.Context) error {

	var uUId requests.UUIDRequest

	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil || uUId.UUID == "" {
		uUId.UUID = c.Param("uuid")
	}
	if uUId.UUID == "" {
		return webResponse.SendBadRequestResponse(c, "Invalid designation identifier", nil)
	}

	payload := new(requests.SentToCheckerDesignationRequest)

	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	notifyType := payload.NotifyType

	if notifyType == 0 {

		authSetting, _ := h.Service.AuthenticationSetting()

		if authSetting != nil && strings.ToLower(authSetting.OtpType) == "sms" {
			notifyType = 2
		} else {
			notifyType = 1
		}
	}

	remarks := strings.TrimSpace(payload.Remarks)

	d, err := h.Service.SentToCheckerDesignation(uUId, payload.Checker, notifyType, remarks)

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return webResponse.SendNotFoundResponse(c, "Designation not found or already deleted", nil)
		}
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	_ = h.Service.CreateUserActivity(c, constants.ActivitySentToChecker, "designations", d.Id, "Sent To Checker Designation Information")

	return webResponse.SendSuccessResponse(c, "Sent To Checker Successfully", nil)
}

func (h Handler) DesignationVerifyPage(c echo.Context) error {
	var uUId requests.UUIDRequest
	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil || uUId.UUID == "" {
		uUId.UUID = c.Param("uuid")
	}

	d, err := h.Service.GetDesignationByIdOrUUID(uUId)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}
	if d == nil {
		return webResponse.SendNotFoundResponse(c, "Designation not found", nil)
	}

	historyData, _ := h.Service.GetLastApprovedDesignationHistory(d.Id)

	opStatusText := "New Add"
	switch d.ActualStatus {
	case 1:
		opStatusText = "New Add"
	case 2:
		opStatusText = "Add-Edit"
	case 3:
		opStatusText = "Edit"
	case 4:
		opStatusText = "Delete"
	case 5:
		opStatusText = "Deactivate"
	case 6:
		opStatusText = "Active"
	}

	checkerText := "None"
	if d.Checker > 0 {
		checkerUser, _ := h.Service.GetUserByID(d.Checker)
		if checkerUser != nil && checkerUser.Name != nil {
			checkerText = *checkerUser.Name
		} else {
			checkerText = fmt.Sprintf("User ID %d", d.Checker)
		}
	}

	html, err := webResponse.RenderHTML(c, "reference/designation/verify", map[string]interface{}{
		"SingleData":   d,
		"HistoryData":  historyData,
		"OpStatusText": opStatusText,
		"CheckerName":  checkerText,
		"uuid":         d.UUID,
	})
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Verify designation form", html)
}

func (h Handler) DesignationVerify(c echo.Context) error {

	var uUId requests.UUIDRequest

	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil || uUId.UUID == "" {
		uUId.UUID = c.Param("uuid")
	}
	if uUId.UUID == "" {
		return webResponse.SendBadRequestResponse(c, "Invalid designation identifier", nil)
	}

	payload := new(requests.VerifyDesignationRequest)

	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	verifyType := payload.VerifyType
	returnReason := strings.TrimSpace(payload.ReturnReason)

	if verifyType == 0 && returnReason == "" {
		return webResponse.SendBadRequestResponse(c, "Return reason is required", nil)
	}

	d, des, err := h.Service.VerifyDesignation(uUId, verifyType, returnReason)

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return webResponse.SendNotFoundResponse(c, "Designation not found or already deleted", nil)
		}
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	if verifyType == 1 {

		_ = h.Service.CreateUserActivity(c, constants.ActivityApprove, "designations", d.Id, des)
	} else {
		_ = h.Service.CreateUserActivity(c, constants.ActivityReturn, "designations", d.Id, des)
	}

	if verifyType == 1 {
		return webResponse.SendSuccessResponse(c, "Approved Successfully", nil)
	}

	return webResponse.SendSuccessResponse(c, "Returned Successfully", nil)
}

func (h Handler) DesignationActivePage(c echo.Context) error {
	var uUId requests.UUIDRequest
	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil || uUId.UUID == "" {
		uUId.UUID = c.Param("uuid")
	}

	d, err := h.Service.GetDesignationByIdOrUUID(uUId)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}
	if d == nil {
		return webResponse.SendNotFoundResponse(c, "Designation not found", nil)
	}

	html, err := webResponse.RenderHTML(c, "reference/designation/active", map[string]interface{}{
		"SingleData": d,
		"uuid":       d.UUID,
	})
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Active designation form", html)
}

func (h Handler) DesignationActive(c echo.Context) error {

	var uUId requests.UUIDRequest

	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil || uUId.UUID == "" {
		uUId.UUID = c.Param("uuid")
	}
	if uUId.UUID == "" {
		return webResponse.SendBadRequestResponse(c, "Invalid designation identifier", nil)
	}

	payload := new(requests.StatusReasonDesignationRequest)

	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	reason := strings.TrimSpace(payload.GetReason())

	if reason == "" {
		return webResponse.SendBadRequestResponse(c, "Active reason is required", nil)
	}

	d, err := h.Service.ActiveDesignation(uUId, reason)

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return webResponse.SendNotFoundResponse(c, "Designation not found or already deleted", nil)
		}
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	_ = h.Service.CreateUserActivity(c, constants.ActivityActivate, "designations", d.Id, "Active Designation Information")

	return webResponse.SendSuccessResponse(c, "Activated Successfully", nil)
}

func (h Handler) DesignationDeactivatePage(c echo.Context) error {
	var uUId requests.UUIDRequest
	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil || uUId.UUID == "" {
		uUId.UUID = c.Param("uuid")
	}

	d, err := h.Service.GetDesignationByIdOrUUID(uUId)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}
	if d == nil {
		return webResponse.SendNotFoundResponse(c, "Designation not found", nil)
	}

	html, err := webResponse.RenderHTML(c, "reference/designation/deactive", map[string]interface{}{
		"SingleData": d,
		"uuid":       d.UUID,
	})
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Deactivate designation form", html)
}

func (h Handler) DesignationDeactivate(c echo.Context) error {

	var uUId requests.UUIDRequest

	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil || uUId.UUID == "" {
		uUId.UUID = c.Param("uuid")
	}
	if uUId.UUID == "" {
		return webResponse.SendBadRequestResponse(c, "Invalid designation identifier", nil)
	}

	payload := new(requests.StatusReasonDesignationRequest)

	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	reason := strings.TrimSpace(payload.GetReason())

	if reason == "" {
		return webResponse.SendBadRequestResponse(c, "Deactivate reason is required", nil)
	}

	d, err := h.Service.DeactiveDesignation(uUId, reason)

	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "not found") {
			return webResponse.SendNotFoundResponse(c, "Designation not found or already deleted", nil)
		}
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	_ = h.Service.CreateUserActivity(c, constants.ActivityDeactivate, "designations", d.Id, "Deactivate Designation Information")

	return webResponse.SendSuccessResponse(c, "Deactivated Successfully", nil)
}

func (h Handler) DesignationHistory(c echo.Context) error {
	var uUId requests.UUIDRequest
	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil || uUId.UUID == "" {
		uUId.UUID = c.Param("uuid")
	}
	if uUId.UUID == "" {
		return webResponse.SendBadRequestResponse(c, "Invalid designation identifier", nil)
	}

	histories, err := h.Service.GetDesignationHistory(uUId)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	activities, _ := h.Service.GetDesignationActivities(uUId)

	html, err := webResponse.RenderHTML(c, "reference/designation/history", map[string]interface{}{
		"Histories":  histories,
		"Activities": activities,
	})
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Designation history", html)
}

func (h Handler) DesignationActivity(c echo.Context) error {
	var uUId requests.UUIDRequest
	if err := (&echo.DefaultBinder{}).BindPathParams(c, &uUId); err != nil || uUId.UUID == "" {
		uUId.UUID = c.Param("uuid")
	}
	if uUId.UUID == "" {
		return webResponse.SendBadRequestResponse(c, "Invalid designation identifier", nil)
	}

	activities, err := h.Service.GetDesignationActivities(uUId)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	histories, _ := h.Service.GetDesignationHistory(uUId)

	html, err := webResponse.RenderHTML(c, "reference/designation/history", map[string]interface{}{
		"Histories":  histories,
		"Activities": activities,
	})
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Designation activities", html)
}
