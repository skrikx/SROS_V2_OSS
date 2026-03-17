package trace

import (
	"srosv2/internal/shared/ids"
)

type ReplayResult struct {
	RunID        ids.RunID        `json:"run_id"`
	EventCount   int              `json:"event_count"`
	States       []string         `json:"states"`
	ReceiptRefs  []string         `json:"receipt_refs,omitempty"`
	ArtifactRefs []string         `json:"artifact_refs,omitempty"`
	Mode         string           `json:"mode"`
}

type Replayer struct {
	reader *Reader
}

func NewReplayer(reader *Reader) *Replayer {
	return &Replayer{reader: reader}
}

func (r *Replayer) Replay(runID ids.RunID) (ReplayResult, error) {
	events, err := r.reader.Events(runID)
	if err != nil {
		return ReplayResult{}, err
	}
	result := ReplayResult{RunID: runID, EventCount: len(events), States: []string{}, ReceiptRefs: []string{}, ArtifactRefs: []string{}, Mode: "reconstructed_lineage"}
	for _, event := range events {
		if state, ok := event.Payload["to"].(string); ok {
			result.States = append(result.States, state)
		}
		if event.ReceiptRef != "" {
			result.ReceiptRefs = append(result.ReceiptRefs, string(event.ReceiptRef))
		}
		for _, ref := range event.ArtifactRefs {
			result.ArtifactRefs = append(result.ArtifactRefs, string(ref))
		}
	}
	return result, nil
}
