package config

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func Load(opts LoadOptions) (LoadResult, error) {
	cwd := opts.CWD
	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return LoadResult{}, fmt.Errorf("resolve working directory: %w", err)
		}
	}

	cfg := DefaultConfig(cwd)
	source := "defaults"
	warnings := make([]string, 0)

	lookup := opts.LookupEnv
	if lookup == nil {
		lookup = os.Getenv
	}

	configPath := strings.TrimSpace(opts.ExplicitPath)
	if configPath == "" {
		configPath = strings.TrimSpace(lookup("SROS_CONFIG"))
	}

	if configPath != "" {
		loaded, err := loadFromFile(configPath)
		if err != nil {
			return LoadResult{}, err
		}
		cfg = merge(cfg, loaded)
		source = filepath.Clean(configPath)
	} else {
		for _, candidate := range DefaultConfigPaths(cwd) {
			if _, err := os.Stat(candidate); err == nil {
				loaded, loadErr := loadFromFile(candidate)
				if loadErr != nil {
					return LoadResult{}, loadErr
				}
				cfg = merge(cfg, loaded)
				source = filepath.Clean(candidate)
				break
			}
		}
	}

	cfg = ApplyEnvOverrides(cfg, lookup)

	if err := Validate(cfg); err != nil {
		return LoadResult{}, err
	}

	return LoadResult{
		Config:   cfg,
		Source:   source,
		Warnings: warnings,
	}, nil
}

func merge(base, overlay Config) Config {
	if overlay.Mode != "" {
		base.Mode = overlay.Mode
	}
	if overlay.WorkspaceRoot != "" {
		base.WorkspaceRoot = overlay.WorkspaceRoot
	}
	if overlay.ArtifactRoot != "" {
		base.ArtifactRoot = overlay.ArtifactRoot
	}
	if overlay.PolicyBundlePath != "" {
		base.PolicyBundlePath = overlay.PolicyBundlePath
	}
	if overlay.MemoryStorePath != "" {
		base.MemoryStorePath = overlay.MemoryStorePath
	}
	if overlay.TraceStorePath != "" {
		base.TraceStorePath = overlay.TraceStorePath
	}
	if overlay.OutputFormat != "" {
		base.OutputFormat = overlay.OutputFormat
	}
	return base
}

func loadFromFile(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return Config{}, fmt.Errorf("read config file %q: %w", path, err)
	}

	values, err := parseConfigDocument(data)
	if err != nil {
		return Config{}, fmt.Errorf("parse config file %q: %w", path, err)
	}

	base := filepath.Dir(path)
	cfg := Config{}
	cfg.Mode = Mode(values["mode"])
	cfg.WorkspaceRoot = resolveIfPresent(base, values["workspace_root"])
	cfg.ArtifactRoot = resolveIfPresent(base, values["artifact_root"])
	cfg.PolicyBundlePath = resolveIfPresent(base, values["policy_bundle_path"])
	cfg.MemoryStorePath = resolveIfPresent(base, values["memory_store_path"])
	cfg.TraceStorePath = resolveIfPresent(base, values["trace_store_path"])
	cfg.OutputFormat = strings.ToLower(strings.TrimSpace(values["output_format"]))

	return cfg, nil
}

func resolveIfPresent(base, value string) string {
	trimmed := strings.TrimSpace(value)
	if trimmed == "" {
		return ""
	}
	return ResolvePath(base, trimmed)
}

func parseConfigDocument(data []byte) (map[string]string, error) {
	trimmed := strings.TrimSpace(string(data))
	if trimmed == "" {
		return map[string]string{}, nil
	}

	if strings.HasPrefix(trimmed, "{") {
		var decoded map[string]any
		if err := json.Unmarshal(data, &decoded); err != nil {
			return nil, err
		}
		out := make(map[string]string, len(decoded))
		for k, v := range decoded {
			out[k] = fmt.Sprintf("%v", v)
		}
		return out, nil
	}

	out := map[string]string{}
	s := bufio.NewScanner(strings.NewReader(trimmed))
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid line %q", line)
		}
		k := strings.TrimSpace(parts[0])
		v := strings.TrimSpace(parts[1])
		v = strings.Trim(v, "\"")
		out[k] = v
	}

	if err := s.Err(); err != nil {
		return nil, err
	}

	return out, nil
}
