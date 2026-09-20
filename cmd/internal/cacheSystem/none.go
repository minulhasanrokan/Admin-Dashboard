package cacheSystem

type NoCache struct{}

func NewNoCache() (CacheSystem, error) {
	return &NoCache{}, nil
}

func (n *NoCache) Set(key string, value any) error {

	_ = key
	_ = value
	return nil
}

func (n *NoCache) Get(key string, dest any) error {

	_ = key
	_ = dest
	return ErrCacheKeyNotFound
}

func (n *NoCache) Delete(key string) error {

	_ = key
	return nil
}

func (n *NoCache) Exists(key string) bool {
	_ = key
	return false
}

func (n *NoCache) Flush() error {
	return nil
}

func (n *NoCache) Close() error {
	return nil
}
