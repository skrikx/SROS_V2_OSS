package provenance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

func (s *Service) ExportBundle(inputPath string) (map[string]any, error) {
	data, err := os.ReadFile(inputPath)
	if err != nil {
		return nil, fmt.Errorf("read receipt bundle input: %w", err)
	}
	var payload map[string]any
	if err := json.Unmarshal(data, &payload); err != nil {
		return nil, fmt.Errorf("decode receipt bundle input: %w", err)
	}
	out := filepath.Join(s.root, "exports", filepath.Base(inputPath))
	if err := os.WriteFile(out, append(data, '\n'), 0o644); err != nil {
		return nil, fmt.Errorf("write exported bundle: %w", err)
	}
	return map[string]any{"exported_to": out, "bundle": payload}, nil
}
