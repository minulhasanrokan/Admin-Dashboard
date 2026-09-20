package cacheSystem

import "errors"

var (
	ErrCacheKeyNotFound = errors.New("cache: key not found")
	ErrCacheExpired     = errors.New("cache: key expired")
)
