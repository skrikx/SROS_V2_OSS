package config

import "strings"

func ApplyEnvOverrides(cfg Config, lookup func(string) string) Config {
	if lookup == nil {
		return cfg
	}

	if v := strings.TrimSpace(lookup("SROS_MODE")); v != "" {
		cfg.Mode = Mode(v)
	}
	if v := strings.TrimSpace(lookup("SROS_WORKSPACE_ROOT")); v != "" {
		cfg.WorkspaceRoot = v
	}
	if v := strings.TrimSpace(lookup("SROS_ARTIFACT_ROOT")); v != "" {
		cfg.ArtifactRoot = v
	}
	if v := strings.TrimSpace(lookup("SROS_POLICY_BUNDLE_PATH")); v != "" {
		cfg.PolicyBundlePath = v
	}
	if v := strings.TrimSpace(lookup("SROS_MEMORY_STORE_PATH")); v != "" {
		cfg.MemoryStorePath = v
	}
	if v := strings.TrimSpace(lookup("SROS_TRACE_STORE_PATH")); v != "" {
		cfg.TraceStorePath = v
	}
	if v := strings.TrimSpace(lookup("SROS_OUTPUT_FORMAT")); v != "" {
		cfg.OutputFormat = strings.ToLower(v)
	}
	if v := strings.TrimSpace(lookup("SROS_DATABASE_ENABLED")); v != "" {
		cfg.Database.Enabled = strings.EqualFold(v, "true")
	}
	if v := strings.TrimSpace(lookup("SROS_DATABASE_DRIVER")); v != "" {
		cfg.Database.Driver = v
	}
	if v := strings.TrimSpace(lookup("SROS_DATABASE_URL")); v != "" {
		cfg.Database.URL = v
	}
	if v := strings.TrimSpace(lookup("SROS_DATABASE_MIGRATIONS_DIR")); v != "" {
		cfg.Database.MigrationsDir = v
	}
	if v := strings.TrimSpace(lookup("SROS_DATABASE_SCHEMA")); v != "" {
		cfg.Database.Schema = v
	}

	return cfg
}
