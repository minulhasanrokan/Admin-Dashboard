package helpers

import (
	"admin-dashboard/cmd/internal/constants"
	"strings"
	"time"
)

func FormatViewDate(t time.Time) string {
	return t.Format(constants.ViewDateFormat)
}

func FormatViewDateTime(t time.Time) string {
	return t.Format(constants.ViewDateTimeFormat)
}

func FormatDBDate(t time.Time) string {
	return t.Format(constants.DbDateFormat)
}

func FormatDBDateTime(t time.Time) string {
	return t.Format(constants.DbDateTimeFormat)
}

func NormalizeToGoLayout(layout string) string {
	lower := strings.ToLower(layout)
	if strings.Contains(lower, "dd") || strings.Contains(lower, "yyyy") {
		r := strings.NewReplacer(
			"yyyy", "2006",
			"YYYY", "2006",
			"yy", "06",
			"YY", "06",
			"mmmm", "January",
			"MMMM", "January",
			"mmm", "Jan",
			"MMM", "Jan",
			"mm", "01",
			"MM", "01",
			"dd", "02",
			"DD", "02",
			"hh", "15",
			"HH", "15",
			"15:04", "15:04",
			"ss", "05",
			"SS", "05",
		)
		return r.Replace(layout)
	}

	res := layout
	res = strings.ReplaceAll(res, "20/", "02/")
	res = strings.ReplaceAll(res, "20-", "02-")
	res = strings.ReplaceAll(res, "20.", "02.")
	res = strings.ReplaceAll(res, "/12/", "/01/")
	res = strings.ReplaceAll(res, "-12-", "-01-")
	return res
}

func ParseInputDate(input string) (time.Time, error) {
	layout := NormalizeToGoLayout(constants.InputDateFormat)
	return time.Parse(layout, input)
}

func ParseInputDateTime(input string) (time.Time, error) {
	layout := NormalizeToGoLayout(constants.InputDateTimeFormat)
	return time.Parse(layout, input)
}

func NowDB() string {
	return time.Now().UTC().Format(constants.DbDateTimeFormat)
}

func NowView() string {
	return time.Now().UTC().Format(constants.ViewDateTimeFormat)
}

func GoLayoutToJsFormat(goLayout string) string {
	lower := strings.ToLower(goLayout)
	if strings.Contains(lower, "dd") || strings.Contains(lower, "yyyy") {
		return strings.ToUpper(goLayout)
	}

	r := strings.NewReplacer(
		"20/01/2006", "DD/MM/YYYY",
		"20/12/2026", "DD/MM/YYYY",
		"02/01/2006", "DD/MM/YYYY",
		"02-01-2006", "DD-MM-YYYY",
		"20-01-2006", "DD-MM-YYYY",
		"2006-01-02", "YYYY-MM-DD",
		"02 Jan 2006", "DD MMM YYYY",
		"2006", "YYYY",
		"06", "YY",
		"January", "MMMM",
		"Jan", "MMM",
		"02", "DD",
		"20", "DD",
		"01", "MM",
		"12", "MM",
		"15:04", "HH:mm",
		"15", "HH",
		"03", "hh",
		"04", "mm",
		"05", "ss",
		"PM", "A",
		"pm", "a",
	)
	res := r.Replace(goLayout)
	res = strings.ReplaceAll(res, "D0", "DD")
	return res
}

func GetDateTimeConstants() map[string]string {
	return map[string]string{
		"ViewDateFormat":        constants.ViewDateFormat,
		"ViewDateTimeFormat":    constants.ViewDateTimeFormat,
		"DbDateFormat":          constants.DbDateFormat,
		"DbDateTimeFormat":      constants.DbDateTimeFormat,
		"InputDateFormat":       constants.InputDateFormat,
		"InputDateTimeFormat":   constants.InputDateTimeFormat,
		"JsInputDateFormat":     GoLayoutToJsFormat(constants.InputDateFormat),
		"JsInputDateTimeFormat": GoLayoutToJsFormat(constants.InputDateTimeFormat),
		"JsViewDateFormat":      GoLayoutToJsFormat(constants.ViewDateFormat),
		"JsViewDateTimeFormat":  GoLayoutToJsFormat(constants.ViewDateTimeFormat),
	}
}
