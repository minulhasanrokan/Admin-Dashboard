package webRoute

import (
	"admin-dashboard/cmd/app/handlers"
	"admin-dashboard/cmd/internal/applicationData"
	"admin-dashboard/cmd/internal/config"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/labstack/echo/v4"
)

func setupTestApp() (*echo.Echo, *config.Config) {
	e := echo.New()
	cfg := &config.Config{
		WebRateLimit:             100,
		WebXSSProtection:         "1; mode=block",
		WebContentTypeNosniff:    "nosniff",
		WebXFrameOptions:         "DENY",
		WebContentSecurityPolicy: "default-src 'none'; script-src 'self'; style-src 'self'; img-src 'self' data:; font-src 'self' data:; connect-src 'self'; object-src 'none'; frame-ancestors 'none'; form-action 'self'; base-uri 'self'; frame-src 'none'; media-src 'self'; manifest-src 'self';",
		WebReferrerPolicy:        "strict-origin-when-cross-origin",
		WebCookieName:            "_csrf",
		WebCookiePath:            "/",
		WebCookieHTTPOnly:        true,
		WebCookieSecure:          false,
		WebCookieSameSite:        "strict",
		WebCookieMaxAge:          3600,
		WebCSRFTokenLookup:       "header:X-CSRF-Token,form:_csrf",
		WebCSRFTokenLength:       32,
	}

	h := handlers.Handler{
		Config: *cfg,
	}

	appData := &applicationData.Application{
		Handler: h,
	}

	g := e.Group("/")
	MainRoute(g, appData)

	return e, cfg
}

func TestSecurityHeaders(t *testing.T) {
	e, _ := setupTestApp()

	req := httptest.NewRequest(http.MethodGet, "/health-check", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200 OK, got %d", rec.Code)
	}

	// Verify CSP
	// Verify Strict CSP
	csp := rec.Header().Get("Content-Security-Policy")
	if csp == "" {
		t.Errorf("missing Content-Security-Policy header")
	}
	if !strings.Contains(csp, "default-src 'none'") {
		t.Errorf("Strict CSP must have default-src 'none', got: %s", csp)
	}
	if strings.Contains(csp, "'unsafe-inline'") {
		t.Errorf("Strict CSP must NOT contain 'unsafe-inline', got: %s", csp)
	}
	if !strings.Contains(csp, "script-src 'self'") {
		t.Errorf("Strict CSP missing script-src 'self'")
	}
	if !strings.Contains(csp, "style-src 'self'") {
		t.Errorf("Strict CSP missing style-src 'self'")
	}
	if !strings.Contains(csp, "object-src 'none'") {
		t.Errorf("Strict CSP missing object-src 'none'")
	}
	if !strings.Contains(csp, "frame-ancestors 'none'") {
		t.Errorf("Strict CSP missing frame-ancestors 'none'")
	}
	if !strings.Contains(csp, "form-action 'self'") {
		t.Errorf("Strict CSP missing form-action 'self'")
	}
	if !strings.Contains(csp, "base-uri 'self'") {
		t.Errorf("Strict CSP missing base-uri 'self'")
	}

	// Verify X-Frame-Options
	xFrame := rec.Header().Get("X-Frame-Options")
	if xFrame != "DENY" {
		t.Errorf("expected X-Frame-Options: DENY, got: %s", xFrame)
	}

	// Verify X-Content-Type-Options
	nosniff := rec.Header().Get("X-Content-Type-Options")
	if nosniff != "nosniff" {
		t.Errorf("expected X-Content-Type-Options: nosniff, got: %s", nosniff)
	}

	// Verify Referrer-Policy
	referrer := rec.Header().Get("Referrer-Policy")
	if referrer != "strict-origin-when-cross-origin" {
		t.Errorf("expected Referrer-Policy: strict-origin-when-cross-origin, got: %s", referrer)
	}
}

func TestCSRFProtectionRejectsPostWithoutToken(t *testing.T) {
	e, _ := setupTestApp()

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader("username=test&password=secret"))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)

	// Since CSRF token is missing, request must be rejected (403 Forbidden or custom error)
	if rec.Code == http.StatusOK {
		t.Errorf("expected CSRF rejection for POST request without token, but got 200 OK")
	}
}
