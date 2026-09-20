package config_test

import (
	"os"
	"strings"
	"testing"

	"admin-dashboard/cmd/internal/config"
)

func TestConfigDefaultStrictCSP(t *testing.T) {
	origCsp := os.Getenv("WEB_CONTENT_SECURITY_POLICY")
	defer os.Setenv("WEB_CONTENT_SECURITY_POLICY", origCsp)

	// Test fallback when env var is empty
	_ = os.Unsetenv("WEB_CONTENT_SECURITY_POLICY")
	cfg, err := config.Init()
	if err != nil {
		t.Fatalf("config.Init failed: %v", err)
	}

	if cfg.WebContentSecurityPolicy == "" {
		t.Fatal("expected default WebContentSecurityPolicy to be non-empty")
	}
	if !strings.Contains(cfg.WebContentSecurityPolicy, "default-src 'none'") {
		t.Errorf("expected default CSP to contain \"default-src 'none'\", got: %s", cfg.WebContentSecurityPolicy)
	}
	if strings.Contains(cfg.WebContentSecurityPolicy, "'unsafe-inline'") {
		t.Errorf("expected default CSP not to contain 'unsafe-inline', got: %s", cfg.WebContentSecurityPolicy)
	}

	// Test when env var is explicitly provided
	customCSP := "default-src 'none'; script-src 'self';"
	os.Setenv("WEB_CONTENT_SECURITY_POLICY", customCSP)
	cfgCustom, err := config.Init()
	if err != nil {
		t.Fatalf("config.Init failed: %v", err)
	}
	if cfgCustom.WebContentSecurityPolicy != customCSP {
		t.Errorf("expected custom CSP %q, got: %q", customCSP, cfgCustom.WebContentSecurityPolicy)
	}
}
