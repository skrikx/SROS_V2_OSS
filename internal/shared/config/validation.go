package config

import (
	"fmt"
	"strings"
)

func Validate(cfg Config) error {
	if cfg.Mode != ModeLocalCLI {
		return fmt.Errorf("mode must be %q", ModeLocalCLI)
	}

	if strings.TrimSpace(cfg.WorkspaceRoot) == "" {
		return fmt.Errorf("workspace_root is required")
	}
	if strings.TrimSpace(cfg.ArtifactRoot) == "" {
		return fmt.Errorf("artifact_root is required")
	}
	if strings.TrimSpace(cfg.PolicyBundlePath) == "" {
		return fmt.Errorf("policy_bundle_path is required")
	}
	if strings.TrimSpace(cfg.MemoryStorePath) == "" {
		return fmt.Errorf("memory_store_path is required")
	}
	if strings.TrimSpace(cfg.TraceStorePath) == "" {
		return fmt.Errorf("trace_store_path is required")
	}

	if cfg.OutputFormat != "text" && cfg.OutputFormat != "json" {
		return fmt.Errorf("output_format must be one of [text, json]")
	}

	return nil
}
