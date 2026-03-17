package mirror

type DriftFlag struct {
	Level   string   `json:"level"`
	Reasons []string `json:"reasons"`
}

type RuntimeSnapshot struct {
	RunID           string   `json:"run_id"`
	SessionID       string   `json:"session_id,omitempty"`
	RuntimeState    string   `json:"runtime_state"`
	PlanPath        string   `json:"plan_path,omitempty"`
	LastDecision    string   `json:"last_decision,omitempty"`
	MemoryMutations int      `json:"memory_mutations"`
	BranchCount     int      `json:"branch_count"`
	PendingApproval bool     `json:"pending_approval"`
	Signals         []string `json:"signals,omitempty"`
}

func DetectDrift(snapshot RuntimeSnapshot) DriftFlag {
	reasons := []string{}
	level := "low"
	if snapshot.RuntimeState == "failed_safe" {
		level = "high"
		reasons = append(reasons, "runtime entered failed_safe")
	}
	if snapshot.PendingApproval {
		if level != "high" {
			level = "medium"
		}
		reasons = append(reasons, "operator checkpoint pending")
	}
	if snapshot.MemoryMutations == 0 {
		if level == "low" {
			level = "medium"
		}
		reasons = append(reasons, "no explicit memory lineage recorded yet")
	}
	if snapshot.BranchCount > 3 {
		if level != "high" {
			level = "medium"
		}
		reasons = append(reasons, "branch fanout is above local default")
	}
	return DriftFlag{Level: level, Reasons: reasons}
}
