package main

import (
	"fmt"
	"io"
	"os"
	"strings"

	"srosv2/cmd/sros/commands"
	"srosv2/internal/core/boot"
	"srosv2/internal/shared/config"
)

type RunOptions struct {
	Stdout io.Writer
	Stderr io.Writer
	Env    []string
	CWD    string
}

type globalOptions struct {
	ConfigPath string
	Format     string
	Workspace  string
	Help       bool
}

func parseGlobalOptions(args []string) (globalOptions, []string, error) {
	opts := globalOptions{}
	remaining := make([]string, 0, len(args))

	for i := 0; i < len(args); i++ {
		tok := args[i]
		if !strings.HasPrefix(tok, "-") {
			remaining = append(remaining, args[i:]...)
			break
		}

		if tok == "--" {
			remaining = append(remaining, args[i+1:]...)
			break
		}

		if tok == "-h" || tok == "--help" {
			opts.Help = true
			continue
		}

		key, val, hasValue := splitFlag(tok)
		switch key {
		case "--config":
			if !hasValue {
				i++
				if i >= len(args) {
					return globalOptions{}, nil, fmt.Errorf("--config requires a value")
				}
				val = args[i]
			}
			opts.ConfigPath = strings.TrimSpace(val)
		case "--format":
			if !hasValue {
				i++
				if i >= len(args) {
					return globalOptions{}, nil, fmt.Errorf("--format requires a value")
				}
				val = args[i]
			}
			opts.Format = strings.ToLower(strings.TrimSpace(val))
		case "--workspace":
			if !hasValue {
				i++
				if i >= len(args) {
					return globalOptions{}, nil, fmt.Errorf("--workspace requires a value")
				}
				val = args[i]
			}
			opts.Workspace = strings.TrimSpace(val)
		default:
			remaining = append(remaining, args[i:]...)
			i = len(args)
		}
	}

	return opts, remaining, nil
}

func splitFlag(token string) (string, string, bool) {
	parts := strings.SplitN(token, "=", 2)
	if len(parts) == 1 {
		return token, "", false
	}
	return parts[0], parts[1], true
}

func buildCommandContext(global globalOptions, opts RunOptions) (*commands.Context, error) {
	stdout := opts.Stdout
	stderr := opts.Stderr
	if stdout == nil {
		stdout = os.Stdout
	}
	if stderr == nil {
		stderr = os.Stderr
	}

	cwd := opts.CWD
	if cwd == "" {
		var err error
		cwd, err = os.Getwd()
		if err != nil {
			return nil, fmt.Errorf("resolve cwd: %w", err)
		}
	}

	envLookup := envLookup(opts.Env)
	loaded, err := config.Load(config.LoadOptions{
		CWD:          cwd,
		ExplicitPath: global.ConfigPath,
		LookupEnv:    envLookup,
	})
	if err != nil {
		return nil, err
	}

	if global.Workspace != "" {
		loaded.Config.WorkspaceRoot = global.Workspace
	}
	if global.Format != "" {
		loaded.Config.OutputFormat = global.Format
	}

	if err := config.Validate(loaded.Config); err != nil {
		return nil, err
	}

	bundle, err := boot.Bootstrap(loaded.Config)
	if err != nil {
		return nil, err
	}

	ctx := &commands.Context{
		Config:       loaded.Config,
		ConfigSource: loaded.Source,
		Warnings:     loaded.Warnings,
		Bundle:       bundle,
		Stdout:       stdout,
		Stderr:       stderr,
		OutputFormat: loaded.Config.OutputFormat,
		CWD:          cwd,
	}
	return ctx, nil
}

func envLookup(env []string) func(string) string {
	lookup := map[string]string{}
	for _, kv := range env {
		parts := strings.SplitN(kv, "=", 2)
		if len(parts) != 2 {
			continue
		}
		lookup[parts[0]] = parts[1]
	}
	return func(key string) string {
		return lookup[key]
	}
}
