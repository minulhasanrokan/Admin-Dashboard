package session

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"admin-dashboard/cmd/internal/config"
	"admin-dashboard/cmd/internal/constants"
	"admin-dashboard/cmd/internal/dtos"

	"github.com/gorilla/sessions"
)

func TestSerializer(t *testing.T) {
	id, err := generateSessionID()
	if err != nil {
		t.Fatalf("failed to generate session id: %v", err)
	}
	if len(id) != 64 {
		t.Fatalf("expected 64 char session id, got %d chars: %s", len(id), id)
	}

	values := map[interface{}]interface{}{
		constants.SessionKeyUserID:   uint64(42),
		constants.SessionKeyUserName: "admin_user",
		constants.SessionKeyEmail:    "admin@example.com",
		constants.SessionKeyRouteActionDetails: dtos.RightInfo{
			GroupName:    "Admin",
			CategoryName: "Users",
			RightName:    "ViewUsers",
			RightAction:  "/users",
		},
	}

	encoded, err := encodeSession(values)
	if err != nil {
		t.Fatalf("failed to encode session: %v", err)
	}

	var decoded map[interface{}]interface{}
	if err := decodeSession(encoded, &decoded); err != nil {
		t.Fatalf("failed to decode session: %v", err)
	}

	if decoded[constants.SessionKeyUserID] != uint64(42) {
		t.Errorf("expected user id 42, got %v", decoded[constants.SessionKeyUserID])
	}
	if decoded[constants.SessionKeyUserName] != "admin_user" {
		t.Errorf("expected username admin_user, got %v", decoded[constants.SessionKeyUserName])
	}

	rightInfo, ok := decoded[constants.SessionKeyRouteActionDetails].(dtos.RightInfo)
	if !ok || rightInfo.GroupName != "Admin" || rightInfo.RightName != "ViewUsers" {
		t.Errorf("unexpected right info: %+v", decoded[constants.SessionKeyRouteActionDetails])
	}

	userID := extractUserID(values)
	if userID == nil || *userID != 42 {
		t.Errorf("expected extracted user id 42, got %v", userID)
	}
}

func TestFileStore(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "session_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	secret := []byte("test-secret-key-32-bytes-long!!")
	store, err := NewFileStore(tempDir, 3600, secret)
	if err != nil {
		t.Fatalf("failed to create file store: %v", err)
	}

	sessionName := "test-session"

	req1 := httptest.NewRequest("GET", "http://localhost/", nil)
	rec1 := httptest.NewRecorder()

	sess1, err := store.New(req1, sessionName)
	if err != nil {
		t.Fatalf("store.New failed: %v", err)
	}
	if !sess1.IsNew {
		t.Errorf("expected session to be new")
	}

	sess1.Values[constants.SessionKeyUserID] = uint64(100)
	sess1.Values[constants.SessionKeyUserName] = "john_doe"

	if err := store.Save(req1, rec1, sess1); err != nil {
		t.Fatalf("store.Save failed: %v", err)
	}

	resp1 := rec1.Result()
	cookies := resp1.Cookies()
	if len(cookies) == 0 {
		t.Fatalf("no cookie set in response")
	}
	sessionCookie := cookies[0]
	if sessionCookie.Name != sessionName {
		t.Errorf("expected cookie name %s, got %s", sessionName, sessionCookie.Name)
	}

	expectedFile := filepath.Join(tempDir, "sess_"+sess1.ID)
	if _, err := os.Stat(expectedFile); err != nil {
		t.Errorf("expected session file %s to exist: %v", expectedFile, err)
	}

	req2 := httptest.NewRequest("GET", "http://localhost/", nil)
	req2.AddCookie(sessionCookie)

	sess2, err := store.New(req2, sessionName)
	if err != nil {
		t.Fatalf("store.New with cookie failed: %v", err)
	}
	if sess2.IsNew {
		t.Errorf("expected session to NOT be new when loaded from cookie")
	}
	if sess2.Values[constants.SessionKeyUserID] != uint64(100) {
		t.Errorf("expected user id 100, got %v", sess2.Values[constants.SessionKeyUserID])
	}
	if sess2.Values[constants.SessionKeyUserName] != "john_doe" {
		t.Errorf("expected username john_doe, got %v", sess2.Values[constants.SessionKeyUserName])
	}

	rec3 := httptest.NewRecorder()
	sess2.Options.MaxAge = -1
	if err := store.Save(req2, rec3, sess2); err != nil {
		t.Fatalf("failed to delete session: %v", err)
	}

	if _, err := os.Stat(expectedFile); !os.IsNotExist(err) {
		t.Errorf("expected session file %s to be deleted after MaxAge = -1", expectedFile)
	}

	resp3 := rec3.Result()
	var deletedCookie *http.Cookie
	for _, c := range resp3.Cookies() {
		if c.Name == sessionName {
			deletedCookie = c
			break
		}
	}
	if deletedCookie == nil || deletedCookie.MaxAge > 0 {
		t.Errorf("expected deleted cookie with MaxAge <= 0")
	}
}

