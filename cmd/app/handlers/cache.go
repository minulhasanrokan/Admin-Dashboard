package handlers

import (
	"admin-dashboard/cmd/common/response/webResponse"
	"admin-dashboard/cmd/internal/constants"
	"admin-dashboard/cmd/internal/dtos"
	"admin-dashboard/cmd/internal/session"

	"github.com/labstack/echo/v4"
)

func (h Handler) ClearCachePage(c echo.Context) error {

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

	html, err := webResponse.RenderHTML(c, "settings/clear_cache", map[string]any{"RightDetails": rightDetails})

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Clear cache", html)
}

func (h Handler) ClearCache(c echo.Context) error {

	err := h.Service.ClearCache()

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	err = h.Service.CreateUserActivity(c, constants.ActivityDelete, "caches", 0, "Clear all application cache")

	if err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "All cache cleared successfully.", nil)
}
