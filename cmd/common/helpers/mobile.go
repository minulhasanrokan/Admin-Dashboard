package helpers

import "strings"

func NormalizeMobile(mobile string) string {
	mobile = strings.TrimSpace(mobile)
	mobile = strings.ReplaceAll(mobile, " ", "")
	mobile = strings.ReplaceAll(mobile, "-", "")
	mobile = strings.ReplaceAll(mobile, "(", "")
	mobile = strings.ReplaceAll(mobile, ")", "")

	if !strings.HasPrefix(mobile, "+88") {
		mobile = "+88" + mobile
	}
	return mobile
}
