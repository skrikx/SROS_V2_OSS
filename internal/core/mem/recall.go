package mem

import (
	"sort"
	"strings"
)

type RecallIndex struct {
	terms map[string]map[string]bool
}

func newRecallIndex() *RecallIndex {
	return &RecallIndex{terms: map[string]map[string]bool{}}
}

func (r *RecallIndex) Add(key, value string) string {
	ref := "recall:" + key
	for _, token := range tokenize(key + " " + value) {
		bucket := r.terms[token]
		if bucket == nil {
			bucket = map[string]bool{}
			r.terms[token] = bucket
		}
		bucket[key] = true
	}
	return ref
}

func (r *RecallIndex) Query(query string) []string {
	keys := map[string]bool{}
	for _, token := range tokenize(query) {
		for key := range r.terms[token] {
			keys[key] = true
		}
	}
	out := make([]string, 0, len(keys))
	for key := range keys {
		out = append(out, key)
	}
	sort.Strings(out)
	return out
}

func tokenize(v string) []string {
	fields := strings.Fields(strings.ToLower(v))
	out := make([]string, 0, len(fields))
	for _, field := range fields {
		clean := strings.Trim(field, ".,:;!?'\"")
		if clean != "" {
			out = append(out, clean)
		}
	}
	return out
}
