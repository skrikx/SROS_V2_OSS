package mem_test

import (
	"testing"

	cmemory "srosv2/contracts/memory"
	"srosv2/internal/core/mem"
)

func TestScopeBindingValidate(t *testing.T) {
	binding := mem.ScopeBinding{Scope: cmemory.ScopeWorkspace, TenantID: "local", WorkspaceID: "default"}
	if err := binding.Validate(); err != nil {
		t.Fatalf("validate: %v", err)
	}
}
