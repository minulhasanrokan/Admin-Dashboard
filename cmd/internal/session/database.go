package session

import (
	"encoding/base64"
	"errors"
	"fmt"
	"net/http"
	"time"

	"admin-dashboard/cmd/internal/models"

	"github.com/gorilla/securecookie"
	"github.com/gorilla/sessions"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

type DatabaseStore struct {
	db      *gorm.DB
	Codecs  []securecookie.Codec
	Options *sessions.Options
	table   string
}

func NewDatabaseStore(db *gorm.DB, table string, maxAge int, keyPairs ...[]byte) (*DatabaseStore, error) {
	if db == nil {
		return nil, errors.New("database connection cannot be nil for database session driver")
	}

	if table == "" {
		table = "sessions"
	}

	if !db.Migrator().HasTable(table) {
		if err := db.Table(table).AutoMigrate(&models.Session{}); err != nil {
			return nil, fmt.Errorf("failed to migrate session table '%s': %w", table, err)
		}
	}

	if maxAge <= 0 {
		maxAge = 86400 * 30
	}

	return &DatabaseStore{
		db:     db,
		Codecs: securecookie.CodecsFromPairs(keyPairs...),
		Options: &sessions.Options{
			Path:     "/",
			MaxAge:   maxAge,
			HttpOnly: true,
		},
		table: table,
	}, nil
}

func (s *DatabaseStore) Get(r *http.Request, name string) (*sessions.Session, error) {
	return sessions.GetRegistry(r).Get(s, name)
}

func (s *DatabaseStore) New(r *http.Request, name string) (*sessions.Session, error) {
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

func (s *DatabaseStore) load(session *sessions.Session) error {
	var record models.Session
	err := s.db.Table(s.table).
		Where("id = ? AND expires_at > ?", session.ID, time.Now()).
		First(&record).Error

	if err != nil {
		return err
	}

	payloadBytes, err := base64.StdEncoding.DecodeString(record.Payload)
	if err != nil {
		return fmt.Errorf("failed to decode session payload: %w", err)
	}

	var values map[interface{}]interface{}
	if err := decodeSession(payloadBytes, &values); err != nil {
		return fmt.Errorf("failed to deserialize session: %w", err)
	}

	session.Values = values
	return nil
}

func (s *DatabaseStore) Save(r *http.Request, w http.ResponseWriter, session *sessions.Session) error {
	if session.Options.MaxAge <= 0 {
		if session.ID != "" {
			_ = s.db.Table(s.table).Where("id = ?", session.ID).Delete(&models.Session{}).Error
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

	base64Payload := base64.StdEncoding.EncodeToString(encodedBytes)
	now := time.Now()
	expiresAt := now.Add(time.Duration(session.Options.MaxAge) * time.Second)
	userID := extractUserID(session.Values)
	ip := getClientIP(r)
	userAgent := r.UserAgent()

	record := models.Session{
		Id:           session.ID,
		UserId:       userID,
		IpAddress:    ip,
		UserAgent:    userAgent,
		Payload:      base64Payload,
		LastActivity: now.Unix(),
		ExpiresAt:    expiresAt,
		UpdatedAt:    now,
	}

	err = s.db.Table(s.table).
		Clauses(clause.OnConflict{
			Columns:   []clause.Column{{Name: "id"}},
			DoUpdates: clause.AssignmentColumns([]string{"user_id", "ip_address", "user_agent", "payload", "last_activity", "expires_at", "updated_at"}),
		}).
		Create(&record).Error

	if err != nil {
		return fmt.Errorf("failed to save session to database: %w", err)
	}

	encodedCookie, err := securecookie.EncodeMulti(session.Name(), session.ID, s.Codecs...)
	if err != nil {
		return fmt.Errorf("failed to encode session cookie: %w", err)
	}

	http.SetCookie(w, sessions.NewCookie(session.Name(), encodedCookie, session.Options))
	return nil
}

func (s *DatabaseStore) CleanExpired() (int64, error) {
	res := s.db.Table(s.table).Where("expires_at < ?", time.Now()).Delete(&models.Session{})
	return res.RowsAffected, res.Error
}
