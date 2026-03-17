package orch

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

type CheckpointRoute struct {
	SessionID      string    `json:"session_id"`
	WorkUnitID     string    `json:"work_unit_id"`
	Route          string    `json:"route"`
	ApprovalPath   string    `json:"approval_path"`
	RequestedAt    time.Time `json:"requested_at"`
	Reason         string    `json:"reason"`
	Capability     string    `json:"capability"`
	SandboxProfile string    `json:"sandbox_profile"`
}

type CheckpointRouter struct {
	root string
}

func NewCheckpointRouter(root string) (*CheckpointRouter, error) {
	if root == "" {
		root = filepath.Join("artifacts", "runtime", "approvals")
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return nil, fmt.Errorf("create checkpoint route root: %w", err)
	}
	return &CheckpointRouter{root: root}, nil
}

func (r *CheckpointRouter) RouteAsk(route CheckpointRoute) (CheckpointRoute, error) {
	path := filepath.Join(r.root, route.SessionID+".json")
	route.ApprovalPath = path
	payload := map[string]any{
		"session_id":       route.SessionID,
		"work_unit_id":     route.WorkUnitID,
		"route":            route.Route,
		"approved":         false,
		"reason":           route.Reason,
		"capability":       route.Capability,
		"sandbox_profile":  route.SandboxProfile,
		"requested_at":     route.RequestedAt.UTC().Format(time.RFC3339),
		"operator_message": "set approved=true and resume explicitly to continue",
	}
	data, err := json.MarshalIndent(payload, "", "  ")
	if err != nil {
		return CheckpointRoute{}, fmt.Errorf("marshal checkpoint route: %w", err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return CheckpointRoute{}, fmt.Errorf("write checkpoint route: %w", err)
	}
	return route, nil
}
