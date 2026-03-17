package cli_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestConfigUsesExplicitPath(t *testing.T) {
	root := repoRoot(t)
	tmp := t.TempDir()
	cfgPath := filepath.Join(tmp, "sros.yaml")
	content := "mode: local_cli\nworkspace_root: ./workspace\nartifact_root: ./workspace/artifacts\npolicy_bundle_path: ./workspace/contracts/policy/local.bundle.json\nmemory_store_path: ./workspace/artifacts/memory\ntrace_store_path: ./workspace/artifacts/trace\noutput_format: text\n"
	if err := os.WriteFile(cfgPath, []byte(content), 0o644); err != nil {
		t.Fatalf("write config: %v", err)
	}

	code, stdout, stderr := runCLI(t, []string{"--config", cfgPath, "config"}, []string{"PATH=" + os.Getenv("PATH")})
	if code != 0 {
		t.Fatalf("expected exit 0, got %d stderr=%s", code, stderr)
	}

	if !strings.Contains(stdout, filepath.Join(tmp, "workspace")) {
		t.Fatalf("expected workspace path from config, got:\n%s", stdout)
	}
	_ = root
}

func TestConfigEnvOverrideWins(t *testing.T) {
	override := filepath.Join(t.TempDir(), "override-workspace")
	env := []string{
		"SROS_WORKSPACE_ROOT=" + override,
		"SROS_OUTPUT_FORMAT=text",
		"PATH=" + os.Getenv("PATH"),
	}
	code, stdout, stderr := runCLI(t, []string{"config"}, env)
	if code != 0 {
		t.Fatalf("expected exit 0, got %d stderr=%s", code, stderr)
	}
	if !strings.Contains(stdout, override) {
		t.Fatalf("expected env override workspace path, got:\n%s", stdout)
	}
}
