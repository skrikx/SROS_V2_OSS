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

type BranchManager struct {
	root string
	now  func() time.Time
}

func newBranchManager(root string, now func() time.Time) (*BranchManager, error) {
	path := filepath.Join(root, "branches")
	if err := os.MkdirAll(path, 0o755); err != nil {
		return nil, fmt.Errorf("create branch manager root: %w", err)
	}
	return &BranchManager{root: path, now: now}, nil
}

func (b *BranchManager) Save(record cmemory.BranchRecord) error {
	if record.CreatedAt.IsZero() {
		record.CreatedAt = b.now().UTC()
	}
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal branch record: %w", err)
	}
	path := filepath.Join(b.root, string(record.BranchID)+".json")
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("write branch record: %w", err)
	}
	return nil
}

func (b *BranchManager) List() ([]cmemory.BranchRecord, error) {
	files, err := os.ReadDir(b.root)
	if err != nil {
		return nil, fmt.Errorf("read branches: %w", err)
	}
	out := make([]cmemory.BranchRecord, 0, len(files))
	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".json" {
			continue
		}
		data, err := os.ReadFile(filepath.Join(b.root, file.Name()))
		if err != nil {
			return nil, fmt.Errorf("read branch record: %w", err)
		}
		var record cmemory.BranchRecord
		if err := json.Unmarshal(data, &record); err != nil {
			return nil, fmt.Errorf("decode branch record: %w", err)
		}
		out = append(out, record)
	}
	sort.Slice(out, func(i, j int) bool { return string(out[i].BranchID) < string(out[j].BranchID) })
	return out, nil
}

func (b *BranchManager) Load(branchID ids.BranchID) (cmemory.BranchRecord, error) {
	path := filepath.Join(b.root, string(branchID)+".json")
	data, err := os.ReadFile(path)
	if err != nil {
		return cmemory.BranchRecord{}, fmt.Errorf("read branch record %s: %w", branchID, err)
	}
	var record cmemory.BranchRecord
	if err := json.Unmarshal(data, &record); err != nil {
		return cmemory.BranchRecord{}, fmt.Errorf("decode branch record %s: %w", branchID, err)
	}
	return record, nil
}
