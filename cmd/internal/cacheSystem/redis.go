package cacheSystem

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisCache struct {
	client *redis.Client
	ctx    context.Context
	ttl    time.Duration
}

func NewRedisCache(host, port, password, db string, ttl time.Duration) (*RedisCache, error) {
	dbIndex := 0

	if db != "" {
		n, err := strconv.Atoi(db)
		if err != nil {
			return nil, fmt.Errorf("invalid redis db '%s': %w", db, err)
		}
		dbIndex = n
	}

	rdb := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       dbIndex,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	if err := rdb.Ping(ctx).Err(); err != nil {
		_ = rdb.Close()
		return nil, fmt.Errorf("failed to connect to redis (%s:%s): %w", host, port, err)
	}

	return &RedisCache{
		client: rdb,
		ctx:    context.Background(),
		ttl:    ttl * time.Minute,
	}, nil
}

func (r *RedisCache) Set(key string, value any) error {
	b, err := json.Marshal(value)
	if err != nil {
		return err
	}

	return r.client.Set(r.ctx, key, b, r.ttl).Err()
}

func (r *RedisCache) Get(key string, dest any) error {

	val, err := r.client.Get(r.ctx, key).Result()

	if errors.Is(err, redis.Nil) {

		_ = r.client.Del(r.ctx, key).Err()
		return ErrCacheKeyNotFound
	}

	if err != nil {
		return err
	}

	if err := json.Unmarshal([]byte(val), dest); err != nil {
		return err
	}

	return nil
}

func (r *RedisCache) Delete(key string) error {

	keys, err := r.client.Keys(r.ctx, key+"*").Result()

	if err != nil {
		return err
	}

	if len(keys) == 0 {
		return nil
	}

	return r.client.Del(r.ctx, keys...).Err()
}

func (r *RedisCache) Exists(key string) bool {
	n, _ := r.client.Exists(r.ctx, key).Result()
	return n > 0
}

func (r *RedisCache) Flush() error {
	return r.client.FlushDB(r.ctx).Err()
}

func (r *RedisCache) Close() error {
	return r.client.Close()
}
