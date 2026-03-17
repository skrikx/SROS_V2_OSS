package mirror

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"srosv2/internal/shared/ids"
)

type WitnessEvent struct {
	WitnessID  string        `json:"witness_id"`
	RunID      ids.RunID     `json:"run_id"`
	SessionID  ids.SessionID `json:"session_id,omitempty"`
	Basis      string        `json:"basis"`
	Severity   string        `json:"severity"`
	Message    string        `json:"message"`
	Signals    []string      `json:"signals,omitempty"`
	OccurredAt time.Time     `json:"occurred_at"`
}

func (e *Engine) writeWitness(event WitnessEvent) error {
	data, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal witness event: %w", err)
	}
	path := filepath.Join(e.root, "witness", event.WitnessID+".json")
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("write witness event: %w", err)
	}
	return nil
}
