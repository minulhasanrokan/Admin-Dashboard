package services_test

import (
	"admin-dashboard/cmd/app/repositories"
	"admin-dashboard/cmd/app/requests"
	"admin-dashboard/cmd/app/services"
	"testing"
)

func TestDesignationRequestValidationAndGetters(t *testing.T) {
	createReq := requests.CreateDesignationRequest{
		DesignationName:  "Lead Architect",
		DesignationTitle: "Enterprise Lead Architect",
		DesignationCode:  "ARCH-01",
		DesignationDetails: "Architecture design lead",
	}

	if createReq.GetName() != "Lead Architect" {
		t.Errorf("expected Lead Architect, got: %s", createReq.GetName())
	}
	if createReq.GetTitle() != "Enterprise Lead Architect" {
		t.Errorf("expected Enterprise Lead Architect, got: %s", createReq.GetTitle())
	}
	if createReq.GetCode() != "ARCH-01" {
		t.Errorf("expected ARCH-01, got: %s", createReq.GetCode())
	}
	if createReq.GetDetails() != "Architecture design lead" {
		t.Errorf("expected Architecture design lead, got: %s", createReq.GetDetails())
	}

	updateReq := requests.UpdateDesignationRequest{
		Name:    "Senior Manager",
		Title:   "Branch Senior Manager",
		Code:    "MGR-02",
		Details: "Operations manager",
	}

	if updateReq.GetName() != "Senior Manager" {
		t.Errorf("expected Senior Manager, got: %s", updateReq.GetName())
	}
	if updateReq.GetTitle() != "Branch Senior Manager" {
		t.Errorf("expected Branch Senior Manager, got: %s", updateReq.GetTitle())
	}
	if updateReq.GetCode() != "MGR-02" {
		t.Errorf("expected MGR-02, got: %s", updateReq.GetCode())
	}
	if updateReq.GetDetails() != "Operations manager" {
		t.Errorf("expected Operations manager, got: %s", updateReq.GetDetails())
	}

	statusReq := requests.StatusReasonDesignationRequest{
		ActiveReason: "Board approved restoration",
	}
	if statusReq.GetReason() != "Board approved restoration" {
		t.Errorf("expected Board approved restoration, got: %s", statusReq.GetReason())
	}
}

func TestDesignationServiceWithNilDb(t *testing.T) {
	svc := &services.Service{
		Db: nil,
	}

	// GetDesignationList should be safe with nil/empty DB
	res, err := svc.GetDesignationList(repositories.BackupFilterParams{Page: 1, Limit: 10})
	if err != nil {
		t.Fatalf("unexpected error with nil db: %v", err)
	}
	if res == nil {
		t.Fatalf("expected non-nil result map")
	}

	// Create without required name should error
	_, err = svc.CreateDesignation("", "Title", "CODE", "")
	if err == nil {
		t.Errorf("expected error for empty designation name")
	}

	// Create without required code should error
	_, err = svc.CreateDesignation("Manager", "Title", "", "")
	if err == nil {
		t.Errorf("expected error for empty designation code")
	}
}
