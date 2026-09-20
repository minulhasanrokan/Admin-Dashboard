package handlers

import (
	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/app/services"
	"admin-dashboard/cmd/common/helpers"
	"admin-dashboard/cmd/common/mailService"
	"admin-dashboard/cmd/common/response/webResponse"
	"admin-dashboard/cmd/internal/cacheSystem"
	"admin-dashboard/cmd/internal/config"
	"admin-dashboard/cmd/internal/constants"
	"admin-dashboard/cmd/internal/databases"
	"admin-dashboard/cmd/internal/storage"
	"errors"
	"html/template"
	"io"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
	"gorm.io/gorm"
)

type Handler struct {
	Db          databases.DbConnection
	Logger      echo.Logger
	Repository  repositories.Repository
	Service     services.Service
	SendMail    mailService.SendMailService
	Config      config.Config
	Storage     storage.Storage
	Templates   *template.Template
	CacheSystem cacheSystem.CacheSystem
}
type HealthResponse struct {
	Health bool `json:"health"`
}

func (h Handler) HealthCheck(c echo.Context) error {

	data := HealthResponse{
		Health: true,
	}

	return c.JSON(http.StatusOK, data)
}

func (h Handler) RenderWithErrorCheck(c echo.Context, data map[string]interface{}, recordErr error, okTemplate string) error {

	if data == nil {
		data = map[string]interface{}{}
	}

	data["SystemError"] = recordErr

	setting, err := h.Service.SystemSetting()

	if err == nil {
		data["SystemSetting"] = setting
	}

	authSetting, err := h.Service.AuthenticationSetting()
	if err == nil {
		data["AuthSetting"] = authSetting
	}

	data["ConfigData"] = h.Config
	data["DateTimeConstants"] = helpers.GetDateTimeConstants()
	data["PaginationConstants"] = map[string]any{
		"DefaultLimit": constants.DefaultPageLimit,
		"MaxLimit":     constants.MaxPageLimit,
		"PageLimits":   constants.PageLimits,
	}

	if recordErr != nil {

		if errors.Is(recordErr, gorm.ErrRecordNotFound) {

			return c.Render(http.StatusNotFound, "404-error.gohtml", data)
		}

		return c.Render(http.StatusInternalServerError, "500-error.gohtml", data)
	}

	return c.Render(http.StatusOK, okTemplate+".gohtml", data)
}

func (h Handler) GetFile(c echo.Context) error {

	fileName := c.Param("*")

	file, err := h.Storage.GetFile(fileName)

	if err != nil {

		return webResponse.SendNotFoundResponse(c, "File not found", err)
	}

	defer func(file io.ReadCloser) {
		err := file.Close()
		if err != nil {

			_ = file.Close()
			return
		}
	}(file)

	contentType := "application/octet-stream"

	ext := strings.ToLower(filepath.Ext(fileName))

	switch ext {

	case ".png":
		contentType = "image/png"
	case ".jpg", ".jpeg":
		contentType = "image/jpeg"
	case ".gif":
		contentType = "image/gif"
	case ".webp":
		contentType = "image/webp"
	case ".svg":
		contentType = "image/svg+xml"
	case ".ico":
		contentType = "image/x-icon"
	case ".mp4":
		contentType = "video/mp4"
	case ".webm":
		contentType = "video/webm"
	case ".ogg":
		contentType = "video/ogg"
	case ".mp3":
		contentType = "audio/mpeg"
	case ".wav":
		contentType = "audio/wav"
	case ".pdf":
		contentType = "application/pdf"
	case ".txt":
		contentType = "text/plain"
	case ".csv":
		contentType = "text/csv"
	case ".doc":
		contentType = "application/msword"
	case ".docx":
		contentType = "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".xls":
		contentType = "application/vnd.ms-excel"
	case ".xlsx":
		contentType = "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".ppt":
		contentType = "application/vnd.ms-powerpoint"
	case ".pptx":
		contentType = "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case ".zip":
		contentType = "application/zip"
	case ".rar":
		contentType = "application/vnd.rar"
	case ".7z":
		contentType = "application/x-7z-compressed"
	default:
		contentType = "application/octet-stream"
	}

	return c.Stream(http.StatusOK, contentType, file)
}

func (h Handler) NotFoundPage(c echo.Context) error {

	setting, err := h.Service.SystemSetting()

	if err != nil {

		return webResponse.SendNotFoundResponse(c, "data not found", err)
	}

	data := map[string]interface{}{
		"SystemSetting": setting,
		"ConfigData":    h.Config,
	}

	return c.Render(http.StatusNotFound, "404-error.gohtml", data)
}

func (h Handler) ServerErrorPage(c echo.Context) error {

	setting, err := h.Service.SystemSetting()

	if err != nil {

		return webResponse.SendNotFoundResponse(c, "data not found", err)
	}

	data := map[string]interface{}{
		"SystemSetting": setting,
		"ConfigData":    h.Config,
	}

	return c.Render(http.StatusInternalServerError, "500-error.gohtml", data)
}

func (h Handler) CSRFErrorPage(c echo.Context) error {

	setting, err := h.Service.SystemSetting()

	if err != nil {

		return webResponse.SendNotFoundResponse(c, "data not found", err)
	}

	data := map[string]interface{}{
		"SystemSetting": setting,
		"ConfigData":    h.Config,
	}

	return c.Render(http.StatusForbidden, "csrf-error.gohtml", data)
}

func (h Handler) NotFoundError(c echo.Context) error {

	return webResponse.SendNotFoundResponse(c, "data not found", nil)
}

func (h Handler) ServerError(c echo.Context) error {

	return webResponse.SendInternalServerErrorResponse(c, nil)
}

func (h Handler) CSRFError(c echo.Context) error {

	return webResponse.SendCsrfErrorResponse(c, "csrf error", nil)
}
