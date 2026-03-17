package mem

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	cmemory "srosv2/contracts/memory"
	"srosv2/internal/shared/ids"
)

type MemoryRecord struct {
	Key            string               `json:"key"`
	Value          string               `json:"value"`
	Scope          ScopeBinding         `json:"scope"`
	BranchID       ids.BranchID         `json:"branch_id"`
	LastMutationID ids.MemoryMutationID `json:"last_mutation_id"`
	UpdatedAt      time.Time            `json:"updated_at"`
}

type SeedFile struct {
	OperatorID  ids.OperatorID  `json:"operator_id"`
	TenantID    ids.TenantID    `json:"tenant_id"`
	WorkspaceID ids.WorkspaceID `json:"workspace_id"`
	RunID       ids.RunID       `json:"run_id"`
	SessionID   ids.SessionID   `json:"session_id"`
	BranchID    ids.BranchID    `json:"branch_id"`
	Entries     []struct {
		Scope cmemory.Scope `json:"scope"`
		Key   string        `json:"key"`
		Value string        `json:"value"`
	} `json:"entries"`
}

type BranchFile struct {
	OperatorID       ids.OperatorID       `json:"operator_id"`
	TenantID         ids.TenantID         `json:"tenant_id"`
	WorkspaceID      ids.WorkspaceID      `json:"workspace_id"`
	RunID            ids.RunID            `json:"run_id"`
	SessionID        ids.SessionID        `json:"session_id"`
	BranchID         ids.BranchID         `json:"branch_id"`
	ParentBranchID   ids.BranchID         `json:"parent_branch_id"`
	BaseMutationID   ids.MemoryMutationID `json:"base_mutation_id,omitempty"`
	RewindToMutation ids.MemoryMutationID `json:"rewind_to_mutation,omitempty"`
	Reason           string               `json:"reason,omitempty"`
	Entries          []struct {
		Scope cmemory.Scope `json:"scope"`
		Key   string        `json:"key"`
		Value string        `json:"value"`
	} `json:"entries"`
}

type Store struct {
	root     string
	now      func() time.Time
	ledger   *MutationLedger
	branches *BranchManager
	recall   *RecallIndex
}

func NewStore(root string, now func() time.Time) (*Store, error) {
	if strings.TrimSpace(root) == "" {
		root = filepath.Join("artifacts", "memory")
	}
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	for _, rel := range []string{"records", "snapshots"} {
		if err := os.MkdirAll(filepath.Join(root, rel), 0o755); err != nil {
			return nil, fmt.Errorf("create memory root: %w", err)
		}
	}
	ledger, err := newLedger(root)
	if err != nil {
		return nil, err
	}
	branches, err := newBranchManager(root, now)
	if err != nil {
		return nil, err
	}
	store := &Store{root: root, now: now, ledger: ledger, branches: branches, recall: newRecallIndex()}
	if err := store.reloadRecall(); err != nil {
		return nil, err
	}
	return store, nil
}

func (s *Store) Root() string { return s.root }

func (s *Store) Upsert(input MutationInput) (cmemory.MemoryMutation, error) {
	ref := s.recall.Add(input.Key, input.Value)
	input.RecallIndexRef = ref
	mutation, err := s.ledger.Append(input)
	if err != nil {
		return cmemory.MemoryMutation{}, err
	}
	record := MemoryRecord{
		Key:            input.Key,
		Value:          input.Value,
		Scope:          input.Scope,
		BranchID:       input.Branch.BranchID,
		LastMutationID: mutation.MutationID,
		UpdatedAt:      mutation.OccurredAt,
	}
	if err := s.writeRecord(record); err != nil {
		return cmemory.MemoryMutation{}, err
	}
	branch, err := s.loadOrDefaultBranch(input.Scope, input.OperatorID, input.Branch)
	if err != nil {
		return cmemory.MemoryMutation{}, err
	}
	branch.HeadMutationID = mutation.MutationID
	if err := s.branches.Save(branch); err != nil {
		return cmemory.MemoryMutation{}, err
	}
	return mutation, nil
}