func TestFileStoreExpiration(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "session_exp_test_*")
	if err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	secret := []byte("test-secret-key-32-bytes-long!!")
	store, err := NewFileStore(tempDir, 1, secret)
	if err != nil {
		t.Fatalf("failed to create file store: %v", err)
	}

	sessionName := "exp-session"

	req := httptest.NewRequest("GET", "http://localhost/", nil)
	rec := httptest.NewRecorder()

	sess, _ := store.New(req, sessionName)
	sess.Options.MaxAge = 1
	sess.Values["key"] = "value"
	if err := store.Save(req, rec, sess); err != nil {
		t.Fatalf("store.Save failed: %v", err)
	}

	cookies := rec.Result().Cookies()
	if len(cookies) == 0 {
		t.Fatalf("no cookie set in response")
	}
	cookie := cookies[0]

	time.Sleep(1500 * time.Millisecond)

	req2 := httptest.NewRequest("GET", "http://localhost/", nil)
	req2.AddCookie(cookie)

	sess2, _ := store.New(req2, sessionName)
	if !sess2.IsNew {
		t.Errorf("expected expired session to be treated as new session")
	}
	if sess2.Values["key"] != nil {
		t.Errorf("expected expired session values to be empty, got: %v", sess2.Values["key"])
	}
}

func TestSessionManagerFactory(t *testing.T) {
	cfg := &config.Config{
		AppKey:          "test-secret-app-key-1234567890!",
		SessionDriver:   "file",
		SessionPath:     t.TempDir(),
		WebCookieMaxAge: 3600,
	}

	store, err := NewSessionStore(cfg, nil)
	if err != nil {
		t.Fatalf("failed to create file store via factory: %v", err)
	}
	if _, ok := store.(*FileStore); !ok {
		t.Errorf("expected *FileStore, got %T", store)
	}

	cfg.SessionDriver = "cookie"
	storeCookie, err := NewSessionStore(cfg, nil)
	if err != nil {
		t.Fatalf("failed to create cookie store via factory: %v", err)
	}
	if _, ok := storeCookie.(*sessions.CookieStore); !ok {
		t.Errorf("expected *sessions.CookieStore, got %T", storeCookie)
	}

	cfg.SessionDriver = "unknown_driver"
	storeFallback, err := NewSessionStore(cfg, nil)
	if err != nil {
		t.Fatalf("failed to fallback to file store: %v", err)
	}
	if _, ok := storeFallback.(*FileStore); !ok {
		t.Errorf("expected fallback *FileStore, got %T", storeFallback)
	}

	cfg.SessionDriver = "database"
	_, err = NewSessionStore(cfg, nil)
	if err == nil {
		t.Errorf("expected error when initializing database store with nil db, got nil")
	}
}

func TestFileStoreCleanExpired(t *testing.T) {
	tempDir := t.TempDir()
	secret := []byte("test-secret-key-32-bytes-long!!")
	store, err := NewFileStore(tempDir, 1, secret)
	if err != nil {
		t.Fatalf("failed to create file store: %v", err)
	}

	req := httptest.NewRequest("GET", "http://localhost/", nil)
	rec := httptest.NewRecorder()

	sess, _ := store.New(req, "test-clean")
	sess.Options.MaxAge = 1
	sess.Values["temp"] = "data"
	if err := store.Save(req, rec, sess); err != nil {
		t.Fatalf("store.Save failed: %v", err)
	}

	cleaned := store.CleanExpired()
	if cleaned != 0 {
		t.Errorf("expected 0 cleaned files, got %d", cleaned)
	}

	time.Sleep(1500 * time.Millisecond)

	cleaned = store.CleanExpired()
	if cleaned != 1 {
		t.Errorf("expected 1 cleaned file, got %d", cleaned)
	}
}
