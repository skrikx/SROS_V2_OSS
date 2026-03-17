package mem

import "sort"

type CompactionPlan struct {
	Keys          []string `json:"keys"`
	MutationCount int      `json:"mutation_count"`
	Summary       string   `json:"summary"`
}

func BuildCompactionPlan(records []MemoryRecord, mutationCount int) CompactionPlan {
	keys := make([]string, 0, len(records))
	for _, record := range records {
		keys = append(keys, record.Key)
	}
	sort.Strings(keys)
	return CompactionPlan{
		Keys:          keys,
		MutationCount: mutationCount,
		Summary:       "compaction preserves lineage while reducing recall surface",
	}
}