func (s *Store) loadOrDefaultBranch(scope ScopeBinding, operatorID ids.OperatorID, ref cmemory.BranchReference) (cmemory.BranchRecord, error) {
	branchID := ref.BranchID
	if branchID == "" {
		branchID = ids.BranchID("branch_main")
	}
	record, err := s.branches.Load(branchID)
	if err == nil {
		return record, nil
	}
	return cmemory.BranchRecord{
		BranchID:       branchID,
		ParentBranchID: ref.ParentBranchID,
		HeadMutationID: ref.HeadMutationID,
		CreatedBy:      operatorID,
		TenantID:       scope.TenantID,
		WorkspaceID:    scope.WorkspaceID,
		CreatedAt:      s.now().UTC(),
		Reason:         "branch initialized",
	}, nil
}

func (s *Store) Records() ([]MemoryRecord, error) {
	files, err := os.ReadDir(filepath.Join(s.root, "records"))
	if err != nil {
		return nil, fmt.Errorf("read memory records: %w", err)
	}
	out := make([]MemoryRecord, 0, len(files))
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}
		var record MemoryRecord
		data, err := os.ReadFile(filepath.Join(s.root, "records", file.Name()))
		if err != nil {
			return nil, fmt.Errorf("read memory record: %w", err)
		}
		if err := json.Unmarshal(data, &record); err != nil {
			return nil, fmt.Errorf("decode memory record: %w", err)
		}
		out = append(out, record)
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Key < out[j].Key })
	return out, nil
}

func (s *Store) Recall(query string) (map[string]any, error) {
	records, err := s.Records()
	if err != nil {
		return nil, err
	}
	keys := s.recall.Query(query)
	matched := []MemoryRecord{}
	for _, record := range records {
		for _, key := range keys {
			if record.Key == key {
				matched = append(matched, record)
			}
		}
	}
	mutations, err := s.ledger.Entries()
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"query":          query,
		"matches":        matched,
		"match_count":    len(matched),
		"mutation_count": len(mutations),
	}, nil
}

func (s *Store) ImportSeed(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read seed file: %w", err)
	}
	var seed SeedFile
	if err := json.Unmarshal(data, &seed); err != nil {
		return nil, fmt.Errorf("decode seed file: %w", err)
	}
	var last ids.MemoryMutationID
	for _, entry := range seed.Entries {
		mutation, err := s.Upsert(MutationInput{
			Scope: ScopeBinding{
				Scope:       entry.Scope,
				TenantID:    seed.TenantID,
				WorkspaceID: seed.WorkspaceID,
				RunID:       seed.RunID,
				SessionID:   seed.SessionID,
			},
			OperatorID:       seed.OperatorID,
			Kind:             cmemory.MutationKindUpsert,
			Branch:           cmemory.BranchReference{BranchID: seed.BranchID, HeadMutationID: last},
			ParentMutationID: last,
			Key:              entry.Key,
			Value:            entry.Value,
			Reason:           "workspace seed import",
			OccurredAt:       s.now().UTC(),
		})
		if err != nil {
			return nil, err
		}
		last = mutation.MutationID
	}
	return s.Recall("workspace")
}

