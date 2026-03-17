package gov

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"srosv2/contracts/policy"
)

type AuditRecord struct {
	Decision policy.PolicyDecision `json:"decision"`
	Recorded time.Time             `json:"recorded"`
}

func writeAudit(root string, decision policy.PolicyDecision, now time.Time) error {
	if root == "" {
		return nil
	}
	if err := os.MkdirAll(root, 0o755); err != nil {
		return fmt.Errorf("create gov audit root: %w", err)
	}
	record := AuditRecord{Decision: decision, Recorded: now.UTC()}
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal audit record: %w", err)
	}
	path := filepath.Join(root, string(decision.DecisionID)+".json")
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("write audit record: %w", err)
	}
	return nil
}
