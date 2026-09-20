package handlers

import (
	"admin-dashboard/cmd/app/requests"
	"admin-dashboard/cmd/common/response/webResponse"

	"github.com/labstack/echo/v4"
)

func (h Handler) GetDropdown(c echo.Context) error {
	var requestData requests.GetDropdownRequest

	if err := c.Bind(&requestData); err != nil {
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	data, err := h.Service.GetDynamicDropdown(requestData)

	if err != nil {
		if err.Error() == h.Db.NotFoundMessage() {
			return webResponse.SendSuccessResponse(c, "No Data Found", nil)
		}
		return webResponse.SendInternalServerErrorResponse(c, err)
	}

	return webResponse.SendSuccessResponse(c, "Dropdown data", data)
}
