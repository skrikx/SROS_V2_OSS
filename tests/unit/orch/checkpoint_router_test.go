package orch_test

import (
	"os"
	"testing"
	"time"

	"srosv2/internal/core/orch"
)

func TestCheckpointRouterWritesApprovalArtifact(t *testing.T) {
	router, err := orch.NewCheckpointRouter(t.TempDir())
	if err != nil {
		t.Fatalf("new router: %v", err)
	}
	route, err := router.RouteAsk(orch.CheckpointRoute{
		SessionID:      "sess_001",
		WorkUnitID:     "wu-001",
		Route:          "local_cli_checkpoint",
		RequestedAt:    time.Date(2026, 3, 17, 12, 0, 0, 0, time.UTC),
		Reason:         "operator checkpoint required",
		Capability:     "connector.invoke",
		SandboxProfile: "net-observe",
	})
	if err != nil {
		t.Fatalf("route ask: %v", err)
	}
	if _, err := os.Stat(route.ApprovalPath); err != nil {
		t.Fatalf("approval artifact missing: %v", err)
	}
}
