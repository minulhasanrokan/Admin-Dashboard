package cacheSystem

type CacheSystem interface {
	Set(key string, value any) error
	Get(key string, dest any) error
	Delete(key string) error
	Exists(key string) bool
	Flush() error
	Close() error
}
