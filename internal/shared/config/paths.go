package config

import (
	"path/filepath"
)

func DefaultConfigPaths(workspaceRoot string) []string {
	root := filepath.Clean(workspaceRoot)
	return []string{
		filepath.Join(root, "sros.yaml"),
		filepath.Join(root, ".sros", "config.yaml"),
	}
}

func ResolvePath(baseDir, raw string) string {
	if raw == "" {
		return raw
	}
	if filepath.IsAbs(raw) {
		return filepath.Clean(raw)
	}
	return filepath.Clean(filepath.Join(baseDir, raw))
}
