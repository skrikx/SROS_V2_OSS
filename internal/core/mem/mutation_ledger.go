package mem

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	cmemory "srosv2/contracts/memory"
	"srosv2/internal/shared/ids"
)

type MutationInput struct {
	Scope            ScopeBinding
	OperatorID       ids.OperatorID
	Kind             cmemory.MutationKind
	Branch           cmemory.BranchReference
	ParentMutationID ids.MemoryMutationID
	Key              string
	Value            string
	Reason           string
	RecallIndexRef   string
	OccurredAt       time.Time
}

type MutationLedger struct {
	root string
	hook func(cmemory.MemoryMutation)
}

func newLedger(root string) (*MutationLedger, error) {
	path := filepath.Join(root, "mutations")
	if err := os.MkdirAll(path, 0o755); err != nil {
		return nil, fmt.Errorf("create mutation ledger: %w", err)
	}
	return &MutationLedger{root: path}, nil
}

func (l *MutationLedger) Append(input MutationInput) (cmemory.MemoryMutation, error) {
	if err := input.Scope.Validate(); err != nil {
		return cmemory.MemoryMutation{}, err
	}
	if input.OperatorID == "" {
		return cmemory.MemoryMutation{}, fmt.Errorf("operator id is required")
	}
	at := input.OccurredAt.UTC()
	if at.IsZero() {
		at = time.Now().UTC()
	}
	id := ids.MemoryMutationID("mm_" + shortHash(fmt.Sprintf("%s|%s|%s|%d", input.Scope.WorkspaceID, input.Branch.BranchID, input.Key, at.UnixNano())))
	lineageRef := fmt.Sprintf("lineage:%s", id)
	mutation := cmemory.MemoryMutation{
		ContractVersion:  "v2.0",
		MutationID:       id,
		RunID:            input.Scope.RunID,
		SessionID:        input.Scope.SessionID,
		TenantID:         input.Scope.TenantID,
		WorkspaceID:      input.Scope.WorkspaceID,
		OperatorID:       input.OperatorID,
		Scope:            input.Scope.Scope,
		Kind:             input.Kind,
		Branch:           input.Branch,
		ParentMutationID: input.ParentMutationID,
		LineageRef:       lineageRef,
		RecallIndexRef:   input.RecallIndexRef,
		Key:              input.Key,
		Value:            input.Value,
		Reason:           input.Reason,
		OccurredAt:       at,
	}
	if errs := cmemory.ValidateMutation(mutation); len(errs) > 0 {
		return cmemory.MemoryMutation{}, errs[0]
	}
	path := filepath.Join(l.root, string(mutation.MutationID)+".json")
	data, err := json.MarshalIndent(mutation, "", "  ")
	if err != nil {
		return cmemory.MemoryMutation{}, fmt.Errorf("marshal mutation: %w", err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return cmemory.MemoryMutation{}, fmt.Errorf("write mutation: %w", err)
	}
	if l.hook != nil {
		l.hook(mutation)
	}
	return mutation, nil
}

func (l *MutationLedger) Entries() ([]cmemory.MemoryMutation, error) {
	files, err := os.ReadDir(l.root)
	if err != nil {
		return nil, fmt.Errorf("read mutation ledger: %w", err)
	}
	out := make([]cmemory.MemoryMutation, 0, len(files))
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}
		var entry cmemory.MemoryMutation
		data, err := os.ReadFile(filepath.Join(l.root, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("read mutation file: %w", err)
		}
		if err := json.Unmarshal(data, &entry); err != nil {
			return nil, fmt.Errorf("decode mutation file: %w", err)
		}
		out = append(out, entry)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].OccurredAt.Before(out[j].OccurredAt) })
	return out, nil
}

func (l *MutationLedger) SetHook(hook func(cmemory.MemoryMutation)) {
	l.hook = hook
}
