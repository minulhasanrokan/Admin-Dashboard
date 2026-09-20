package webRoute

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestActivityLogsAndLoginHistoryRoutesSecurity(t *testing.T) {
	e, _ := setupTestApp()

	endpoints := []string{
		"/settings/activity-logs",
		"/settings/activity-logs/list",
		"/settings/login-history",
		"/settings/login-history/list",
	}

	for _, endpoint := range endpoints {
		req := httptest.NewRequest(http.MethodGet, endpoint, nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		// Unauthenticated requests should be blocked or redirected
		if rec.Code == http.StatusOK {
			t.Errorf("expected unauthorized or error status for %s, got %d", endpoint, rec.Code)
		}

		// Security headers should still be present
		csp := rec.Header().Get("Content-Security-Policy")
		if csp == "" {
			t.Errorf("missing Content-Security-Policy header on %s response", endpoint)
		}
	}
}
