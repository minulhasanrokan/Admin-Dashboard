package session

import (
	"bytes"
	"crypto/rand"
	"encoding/gob"
	"encoding/hex"
	"net"
	"net/http"
	"strings"

	"admin-dashboard/cmd/internal/constants"
)

func generateSessionID() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

func encodeSession(values map[interface{}]interface{}) ([]byte, error) {
	if len(values) == 0 {
		return []byte{}, nil
	}
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(values); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func decodeSession(data []byte, values *map[interface{}]interface{}) error {
	if len(data) == 0 {
		*values = make(map[interface{}]interface{})
		return nil
	}
	dec := gob.NewDecoder(bytes.NewReader(data))
	return dec.Decode(values)
}

func extractUserID(values map[interface{}]interface{}) *uint64 {
	if val, ok := values[constants.SessionKeyUserID]; ok {
		switch v := val.(type) {
		case uint64:
			return &v
		case uint:
			u64 := uint64(v)
			return &u64
		case int64:
			if v >= 0 {
				u64 := uint64(v)
				return &u64
			}
		case int:
			if v >= 0 {
				u64 := uint64(v)
				return &u64
			}
		}
	}
	return nil
}

func getClientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		if len(parts) > 0 {
			ip := strings.TrimSpace(parts[0])
			if net.ParseIP(ip) != nil {
				return ip
			}
		}
	}
	if xri := r.Header.Get("X-Real-IP"); xri != "" {
		ip := strings.TrimSpace(xri)
		if net.ParseIP(ip) != nil {
			return ip
		}
	}
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err == nil && net.ParseIP(ip) != nil {
		return ip
	}
	return r.RemoteAddr
}
