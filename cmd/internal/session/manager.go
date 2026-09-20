package session

import (
	"admin-dashboard/cmd/internal/config"
	"fmt"
	"strings"

	"github.com/gorilla/sessions"
	"gorm.io/gorm"
)

func NewSessionStore(cfg *config.Config, db *gorm.DB) (sessions.Store, error) {
	driver := strings.ToLower(strings.TrimSpace(cfg.SessionDriver))
	secret := []byte(cfg.AppKey)
	if len(secret) == 0 {
		return nil, fmt.Errorf("app key is required for session security")
	}

	maxAge := cfg.WebCookieMaxAge
	if maxAge <= 0 {
		if cfg.SessionLifetime > 0 {
			maxAge = cfg.SessionLifetime * 60
		} else {
			maxAge = 86400 * 30
		}
	}

	switch driver {
	case "file":
		return NewFileStore(cfg.SessionPath, maxAge, secret)

	case "database", "db":
		return NewDatabaseStore(db, cfg.SessionTable, maxAge, secret)

	case "redis":
		return NewRedisStore(
			cfg.RedisHost,
			cfg.RedisPort,
			cfg.RedisPassword,
			cfg.RedisDb,
			maxAge,
			secret,
		)

	case "cookie":
		cookieStore := sessions.NewCookieStore(secret)
		cookieStore.Options = &sessions.Options{
			Path:     cfg.WebCookiePath,
			Domain:   cfg.WebCookieDomain,
			MaxAge:   maxAge,
			Secure:   cfg.WebCookieSecure,
			HttpOnly: cfg.WebCookieHTTPOnly,
		}
		return cookieStore, nil

	default:
		return NewFileStore(cfg.SessionPath, maxAge, secret)
	}
}
