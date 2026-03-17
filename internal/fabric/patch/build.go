package patch

func BuildDiff(path string, before, after []byte) map[string]any {
	return map[string]any{
		"path":    path,
		"before":  string(before),
		"after":   string(after),
		"changed": string(before) != string(after),
	}
}
