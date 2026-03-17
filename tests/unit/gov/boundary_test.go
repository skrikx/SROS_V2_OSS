package gov_test

import (
	"testing"

	"srosv2/contracts/policy"
	"srosv2/internal/core/gov"
)

func TestResolveBoundary(t *testing.T) {
	if got := gov.ResolveBoundary("patch.apply"); got != policy.TrustBoundaryLocalFS {
		t.Fatalf("expected local fs, got %s", got)
	}
	if got := gov.ResolveBoundary("connector.invoke"); got != policy.TrustBoundaryExternalNet {
		t.Fatalf("expected external net, got %s", got)
	}
}
