package helpers

import (
	"net/http"
	"strconv"
	"strings"
)

func StrToInt(value string) (int, error) {

	v, err := strconv.Atoi(value)

	if err != nil {
		return 0, err
	}

	return v, nil
}

func StrToBool(value string) (bool, error) {

	v, err := strconv.ParseBool(value)

	if err != nil {
		return false, err
	}

	return v, nil
}

func StrToSameSite(value string) http.SameSite {

	switch strings.ToLower(value) {
	case "strict":
		return http.SameSiteStrictMode

	case "lax":
		return http.SameSiteLaxMode

	case "none":
		return http.SameSiteNoneMode

	default:
		return http.SameSiteDefaultMode
	}
}

func Str(s string) *string {
	return &s
}

func UInt(v uint64) *uint64 {
	return &v
}

func Int(v uint64) *int {
	val := int(v)
	return &val
}
