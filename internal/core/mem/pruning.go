package mem

import (
	"sort"
	"time"
)

type PrunePlan struct {
	Candidates []string  `json:"candidates"`
	Cutoff     time.Time `json:"cutoff"`
}

func BuildPrunePlan(records []MemoryRecord, cutoff time.Time) PrunePlan {
	out := PrunePlan{Candidates: []string{}, Cutoff: cutoff.UTC()}
	for _, record := range records {
		if record.UpdatedAt.Before(cutoff.UTC()) {
			out.Candidates = append(out.Candidates, record.Key)
		}
	}
	sort.Strings(out.Candidates)
	return out
}
