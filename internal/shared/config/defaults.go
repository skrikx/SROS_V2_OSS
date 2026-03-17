package config

import "path/filepath"

func DefaultConfig(cwd string) Config {
	if cwd == "" {
		cwd = "."
	}

	root := filepath.Clean(cwd)
	return Config{
		Mode:             ModeLocalCLI,
		WorkspaceRoot:    root,
		ArtifactRoot:     filepath.Join(root, "artifacts"),
		PolicyBundlePath: filepath.Join(root, "contracts", "policy", "local.bundle.json"),
		MemoryStorePath:  filepath.Join(root, "artifacts", "memory"),
		TraceStorePath:   filepath.Join(root, "artifacts", "trace"),
		OutputFormat:     "text",
	}
}
