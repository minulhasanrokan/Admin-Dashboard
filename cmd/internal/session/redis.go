package session

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client  *redis.Client
	Codecs  []securecookie.Codec
	Options *sessions.Options
	prefix  string
}

func NewRedisStore(host, port, password, db string, maxAge int, keyPairs ...[]byte) (*RedisStore, error) {
	dbIndex := 0
	if db != "" {
		n, err := strconv.Atoi(db)
		if err != nil {
			return nil, fmt.Errorf("invalid redis db '%s': %w", db, err)
		}
		dbIndex = n
	}

	if host == "" {
		host = "127.0.0.1"
	}
	if port == "" {
		port = "6379"
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

	return NewRedisStoreWithClient(rdb, maxAge, keyPairs...), nil
}

func NewRedisStoreWithClient(client *redis.Client, maxAge int, keyPairs ...[]byte) *RedisStore {
	if maxAge <= 0 {
		maxAge = 86400 * 30
	}

	return &RedisStore{
		client:  client,
		Codecs:  securecookie.CodecsFromPairs(keyPairs...),
		prefix:  "session:",
		Options: &sessions.Options{
			Path:     "/",
			MaxAge:   maxAge,
			HttpOnly: true,
		},
	}
}

func (s *RedisStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

func (s *RedisStore) New(r *http.Request, name string) (*sessions.Session, error) {
	session := sessions.NewSession(s, name)
	opts := *s.Options
	session.Options = &opts
	session.IsNew = true

	c, err := r.Cookie(name)
	if err != nil {
		return session, nil
	}

	if err = securecookie.DecodeMulti(name, c.Value, &session.ID, s.Codecs...); err != nil {
		return session, nil
	}

	if err = s.load(session); err != nil {
		session.ID = ""
		session.Values = make(map[interface{}]interface{})
		session.IsNew = true
		return session, nil
	}

	session.IsNew = false
	return session, nil
}

func (s *RedisStore) load(session *sessions.Session) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	data, err := s.client.Get(ctx, s.prefix+session.ID).Bytes()
	if errors.Is(err, redis.Nil) {
		return fmt.Errorf("session key not found")
	}
	if err != nil {
		return err
	}

	var values map[interface{}]interface{}
	if err := decodeSession(data, &values); err != nil {
		return fmt.Errorf("failed to deserialize session: %w", err)
	}

	session.Values = values
	return nil
}

func (s *RedisStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	_ = r
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if session.Options.MaxAge <= 0 {
		if session.ID != "" {
			_ = s.client.Del(ctx, s.prefix+session.ID).Err()
		}
		http.SetCookie(w, sessions.NewCookie(session.Name(), "", session.Options))
		return nil
	}

	if session.ID == "" {
		id, err := generateSessionID()
		if err != nil {
			return err
		}
		session.ID = id
	}

	encodedBytes, err := encodeSession(session.Values)
	if err != nil {
		return fmt.Errorf("failed to encode session: %w", err)
	}

	ttl := time.Duration(session.Options.MaxAge) * time.Second
	if err := s.client.Set(ctx, s.prefix+session.ID, encodedBytes, ttl).Err(); err != nil {
		return fmt.Errorf("failed to save session to redis: %w", err)
	}

	encodedCookie, err := securecookie.EncodeMulti(session.Name(), session.ID, s.Codecs...)
	if err != nil {
		return fmt.Errorf("failed to encode session cookie: %w", err)
	}

	http.SetCookie(w, sessions.NewCookie(session.Name(), encodedCookie, session.Options))
	return nil
}

func (s *RedisStore) Close() error {
	return s.client.Close()
}
