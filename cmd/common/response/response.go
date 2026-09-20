package response

type Response map[string]interface{}

type ApiJsonSuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type WebJsonSuccessResponse struct {
	Success bool        `json:"success"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
	CSRF    string      `json:"_csrf"`
}

type ApiJsonFailedValidationResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Errors  []*ValidationError `json:"errors"`
}

type WebJsonFailedValidationResponse struct {
	Success bool               `json:"success"`
	Message string             `json:"message"`
	Errors  []*ValidationError `json:"errors"`
	CSRF    string             `json:"_csrf"`
}

type ApiJsonErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Error   string `json:"error"`
}

type WebJsonErrorResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	CSRF    string `json:"_csrf"`
	Error   string `json:"error"`
}

type ValidationError struct {
	Error     string `json:"error"`
	Key       string `json:"key"`
	Condition string `json:"condition"`
}
