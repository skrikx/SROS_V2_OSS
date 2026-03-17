package regression_test

import (
	"os"
	"strings"
	"testing"
)

func TestFrontDoorRegression(t *testing.T) {
	data, err := os.ReadFile("../../README.md")
	if err != nil {
		t.Fatalf("read README: %v", err)
	}
	readme := string(data)
	for _, token := range []string{"Why try it", "Fast path", "examples/showcase", "scripts/first_run_smoke.sh"} {
		if !strings.Contains(readme, token) {
			t.Fatalf("README missing token %q", token)
		}
	}
}
