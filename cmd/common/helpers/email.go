package helpers

import "strings"

func MaskEmail(email string) string {
	parts := strings.Split(email, "@")
	if len(parts) != 2 {
		return email
	}

	name := parts[0]
	domain := parts[1]

	if len(name) <= 2 {
		return name[:1] + "***@" + domain
	}

	return name[:1] + strings.Repeat("*", len(name)-2) + name[len(name)-1:] + "@" + domain
}
