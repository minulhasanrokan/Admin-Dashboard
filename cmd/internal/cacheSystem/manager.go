package cacheSystem

import (
	"admin-dashboard/cmd/internal/config"
	"strings"
	"time"
)

func NewManager(cfg *config.Config) (CacheSystem, error) {
	driver := strings.ToLower(strings.TrimSpace(cfg.CacheDriver))

	switch driver {
	case "file":
		return NewFileCache(cfg.CachePath, time.Duration(cfg.CacheTime))

	case "redis":
		return NewRedisCache(
			cfg.RedisHost,
			cfg.RedisPort,
			cfg.RedisPassword,
			cfg.RedisDb,
			time.Duration(cfg.CacheTime),
		)

	default:
		return NewNoCache()
	}
}