func (s *Store) ApplyBranchFile(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read branch file: %w", err)
	}
	var file BranchFile
	if err := json.Unmarshal(data, &file); err != nil {
		return nil, fmt.Errorf("decode branch file: %w", err)
	}
	record := cmemory.BranchRecord{
		BranchID:       file.BranchID,
		ParentBranchID: file.ParentBranchID,
		HeadMutationID: file.BaseMutationID,
		CreatedBy:      file.OperatorID,
		TenantID:       file.TenantID,
		WorkspaceID:    file.WorkspaceID,
		CreatedAt:      s.now().UTC(),
		Reason:         file.Reason,
	}
	if err := s.branches.Save(record); err != nil {
		return nil, err
	}
	head := file.BaseMutationID
	for _, entry := range file.Entries {
		mutation, err := s.Upsert(MutationInput{
			Scope: ScopeBinding{
				Scope:       entry.Scope,
				TenantID:    file.TenantID,
				WorkspaceID: file.WorkspaceID,
				RunID:       file.RunID,
				SessionID:   file.SessionID,
			},
			OperatorID:       file.OperatorID,
			Kind:             cmemory.MutationKindUpsert,
			Branch:           cmemory.BranchReference{BranchID: file.BranchID, ParentBranchID: file.ParentBranchID, HeadMutationID: head},
			ParentMutationID: head,
			Key:              entry.Key,
			Value:            entry.Value,
			Reason:           "branch apply",
			OccurredAt:       s.now().UTC(),
		})
		if err != nil {
			return nil, err
		}
		head = mutation.MutationID
	}
	if file.RewindToMutation != "" {
		if err := s.Rewind(file.BranchID, file.RewindToMutation, file.OperatorID, file.Reason); err != nil {
			return nil, err
		}
	}
	branches, err := s.branches.List()
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"branch":   file.BranchID,
		"head":     head,
		"branches": branches,
	}, nil
}

func (s *Store) Rewind(branchID ids.BranchID, target ids.MemoryMutationID, operatorID ids.OperatorID, reason string) error {
	entries, err := s.ledger.Entries()
	if err != nil {
		return err
	}
	seen := false
	var targetEntry cmemory.MemoryMutation
	for _, entry := range entries {
		if entry.Branch.BranchID == branchID && entry.MutationID == target {
			seen = true
			targetEntry = entry
			break
		}
	}
	if !seen {
		return fmt.Errorf("rewind target %s not found on branch %s", target, branchID)
	}
	branch, err := s.branches.Load(branchID)
	if err != nil {
		return err
	}
	branch.RewoundTo = target
	branch.HeadMutationID = target
	branch.Reason = reason
	if err := s.branches.Save(branch); err != nil {
		return err
	}
	_, err = s.ledger.Append(MutationInput{
		Scope: ScopeBinding{
			Scope:       targetEntry.Scope,
			TenantID:    targetEntry.TenantID,
			WorkspaceID: targetEntry.WorkspaceID,
			RunID:       targetEntry.RunID,
			SessionID:   targetEntry.SessionID,
		},
		OperatorID:       operatorID,
		Kind:             cmemory.MutationKindAnnotate,
		Branch:           cmemory.BranchReference{BranchID: branchID, ParentBranchID: branch.ParentBranchID, HeadMutationID: target},
		ParentMutationID: target,
		Key:              "branch.rewind",
		Value:            string(target),
		Reason:           "rewind applied: " + reason,
		OccurredAt:       s.now().UTC(),
	})
	return err
}

func (s *Store) Branches() ([]cmemory.BranchRecord, error) {
	return s.branches.List()
}

func (s *Store) Ledger() ([]cmemory.MemoryMutation, error) {
	return s.ledger.Entries()
}

func (s *Store) SessionTree() ([]SessionNode, error) {
	entries, err := s.ledger.Entries()
	if err != nil {
		return nil, err
	}
	return BuildSessionTree(entries), nil
}

func (s *Store) PrunePlan(cutoff time.Time) (PrunePlan, error) {
	records, err := s.Records()
	if err != nil {
		return PrunePlan{}, err
	}
	return BuildPrunePlan(records, cutoff), nil
}

func (s *Store) CompactionPlan() (CompactionPlan, error) {
	records, err := s.Records()
	if err != nil {
		return CompactionPlan{}, err
	}
	mutations, err := s.Ledger()
	if err != nil {
		return CompactionPlan{}, err
	}
	return BuildCompactionPlan(records, len(mutations)), nil
}

func (s *Store) writeRecord(record MemoryRecord) error {
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal memory record: %w", err)
	}
	path := filepath.Join(s.root, "records", shortHash(record.Key)+".json")
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("write memory record: %w", err)
	}
	return nil
}

func (s *Store) reloadRecall() error {
	records, err := s.Records()
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return nil
	}
	s.recall = newRecallIndex()
	for _, record := range records {
		s.recall.Add(record.Key, record.Value)
	}
	return nil
}
