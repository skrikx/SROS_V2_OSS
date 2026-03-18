package e2e_test

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"
)

func repoRoot(t *testing.T) string {
	t.Helper()
	wd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	return filepath.Clean(filepath.Join(wd, "..", "..", ".."))
}

func bashScriptCommand(t *testing.T, script string) *exec.Cmd {
	t.Helper()

	const gitBashPath = "C:\\Program Files\\Git\\bin\\bash.exe"
	if runtime.GOOS == "windows" {
		if _, err := os.Stat(gitBashPath); err == nil {
			return exec.Command(gitBashPath, script)
		}
	}

	if bashPath, err := exec.LookPath("bash"); err == nil {
		return exec.Command(bashPath, script)
	}

	t.Fatalf("bash executable not found in PATH or at %q", gitBashPath)
	return nil
}
