package provenance

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"srosv2/contracts/evidence"
	"srosv2/internal/shared/ids"
)

type ArtifactProvenance struct {
	RunID      ids.RunID         `json:"run_id"`
	Artifact   evidence.ArtifactRef `json:"artifact"`
	LinkedAt   time.Time         `json:"linked_at"`
	SourceKind string            `json:"source_kind"`
}

type Service struct {
	root string
	now  func() time.Time
}

func New(root string, now func() time.Time) (*Service, error) {
	if root == "" {
		root = filepath.Join("artifacts", "provenance")
	}
	for _, rel := range []string{"receipts", "bundles", "closures", "exports", "artifacts"} {
		if err := os.MkdirAll(filepath.Join(root, rel), 0o755); err != nil {
			return nil, fmt.Errorf("create provenance root: %w", err)
		}
	}
	if now == nil {
		now = func() time.Time { return time.Now().UTC() }
	}
	return &Service{root: root, now: now}, nil
}

func (s *Service) Root() string { return s.root }

func (s *Service) LinkArtifact(runID ids.RunID, path, mediaType, sourceKind string) (evidence.ArtifactRef, error) {
	info, err := os.Stat(path)
	if err != nil {
		return evidence.ArtifactRef{}, fmt.Errorf("stat artifact %s: %w", path, err)
	}
	data, err := os.ReadFile(path)
	if err != nil {
		return evidence.ArtifactRef{}, fmt.Errorf("read artifact %s: %w", path, err)
	}
	ref := evidence.ArtifactRef{
		ArtifactID: ids.ArtifactID("art_" + digestBytes([]byte(path))[:12]),
		Path:       path,
		DigestAlgo: evidence.DigestAlgorithmSHA256,
		Digest:     digestBytes(data),
		MediaType:  mediaType,
		Bytes:      info.Size(),
	}
	prov := ArtifactProvenance{RunID: runID, Artifact: ref, LinkedAt: s.now().UTC(), SourceKind: sourceKind}
	encoded, err := json.MarshalIndent(prov, "", "  ")
	if err != nil {
		return evidence.ArtifactRef{}, fmt.Errorf("marshal artifact provenance: %w", err)
	}
	out := filepath.Join(s.root, "artifacts", string(ref.ArtifactID)+".json")
	if err := os.WriteFile(out, append(encoded, '\n'), 0o644); err != nil {
		return evidence.ArtifactRef{}, fmt.Errorf("write artifact provenance: %w", err)
	}
	return ref, nil
}
