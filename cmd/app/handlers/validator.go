package handlers

import (
	"admin-dashboard/cmd/common/helpers"
	"admin-dashboard/cmd/common/response"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

func (h Handler) PayloadValidate(c echo.Context, payload interface{}) []*response.ValidationError {

	if err := h.BindBodyRequest(c, payload); err != nil {

		return []*response.ValidationError{
			{
				Error:     err.Error(),
				Key:       "binding",
				Condition: "binding",
			},
		}
	}

	validationErrors := h.ValidateBodyRequest(payload)

	if len(validationErrors) > 0 {
		return validationErrors
	}

	return nil
}

func (h Handler) BindBodyRequest(c echo.Context, payload interface{}) error {

	if err := (&echo.DefaultBinder{}).BindBody(c, payload); err != nil {
		return errors.New("failed to bind body")
	}

	return nil
}

func (h Handler) ValidateBodyRequest(payload interface{}) []*response.ValidationError {

	validate := validator.New(
		validator.WithRequiredStructEnabled(),
	)

	userRole, err := h.Service.UserRole()

	if err != nil {

		return []*response.ValidationError{
			{
				Error:     "Failed to load validation rules",
				Key:       "validator",
				Condition: "internal",
			},
		}
	}

	minimum := int(userRole.PasswordMinLength)
	maximum := int(userRole.PasswordMaxLength)

	mobileRegex := regexp.MustCompile(
		`^\+8801[3-9][0-9]{8}$`,
	)

	if err := validate.RegisterValidation(
		"mobile",
		func(fl validator.FieldLevel) bool {

			mobile := fl.Field().String()

			mobile = helpers.NormalizeMobile(mobile)

			return mobileRegex.MatchString(mobile)
		},
	); err != nil {

		return validatorInitError("mobile")
	}

	if err := validate.RegisterValidation(
		"password",
		func(fl validator.FieldLevel) bool {

			password := strings.TrimSpace(fl.Field().String())

			if len(password) < minimum || len(password) > maximum {

				return false
			}

			var hasUpper bool
			var hasLower bool
			var hasNumber bool

			for _, c := range password {

				switch {

				case c >= 'A' && c <= 'Z':
					hasUpper = true

				case c >= 'a' && c <= 'z':
					hasLower = true

				case c >= '0' && c <= '9':
					hasNumber = true
				}
			}

			return hasUpper &&
				hasLower &&
				hasNumber
		},
	); err != nil {
		return validatorInitError("password")
	}

	if err := validate.RegisterValidation("file", validateFile); err != nil {

		return validatorInitError("file")
	}

	if err := validate.RegisterValidation("minsize", validateMinFileSize); err != nil {

		return validatorInitError("minsize")
	}

	if err := validate.RegisterValidation("maxsize_mb", validateMaxFileSizeMB); err != nil {

		return validatorInitError("maxsize_mb")
	}

	if err := validate.RegisterValidation("imagetype", validateImageType); err != nil {

		return validatorInitError("imagetype")
	}

	if err := validate.RegisterValidation("filetype", validateGeneralFileType); err != nil {

		return validatorInitError("filetype")
	}

	err = validate.Struct(payload)

	if err == nil {

		h.normalizeFields(payload)

		return nil
	}

	var validationErrors []*response.ValidationError

	var ve validator.ValidationErrors

	if !errors.As(err, &ve) {

		return []*response.ValidationError{
			{
				Error:     "Validation failed",
				Key:       "validator",
				Condition: "validation",
			},
		}
	}

	refVal := reflect.ValueOf(payload)

	if refVal.Kind() == reflect.Ptr {
		refVal = refVal.Elem()
	}

	for _, fieldErr := range ve {

		if refVal.Kind() != reflect.Struct {
			continue
		}

		structField, ok :=
			refVal.Type().FieldByName(
				fieldErr.StructField(),
			)

		if !ok {
			continue
		}

		jsonKey := structField.Tag.Get("json")

		if strings.Contains(jsonKey, ",") {

			jsonKey = strings.Split(jsonKey, ",")[0]
		}
		if jsonKey == "" {

			jsonKey = strings.ToLower(fieldErr.StructField())
		}

		param := fieldErr.Param()

		condition := fieldErr.Tag()

		keyTitle := strings.ReplaceAll(jsonKey, "_", " ")

		if len(keyTitle) > 0 {

			keyTitle = strings.ToUpper(keyTitle[:1]) + keyTitle[1:]
		}

		var msg string

		switch condition {

		case "required":
			msg = fmt.Sprintf("%s is required", keyTitle)
		case "email":
			msg = fmt.Sprintf("%s must be a valid email address", keyTitle)
		case "mobile":
			msg = fmt.Sprintf("%s must be a valid mobile number", keyTitle)
		case "numeric":
			msg = fmt.Sprintf("%s can only contain digits (0-9)", keyTitle)
		case "uuid4":
			msg = fmt.Sprintf("%s must be a valid UUID", keyTitle)
		case "len":
			msg = fmt.Sprintf("%s must be exactly %s characters", keyTitle, param)
		case "password":
			msg = fmt.Sprintf("%s must be %d-%d characters and contain uppercase, lowercase and number", keyTitle, minimum, maximum)
		case "min":
			msg = fmt.Sprintf("%s must be at least %s characters", keyTitle, param)
		case "gtefield":
			msg = fmt.Sprintf("%s must be greater than or equal to %s", keyTitle, strings.ToLower(strings.ReplaceAll(param, "_", " ")))
		case "max":
			msg = fmt.Sprintf("%s cannot be greater than %s characters", keyTitle, param)
		case "eqfield":
			msg = fmt.Sprintf("%s must match %s", keyTitle, strings.ToLower(param))
		case "oneof":
			paramText := strings.ReplaceAll(param, " ", " or ")
			msg = fmt.Sprintf("%s must be one of %s", keyTitle, strings.ToLower(paramText))
		case "file":
			msg = fmt.Sprintf("%s is not a valid file", keyTitle)
		case "minsize":
			minSize, _ := strconv.ParseInt(param, 10, 64)
			msg = fmt.Sprintf("%s must be at least %s", keyTitle, formatFileSize(minSize))
		case "maxsize_mb":
			msg = fmt.Sprintf("%s cannot exceed %s MB", keyTitle, param)
		case "imagetype":
			msg = fmt.Sprintf("%s must be a valid image file", keyTitle)
		case "filetype":
			msg = fmt.Sprintf("%s has an unsupported file type", keyTitle)
		default:
			msg = fmt.Sprintf("%s field is %s", keyTitle, condition)
		}

		validationErrors = append(validationErrors, &response.ValidationError{
			Error:     msg,
			Key:       jsonKey,
			Condition: condition,
		},
		)
	}

	return validationErrors
}

func getFileHeader(field reflect.Value) (*multipart.FileHeader, bool) {

	if !field.IsValid() {
		return nil, false
	}

	if field.Kind() == reflect.Ptr {

		if field.IsNil() {
			return nil, true
		}

		file, ok := field.Interface().(*multipart.FileHeader)

		if ok {
			return file, true
		}

		field = field.Elem()
	}

	if field.CanInterface() {

		if file, ok := field.Interface().(multipart.FileHeader); ok {

			return &file, true
		}
	}

	return nil, false
}

func validateFile(fl validator.FieldLevel) bool {

	file, ok := getFileHeader(fl.Field())

	if !ok {
		return false
	}

	if file == nil {
		return true
	}

	return file.Size > 0
}

func validateMinFileSize(fl validator.FieldLevel) bool {

	file, ok := getFileHeader(fl.Field())

	if !ok {
		return false
	}

	if file == nil {
		return true
	}

	minSize, err := strconv.ParseInt(fl.Param(), 10, 64)

	if err != nil {
		return false
	}

	return file.Size >= minSize
}

func validateMaxFileSizeMB(fl validator.FieldLevel) bool {

	file, ok := getFileHeader(fl.Field())

	if !ok {
		return false
	}

	if file == nil {
		return true
	}

	maxMB, err := strconv.ParseFloat(fl.Param(), 64)

	if err != nil {
		return false
	}

	maxBytes := int64(maxMB * 1024 * 1024)

	return file.Size <= maxBytes
}

var allowedImageTypes = map[string]bool{
	"image/jpeg":               true,
	"image/png":                true,
	"image/gif":                true,
	"image/webp":               true,
	"image/bmp":                true,
	"image/tiff":               true,
	"image/svg+xml":            true,
	"image/x-icon":             true,
	"image/vnd.microsoft.icon": true,
}

var allowedFileTypes = map[string]bool{

	"image/jpeg":               true,
	"image/png":                true,
	"image/gif":                true,
	"image/webp":               true,
	"image/bmp":                true,
	"image/tiff":               true,
	"image/svg+xml":            true,
	"image/x-icon":             true,
	"image/vnd.microsoft.icon": true,
	"application/pdf":          true,
	"text/plain":               true,
	"text/csv":                 true,
	"text/xml":                 true,
	"text/html":                true,
	"text/css":                 true,
	"application/csv":          true,
	"application/xml":          true,
	"application/xaml+xml":     true,
	"application/json":         true,
	"application/vnd.ms-excel": true,
	"application/vnd.openxmlformats-officedocument.spreadsheetml.sheet": true,
	"application/vnd.ms-excel.sheet.macroEnabled.12":                    true,
	"application/vnd.ms-excel.sheet.binary.macroEnabled.12":             true,
	"application/msword": true,
	"application/vnd.openxmlformats-officedocument.wordprocessingml.document":   true,
	"application/vnd.ms-word.document.macroEnabled.12":                          true,
	"application/vnd.ms-powerpoint":                                             true,
	"application/vnd.openxmlformats-officedocument.presentationml.presentation": true,
	"application/vnd.ms-powerpoint.presentation.macroEnabled.12":                true,
	"application/zip":                                 true,
	"application/x-zip-compressed":                    true,
	"application/x-7z-compressed":                     true,
	"application/vnd.rar":                             true,
	"application/x-rar-compressed":                    true,
	"application/gzip":                                true,
	"application/x-gzip":                              true,
	"application/x-tar":                               true,
	"application/rtf":                                 true,
	"application/vnd.oasis.opendocument.text":         true,
	"application/vnd.oasis.opendocument.spreadsheet":  true,
	"application/vnd.oasis.opendocument.presentation": true,
	"audio/mpeg":                                      true,
	"audio/wav":                                       true,
	"audio/ogg":                                       true,
	"audio/webm":                                      true,
	"video/mp4":                                       true,
	"video/mpeg":                                      true,
	"video/webm":                                      true,
	"video/ogg":                                       true,
	"video/x-msvideo":                                 true,
}

func validateImageType(fl validator.FieldLevel) bool {

	file, ok := getFileHeader(fl.Field())

	if !ok {
		return false
	}

	if file == nil {
		return true
	}

	contentType := strings.ToLower(file.Header.Get("Content-Type"))

	if !allowedImageTypes[contentType] {
		return false
	}

	if contentType == "image/svg+xml" {
		return true
	}

	return detectActualFileType(
		file,
	)
}

func validateGeneralFileType(fl validator.FieldLevel) bool {

	file, ok := getFileHeader(fl.Field())

	if !ok {
		return false
	}

	if file == nil {
		return true
	}

	contentType := strings.ToLower(file.Header.Get("Content-Type"))

	if !allowedFileTypes[contentType] {
		return false
	}

	if strings.HasPrefix(contentType, "text/") {

		return true
	}

	if contentType == "image/svg+xml" || contentType == "application/xml" || contentType == "application/xaml+xml" || contentType == "application/json" {

		return true
	}

	return detectActualFileType(
		file,
	)
}

func detectActualFileType(file *multipart.FileHeader) bool {

	src, err := file.Open()

	if err != nil {
		return false
	}

	defer func(src multipart.File) {
		_ = src.Close()

	}(src)

	buffer :=
		make([]byte, 512)

	n, err := src.Read(buffer)

	if err != nil && n == 0 {
		return false
	}

	detectedType := strings.ToLower(http.DetectContentType(buffer[:n]))

	if detectedType == "application/octet-stream" {
		return true
	}

	return true
}

func validatorInitError(key string) []*response.ValidationError {

	return []*response.ValidationError{
		{
			Error:     "Validator initialization failed",
			Key:       key,
			Condition: "internal",
		},
	}
}

func (h Handler) normalizeFields(payload interface{}) {

	refVal := reflect.ValueOf(payload)

	if refVal.Kind() != reflect.Ptr {
		return
	}

	refVal = refVal.Elem()

	if refVal.Kind() != reflect.Struct {
		return
	}

	for i := 0; i < refVal.NumField(); i++ {

		field := refVal.Field(i)

		typeField := refVal.Type().Field(i)

		if field.Kind() != reflect.String {
			continue
		}

		tag := typeField.Tag.Get("validate")

		if strings.Contains(tag, "mobile") {

			normalized := helpers.NormalizeMobile(field.String())

			if field.CanSet() {
				field.SetString(normalized)
			}
		}
	}
}

func formatFileSize(size int64) string {

	if size < 1024 {
		return fmt.Sprintf("%d bytes", size)
	}

	if size < 1024*1024 {

		return fmt.Sprintf("%.2f KB", float64(size)/1024)
	}

	return fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
}
