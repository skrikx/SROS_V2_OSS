package cli_integration_test

import (
	"bytes"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func runTraceCLICommand(t *testing.T, root string, args ...string) string {
	t.Helper()
	var out, errBuf bytes.Buffer
	cmd := exec.Command("go", append([]string{"run", "./cmd/sros"}, args...)...)
	cmd.Dir = root
	cmd.Env = append(os.Environ(), "PATH="+os.Getenv("PATH"))
	cmd.Stdout = &out
	cmd.Stderr = &errBuf
	if err := cmd.Run(); err != nil {
		t.Fatalf("cli command failed: %v stderr=%s", err, errBuf.String())
	}
	return out.String()
}

func joinRoot(root string, parts ...string) string {
	items := append([]string{root}, parts...)
	return filepath.Join(items...)
}
