package repositories

func (r *Repository) ClearCache() error {
	return r.CacheSystem.Flush()
}
