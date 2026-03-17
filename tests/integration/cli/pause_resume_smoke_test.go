package cli_integration_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

type runtimeCLIResponse struct {
	Session struct {
		SessionID string `json:"session_id"`
	} `json:"session"`
}

func TestPauseResumeSmoke(t *testing.T) {
	root := repoRoot(t)
	contract := filepath.Join(root, "examples", "runs", "minimal_run_contract.json")

	runOut := runCLICmdJSON(t, root, "run", "--contract", contract)
	var runResp runtimeCLIResponse
	if err := json.Unmarshal(runOut, &runResp); err != nil {
		t.Fatalf("decode run response: %v", err)
	}
	if runResp.Session.SessionID == "" {
		t.Fatal("expected session id from run response")
	}

	runCLICmd(t, root, "pause", "--session", runResp.Session.SessionID, "--reason", "integration smoke pause")
	runCLICmd(t, root, "resume", "--session", runResp.Session.SessionID)
}

func runCLICmd(t *testing.T, root string, args ...string) {
	t.Helper()
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmdArgs := append([]string{"run", "./cmd/sros"}, args...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "PATH="+os.Getenv("PATH"))
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		t.Fatalf("cli command failed: go %v: %v stderr=%s", cmdArgs, err, errOut.String())
	}
}

func runCLICmdJSON(t *testing.T, root string, args ...string) []byte {
	t.Helper()
	var out bytes.Buffer
	var errOut bytes.Buffer
	cmdArgs := append([]string{"run", "./cmd/sros", "--format", "json"}, args...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "PATH="+os.Getenv("PATH"))
	cmd.Stdout = &out
	cmd.Stderr = &errOut
	if err := cmd.Run(); err != nil {
		t.Fatalf("cli json command failed: go %v: %v stderr=%s", cmdArgs, err, errOut.String())
	}
	return out.Bytes()
}
