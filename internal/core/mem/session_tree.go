package mem

import (
	cmemory "srosv2/contracts/memory"
)

type SessionNode struct {
	SessionID string        `json:"session_id"`
	RunID     string        `json:"run_id"`
	Keys      []string      `json:"keys"`
	Children  []SessionNode `json:"children,omitempty"`
}

func BuildSessionTree(entries []cmemory.MemoryMutation) []SessionNode {
	bySession := map[string]*SessionNode{}
	order := []string{}
	for _, entry := range entries {
		sid := string(entry.SessionID)
		node := bySession[sid]
		if node == nil {
			node = &SessionNode{SessionID: sid, RunID: string(entry.RunID), Keys: []string{}}
			bySession[sid] = node
			order = append(order, sid)
		}
		node.Keys = append(node.Keys, entry.Key)
	}
	out := make([]SessionNode, 0, len(order))
	for _, sid := range order {
		out = append(out, *bySession[sid])
	}
	return out
}
