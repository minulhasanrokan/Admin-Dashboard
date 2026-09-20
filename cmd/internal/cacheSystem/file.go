package cacheSystem

import (
	"admin-dashboard/cmd/common/helpers"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type FileCache struct {
	dir string
	mu  sync.RWMutex
	ttl time.Duration
}

func NewFileCache(dir string, ttl time.Duration) (*FileCache, error) {
	if dir == "" {
		return nil, fmt.Errorf("cache directory cannot be empty")
	}

	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create cache directory '%s': %w", dir, err)
	}

	return &FileCache{
		dir: dir,
		ttl: ttl * time.Minute,
	}, nil
}

func (f *FileCache) file(key string) string {
	return filepath.Join(f.dir, key+".cache")
}

func (f *FileCache) Set(key string, value any) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	b, err := helpers.Marshal(value)
	if err != nil {
		return err
	}

	item := Item{
		Value: b,
	}

	if f.ttl > 0 {
		item.Expiration = time.Now().Add(f.ttl)
	}

	data, err := json.Marshal(item)
	if err != nil {
		return err
	}

	return os.WriteFile(f.file(key), data, 0644)
}

func (f *FileCache) Get(key string, dest any) error {

	f.mu.RLock()
	defer f.mu.RUnlock()

	data, err := os.ReadFile(f.file(key))
	if err != nil {
		return ErrCacheKeyNotFound
	}

	var item Item

	err = helpers.Unmarshal(data, &item)

	if err != nil {
		return err
	}

	if item.Expired() {

		err := os.Remove(f.file(key))

		if err != nil {
			return err
		}

		return ErrCacheExpired
	}

	return helpers.Unmarshal(item.Value, dest)
}

func (f *FileCache) Delete(key string) error {
	f.mu.Lock()
	defer f.mu.Unlock()

	files, err := filepath.Glob(filepath.Join(f.dir, key+"*.cache"))
	if err != nil {
		return err
	}

	for _, file := range files {
		if err := os.Remove(file); err != nil && !os.IsNotExist(err) {
			return err
		}
	}

	return nil
}

func (f *FileCache) Exists(key string) bool {
	f.mu.RLock()
	defer f.mu.RUnlock()

	_, err := os.Stat(f.file(key))

	return err == nil
}

func (f *FileCache) Flush() error {
	f.mu.Lock()
	defer f.mu.Unlock()

	files, _ := os.ReadDir(f.dir)

	for _, file := range files {

		err := os.Remove(filepath.Join(f.dir, file.Name()))

		if err != nil {
			return err
		}
	}

	return nil
}

func (f *FileCache) Close() error {
	return nil
}
