package session

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
)

func init() {
	gob.Register(fileSessionEnvelope{})
}

type fileSessionEnvelope struct {
	ExpiresAt int64
	Payload   []byte
}

type FileStore struct {
	Codecs  []securecookie.Codec
	Options *sessions.Options
	dir     string
	mu      sync.RWMutex
}

func NewFileStore(path string, maxAge int, keyPairs ...[]byte) (*FileStore, error) {
	if path == "" {
		path = "storage/sessions"
	}

	if err := os.MkdirAll(path, 0755); err != nil {
		return nil, fmt.Errorf("failed to create session directory '%s': %w", path, err)
	}

	if maxAge <= 0 {
		maxAge = 86400 * 30
	}

	return &FileStore{
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &sessions.Options{
			Path:     "/",
			MaxAge:   maxAge,
			HttpOnly: true,
		},
		dir: path,
	}, nil
}

func (s *FileStore) filePath(sessionID string) string {
	return filepath.Join(s.dir, "sess_"+filepath.Base(sessionID))
}

func (s *FileStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

func (s *FileStore) New(r *http.Request, name string) (*sessions.Session, error) {
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

func (s *FileStore) load(session *sessions.Session) error {
	s.mu.RLock()
	defer s.mu.RUnlock()

	fp := s.filePath(session.ID)
	data, err := os.ReadFile(fp)
	if err != nil {
		return err
	}

	var env fileSessionEnvelope
	if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&env); err != nil {
		return err
	}

	if env.ExpiresAt > 0 && time.Now().Unix() >= env.ExpiresAt {
		_ = os.Remove(fp)
		return fmt.Errorf("session expired")
	}

	var values map[interface{}]interface{}
	if err := decodeSession(env.Payload, &values); err != nil {
		return err
	}

	session.Values = values
	return nil
}

func (s *FileStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	_ = r

	if session.Options.MaxAge <= 0 {
		if session.ID != "" {
			s.mu.Lock()
			_ = os.Remove(s.filePath(session.ID))
			s.mu.Unlock()
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

	encodedPayload, err := encodeSession(session.Values)
	if err != nil {
		return fmt.Errorf("failed to encode session values: %w", err)
	}

	expiresAt := time.Now().Add(time.Duration(session.Options.MaxAge) * time.Second).Unix()
	env := fileSessionEnvelope{
		ExpiresAt: expiresAt,
		Payload:   encodedPayload,
	}

	var envBuf bytes.Buffer
	if err := gob.NewEncoder(&envBuf).Encode(env); err != nil {
		return fmt.Errorf("failed to encode session envelope: %w", err)
	}
	encodedData := envBuf.Bytes()

	s.mu.Lock()
	fp := s.filePath(session.ID)
	tmpFile := fp + ".tmp"
	if err := os.WriteFile(tmpFile, encodedData, 0600); err != nil {
		s.mu.Unlock()
		return fmt.Errorf("failed to write session file: %w", err)
	}
	if err := os.Rename(tmpFile, fp); err != nil {
		_ = os.Remove(tmpFile)
		if err := os.WriteFile(fp, encodedData, 0600); err != nil {
			s.mu.Unlock()
			return fmt.Errorf("failed to save session file: %w", err)
		}
	}
	s.mu.Unlock()

	encodedCookie, err := securecookie.EncodeMulti(session.Name(), session.ID, s.Codecs...)
	if err != nil {
		return fmt.Errorf("failed to encode session cookie: %w", err)
	}

	http.SetCookie(w, sessions.NewCookie(session.Name(), encodedCookie, session.Options))
	return nil
}

func (s *FileStore) CleanExpired() int {
	s.mu.Lock()
	defer s.mu.Unlock()

	entries, err := os.ReadDir(s.dir)
	if err != nil {
		return 0
	}

	now := time.Now().Unix()
	cleaned := 0

	for _, entry := range entries {
		if entry.IsDir() || !strings.HasPrefix(entry.Name(), "sess_") {
			continue
		}
		fp := filepath.Join(s.dir, entry.Name())
		data, err := os.ReadFile(fp)
		if err != nil {
			continue
		}
		var env fileSessionEnvelope
		if err := gob.NewDecoder(bytes.NewReader(data)).Decode(&env); err == nil {
			if env.ExpiresAt > 0 && now >= env.ExpiresAt {
				_ = os.Remove(fp)
				cleaned++
			}
		}
	}

	return cleaned
}
