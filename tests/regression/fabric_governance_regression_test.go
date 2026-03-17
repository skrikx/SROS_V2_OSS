package regression_test

import (
	"os"
	"testing"
)

func TestFabricGovernanceRegressionFixturesExist(t *testing.T) {
	for _, path := range []string{
		"../../tests/golden/fabric/registry_search_result_min.json",
		"../../tests/golden/gov/allow_decision.json",
	} {
		if _, err := os.Stat(path); err != nil {
			t.Fatalf("missing fixture %s: %v", path, err)
		}
	}
}
