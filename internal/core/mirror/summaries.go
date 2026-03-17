package mirror

import "fmt"

type ReflectionSummary struct {
	RunID        string   `json:"run_id"`
	SessionID    string   `json:"session_id,omitempty"`
	DriftLevel   string   `json:"drift_level"`
	Summary      string   `json:"summary"`
	SourceBasis  string   `json:"source_basis"`
	WitnessCount int      `json:"witness_count"`
	Signals      []string `json:"signals,omitempty"`
}

func BuildSummary(snapshot RuntimeSnapshot, drift DriftFlag, witnessCount int) ReflectionSummary {
	return ReflectionSummary{
		RunID:        snapshot.RunID,
		SessionID:    snapshot.SessionID,
		DriftLevel:   drift.Level,
		Summary:      fmt.Sprintf("mirror summary derived from runtime state and local lineage with drift=%s", drift.Level),
		SourceBasis:  "runtime_state+local_lineage",
		WitnessCount: witnessCount,
		Signals:      snapshot.Signals,
	}
}
