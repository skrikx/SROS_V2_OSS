package mirror

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"srosv2/internal/shared/ids"
)

type Engine struct {
	root        string
	now         func() time.Time
	witnessHook func(WitnessEvent)
}

func New(root string, now func() time.Time) (*Engine, error) {
	if root == "" {
		root = filepath.Join("artifacts", "mirror")
	}
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	for _, rel := range []string{"witness", "summaries"} {
		if err := os.MkdirAll(filepath.Join(root, rel), 0o755); err != nil {
			return nil, fmt.Errorf("create mirror root: %w", err)
		}
	}
	return &Engine{root: root, now: now}, nil
}

func (e *Engine) Observe(snapshot RuntimeSnapshot, basis string) (WitnessEvent, ReflectionSummary, error) {
	drift := DetectDrift(snapshot)
	event := WitnessEvent{
		WitnessID:  "wit_" + shortHash(snapshot.RunID+"|"+snapshot.RuntimeState+"|"+basis),
		RunID:      ids.RunID(snapshot.RunID),
		SessionID:  ids.SessionID(snapshot.SessionID),
		Basis:      basis,
		Severity:   drift.Level,
		Message:    "mirror witness derived from runtime state and local lineage",
		Signals:    append([]string{}, snapshot.Signals...),
		OccurredAt: e.now().UTC(),
	}
	if err := e.writeWitness(event); err != nil {
		return WitnessEvent{}, ReflectionSummary{}, err
	}
	summary := BuildSummary(snapshot, drift, 1)
	if err := e.writeSummary(summary); err != nil {
		return WitnessEvent{}, ReflectionSummary{}, err
	}
	return event, summary, nil
}

func (e *Engine) writeSummary(summary ReflectionSummary) error {
	data, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal reflection summary: %w", err)
	}
	path := filepath.Join(e.root, "summaries", shortHash(summary.RunID+"|"+summary.SessionID)+".json")
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("write reflection summary: %w", err)
	}
	return nil
}

func (e *Engine) StatusFromFile(path string) (map[string]any, error) {
	snapshot, err := loadSnapshot(path)
	if err != nil {
		return nil, err
	}
	event, summary, err := e.Observe(snapshot, "file_input")
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"snapshot": snapshot,
		"witness":  event,
		"summary":  summary,
		"drift":    DetectDrift(snapshot),
	}, nil
}

func (e *Engine) WitnessFromFile(path string) (map[string]any, error) {
	return e.StatusFromFile(path)
}

func (e *Engine) SetWitnessHook(hook func(WitnessEvent)) {
	e.witnessHook = hook
}

func loadSnapshot(path string) (RuntimeSnapshot, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return RuntimeSnapshot{}, fmt.Errorf("read mirror input: %w", err)
	}
	var snapshot RuntimeSnapshot
	if err := json.Unmarshal(data, &snapshot); err != nil {
		return RuntimeSnapshot{}, fmt.Errorf("decode mirror input: %w", err)
	}
	return snapshot, nil
}
