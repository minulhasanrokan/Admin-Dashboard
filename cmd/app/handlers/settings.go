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
	"time"

	"github.com/labstack/echo/v4"
)

func (h Handler) SystemSetting(c echo.Context) error {

	err := h.CacheSystem.Flush()
	if err != nil {
		return err
	}

	systemSetting, err := h.Service.SystemSetting()

	if err != nil {
		return webResponse.SendNotFoundResponse(c, "System settings not found", err)
	}

	userSess, sess, err := session.GetUserSession(c)

	if err != nil || userSess.IsEmpty() || !userSess.IsActive || userSess.CategoryRouteAction == "" {
		_ = session.ClearUserSession(c)
		return webResponse.SendUnauthorizedResponse(c, "Session expired", nil)
	}

	var rightDetails []dtos.RightInfo

	if userSess.SuperAdmin == 1 {
		rightDetails, err = h.Service.GetRightDetailsByCategoryKeyKey(userSess.CategoryRouteAction)
	} else {
		rightDetails, err = h.Service.GetRightDetailsByCategoryKeyUserId(userSess.UserId, userSess.CategoryRouteAction)
	}

	if err != nil {
		_ = session.ClearUserSession(c)
		return webResponse.SendUnauthorizedResponse(c, "Session expired", nil)
	}

	templates := map[string]string{
		"settings.system.basic.settings":  "settings/system_view",
		"settings.system.system.settings": "settings/email_setting",
		"settings.system.other.settings":  "settings/clear_cache",
	}

	for _, right := range rightDetails {

		template, ok := templates[right.RightAction]

		if !ok {
			continue
		}

		rightInfo, err := h.Service.GetRightDetailsByKey(right.RightAction)

		if err != nil {
			return webResponse.SendInternalServerErrorResponse(c, err)
		}

		now := time.Now()

		userSess.LastActiveTime = &now
		userSess.RouteAction = right.RightAction
		userSess.RouteActionDetails = rightInfo

		sess.Values[constants.SessionKeyRouteAction] = right.RightAction
		sess.Values[constants.SessionKeyRouteActionDetails] = rightInfo
		sess.Values[constants.SessionKeyLastActiveTime] = now.Format(time.RFC3339)

		if err := sess.Save(c.Request(), c.Response()); err != nil {
			return webResponse.SendInternalServerErrorResponse(c, err)
		}

		c.Set(constants.SessionName, userSess)

		html, err := webResponse.RenderHTML(c, template, map[string]any{"SystemSetting": systemSetting, "RightDetails": rightDetails})

		if err != nil {
			return webResponse.SendInternalServerErrorResponse(c, err)
		}

		return webResponse.SendSuccessResponse(c, "System settings", html)
	}

	return webResponse.SendNotFoundResponse(c, "No page found for this permission", nil)
}

func (h Handler) BusinessSettingPage(c echo.Context) error {

	systemSetting, err := h.Service.SystemSetting()

	if err != nil {
		return webResponse.SendNotFoundResponse(c, "System settings not found", err)
	}

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

	html, err := webResponse.RenderHTML(c, "settings/system_view", map[string]any{"SystemSetting": systemSetting, "RightDetails": rightDetails})

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "System settings", html)
}

func (h Handler) BusinessSystemSettingUpdate(c echo.Context) error {

	payload := new(requests.BusinessSettingRequest)

	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	systemSetting, err := h.Service.BusinessSystemSettingUpdate(payload)

	if err != nil {

		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	err = h.Service.CreateUserActivity(c, constants.ActivityEdit, "system_settings", systemSetting.Id, "Update Business Setting")

	if err != nil {

		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Business settings updated", nil)
}

func (h Handler) ImagesSettingPage(c echo.Context) error {

	systemSetting, err := h.Service.SystemSetting()

	if err != nil {
		return webResponse.SendNotFoundResponse(c, "Image settings not found", err)
	}

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

	html, err := webResponse.RenderHTML(c, "settings/images-settings", map[string]any{"SystemSetting": systemSetting, "RightDetails": rightDetails})

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "System settings", html)
}

func (h Handler) ImagesSystemSettingUpdate(c echo.Context) error {

	payload := new(requests.ImagesSettingRequest)

	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	systemSetting, err := h.Service.ImagesSystemSettingUpdate(payload)

	if err != nil {

		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	err = h.Service.CreateUserActivity(c, constants.ActivityEdit, "system_settings", systemSetting.Id, "Update Business Setting")

	if err != nil {

		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Business settings updated", nil)
}

func (h Handler) BackUpPage(c echo.Context) error {

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

	html, err := webResponse.RenderHTML(c, "settings/backup", map[string]any{"RightDetails": rightDetails})

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "backup", html)
}

func (h Handler) Backup(c echo.Context) error {

	payload := new(requests.BackupSettingRequest)

	if validationErrors := h.PayloadValidate(c, payload); validationErrors != nil {
		return webResponse.SendFailedValidationResponse(c, validationErrors)
	}

	var createdBy *uint64
	if userSess, _, err := session.GetUserSession(c); err == nil && userSess.UserId > 0 {
		createdBy = &userSess.UserId
	}

	backupHistory, err := h.Service.BackupCreate(payload, createdBy)

	if err != nil {

		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	err = h.Service.CreateUserActivity(c, constants.ActivityEdit, "backup_histories", backupHistory.Id, "Create Backup Request")

	if err != nil {

		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Backup request generated successfully", nil)
}

func (h Handler) BackupList(c echo.Context) error {

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

	result, err := h.Service.GetBackupList(params)
	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Backup list retrieved successfully", result)
}

func (h Handler) BackupDownload(c echo.Context) error {

	param := c.Param("uuid")
	if param == "" {
		param = c.Param("id")
	}
	if param == "" {
		return webResponse.SendBadRequestResponse(c, "Invalid backup identifier", nil)
	}

	backupHistory, filePath, err := h.Service.GetBackupFileByUUID(param)
	if err != nil {
		return webResponse.SendNotFoundResponse(c, "Backup file not found", err)
	}

	var downloadedBy uint64
	if userSess, _, err := session.GetUserSession(c); err == nil && userSess.UserId > 0 {
		downloadedBy = userSess.UserId
	}
	_ = h.Service.MarkBackupDownloaded(backupHistory.Id, backupHistory.UUID.String(), downloadedBy, time.Now())

	return c.Attachment(filePath, backupHistory.FileName)
}
