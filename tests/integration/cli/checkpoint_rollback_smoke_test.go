package cli_integration_test

import (
	"encoding/json"
	"path/filepath"
	"testing"
)

type checkpointResponse struct {
	Session struct {
		SessionID string `json:"session_id"`
	} `json:"session"`
	CheckpointID string `json:"checkpoint_id"`
}

func TestCheckpointRollbackSmoke(t *testing.T) {
	root := repoRoot(t)
	contract := filepath.Join(root, "examples", "runs", "checkpointable_run_contract.json")

	runOut := runCLICmdJSON(t, root, "run", "--contract", contract)
	var runResp runtimeCLIResponse
	if err := json.Unmarshal(runOut, &runResp); err != nil {
		t.Fatalf("decode run response: %v", err)
	}

	cpOut := runCLICmdJSON(t, root, "checkpoint", "--session", runResp.Session.SessionID, "--stage", "validated")
	var cpResp checkpointResponse
	if err := json.Unmarshal(cpOut, &cpResp); err != nil {
		t.Fatalf("decode checkpoint response: %v", err)
	}
	if cpResp.CheckpointID == "" {
		t.Fatal("expected checkpoint id")
	}

	runCLICmd(t, root, "rollback", "--session", runResp.Session.SessionID, "--checkpoint", cpResp.CheckpointID)
}
