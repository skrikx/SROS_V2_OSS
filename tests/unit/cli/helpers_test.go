package cli_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func runCLI(t *testing.T, args []string, env []string) (int, string, string) {
	t.Helper()

	root := repoRoot(t)
	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmdArgs := append([]string{"run", "./cmd/sros"}, args...)
	cmd := exec.Command("go", cmdArgs...)
	cmd.Dir = root
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if len(env) == 0 {
		cmd.Env = os.Environ()
	} else {
		cmd.Env = append(os.Environ(), env...)
	}

	err := cmd.Run()
	if err == nil {
		return 0, stdout.String(), stderr.String()
	}
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode(), stdout.String(), stderr.String()
	}
	t.Fatalf("run cli: %v", err)
	return 1, stdout.String(), stderr.String()
}

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	root := filepath.Clean(filepath.Join(wd, "..", "..", ".."))
	return root
}

func loadGolden(t *testing.T, name string) string {
	t.Helper()
	path := filepath.Join(repoRoot(t), "tests", "golden", "cli", name)
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("read golden %s: %v", path, err)
	}
	return string(data)
}
