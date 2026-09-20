package services_test

import (
	"testing"

	"admin-dashboard/cmd/app/services"
	"admin-dashboard/cmd/internal/config"
)

func TestDashboardDataRetrieval(t *testing.T) {
	// Verify that with nil DB (e.g. disconnected), GetDashboardData executes safely without panics
	svc := &services.Service{
		Db:     nil,
		Config: config.Config{},
	}

	data, err := svc.GetDashboardData()
	if err != nil {
		t.Fatalf("expected nil error with empty db, got: %v", err)
	}

	if data.Stats.TotalUsers < 0 {
		t.Errorf("expected non-negative TotalUsers, got: %d", data.Stats.TotalUsers)
	}

	if data.RecentActivities == nil {
		t.Errorf("expected non-nil RecentActivities slice")
	}

	if data.RecentLogins == nil {
		t.Errorf("expected non-nil RecentLogins slice")
	}
}
