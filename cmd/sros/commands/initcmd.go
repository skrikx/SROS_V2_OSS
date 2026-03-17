package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func newInitCommand() *Command {
	return &Command{
		Name:    "init",
		Summary: "Initialize local config for SROS v2",
		Usage:   "sros init [--path <config-file>] [--force]",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("init", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			targetPath := fs.String("path", filepath.Join(ctx.Config.WorkspaceRoot, "sros.yaml"), "config file path")
			force := fs.Bool("force", false, "overwrite existing config")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if fs.NArg() != 0 {
				return OperatorError("init does not accept positional arguments")
			}

			if _, err := os.Stat(*targetPath); err == nil && !*force {
				return OperatorError("config file already exists; use --force to overwrite")
			}

			if err := os.MkdirAll(filepath.Dir(*targetPath), 0o755); err != nil {
				return EnvironmentError(fmt.Sprintf("create config directory: %v", err))
			}

			content := fmt.Sprintf("mode: %s\nworkspace_root: %s\nartifact_root: %s\npolicy_bundle_path: %s\nmemory_store_path: %s\ntrace_store_path: %s\noutput_format: %s\n",
				ctx.Config.Mode,
				ctx.Config.WorkspaceRoot,
				ctx.Config.ArtifactRoot,
				ctx.Config.PolicyBundlePath,
				ctx.Config.MemoryStorePath,
				ctx.Config.TraceStorePath,
				ctx.Config.OutputFormat,
			)
			if err := os.WriteFile(*targetPath, []byte(content), 0o644); err != nil {
				return EnvironmentError(fmt.Sprintf("write config file: %v", err))
			}

			return writeOutput(ctx, "config initialized: "+*targetPath, map[string]any{"config_path": *targetPath})
		},
	}
}
