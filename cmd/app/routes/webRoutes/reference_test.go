package webRoute

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestReferenceDesignationRoutesSecurity(t *testing.T) {
	e, _ := setupTestApp()

	endpoints := []string{
		"/reference/designation",
		"/reference/designation/add",
		"/reference/designation/view",
		"/reference/designation/view/test-uuid",
		"/reference/designation/edit",
		"/reference/designation/edit/test-uuid",
		"/reference/designation/delete",
		"/reference/designation/delete/test-uuid",
		"/reference/designation/stc",
		"/reference/designation/stc/test-uuid",
		"/reference/designation/verify",
		"/reference/designation/verify/test-uuid",
		"/reference/designation/active",
		"/reference/designation/active/test-uuid",
		"/reference/designation/deactivate",
		"/reference/designation/deactivate/test-uuid",
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
