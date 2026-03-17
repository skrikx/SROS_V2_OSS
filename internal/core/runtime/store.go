package runtime

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

type ApprovalCheckpoint struct {
	SessionID   string `json:"session_id"`
	Reason      string `json:"reason"`
	Approved    bool   `json:"approved"`
	ApprovedBy  string `json:"approved_by,omitempty"`
	ApprovedAt  string `json:"approved_at,omitempty"`
	RequestedAt string `json:"requested_at"`
}

type Store struct {
	root    string
	mu      sync.Mutex
	pgStore *PostgresStore
}

func NewStore(root string) (*Store, error) {
	if strings.TrimSpace(root) == "" {
		root = filepath.Join("artifacts", "runtime")
	}
	root = filepath.Clean(root)

	dirs := []string{
		filepath.Join(root, "sessions"),
		filepath.Join(root, "checkpoints"),
		filepath.Join(root, "rollbacks"),
		filepath.Join(root, "approvals"),
	}
	for _, d := range dirs {
		if err := os.MkdirAll(d, 0o755); err != nil {
			return nil, fmt.Errorf("create runtime store directory %s: %w", d, err)
		}
	}

	return &Store{root: root}, nil
}

func (s *Store) Root() string {
	return s.root
}

func (s *Store) SetPostgresStore(store *PostgresStore) {
	s.pgStore = store
}

func (s *Store) SaveSession(session RuntimeSession) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.root, "sessions", session.SessionID+".json")
	if err := writeJSON(path, session); err != nil {
		return err
	}
	latest := filepath.Join(s.root, "latest_session.txt")
	if err := os.WriteFile(latest, []byte(session.SessionID), 0o644); err != nil {
		return fmt.Errorf("write latest session marker: %w", err)
	}
	if s.pgStore != nil {
		if err := s.pgStore.SaveSession(context.Background(), session); err != nil {
			return err
		}
	}
	return nil
}

func (s *Store) LoadSession(sessionID string) (RuntimeSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.root, "sessions", sessionID+".json")
	var session RuntimeSession
	if err := readJSON(path, &session); err != nil {
		return RuntimeSession{}, err
	}
	return session, nil
}

func (s *Store) LatestSession() (RuntimeSession, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	marker := filepath.Join(s.root, "latest_session.txt")
	data, err := os.ReadFile(marker)
	if err != nil {
		return RuntimeSession{}, fmt.Errorf("read latest session marker: %w", err)
	}
	sessionID := strings.TrimSpace(string(data))
	if sessionID == "" {
		return RuntimeSession{}, fmt.Errorf("latest session marker is empty")
	}
	path := filepath.Join(s.root, "sessions", sessionID+".json")
	var session RuntimeSession
	if err := readJSON(path, &session); err != nil {
		return RuntimeSession{}, err
	}
	return session, nil
}

func (s *Store) SaveCheckpoint(cp RuntimeCheckpoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.root, "checkpoints", string(cp.Record.CheckpointID)+".json")
	if err := writeJSON(path, cp); err != nil {
		return err
	}
	if s.pgStore != nil {
		return s.pgStore.SaveCheckpoint(context.Background(), cp)
	}
	return nil
}

func (s *Store) LoadCheckpoint(id string) (RuntimeCheckpoint, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.root, "checkpoints", id+".json")
	var cp RuntimeCheckpoint
	if err := readJSON(path, &cp); err != nil {
		return RuntimeCheckpoint{}, err
	}
	return cp, nil
}

func (s *Store) SaveRollback(rb RuntimeRollback) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.root, "rollbacks", string(rb.Record.RollbackID)+".json")
	if err := writeJSON(path, rb); err != nil {
		return err
	}
	if s.pgStore != nil {
		return s.pgStore.SaveRollback(context.Background(), rb)
	}
	return nil
}

func (s *Store) SaveApproval(a ApprovalCheckpoint) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.root, "approvals", a.SessionID+".json")
	if err := writeJSON(path, a); err != nil {
		return err
	}
	if s.pgStore != nil {
		return s.pgStore.SaveApproval(context.Background(), RuntimeSession{SessionID: a.SessionID, RunID: a.SessionID}, a)
	}
	return nil
}

func (s *Store) LoadApproval(sessionID string) (ApprovalCheckpoint, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	path := filepath.Join(s.root, "approvals", sessionID+".json")
	var a ApprovalCheckpoint
	if err := readJSON(path, &a); err != nil {
		return ApprovalCheckpoint{}, err
	}
	return a, nil
}

func (s *Store) ApprovalPath(sessionID string) string {
	return filepath.Join(s.root, "approvals", sessionID+".json")
}

func writeJSON(path string, value any) error {
	data, err := json.MarshalIndent(value, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal json %s: %w", path, err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("write json %s: %w", path, err)
	}
	return nil
}

func readJSON(path string, out any) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return fmt.Errorf("read json %s: %w", path, err)
	}
	if err := json.Unmarshal(data, out); err != nil {
		return fmt.Errorf("decode json %s: %w", path, err)
	}
	return nil
}
