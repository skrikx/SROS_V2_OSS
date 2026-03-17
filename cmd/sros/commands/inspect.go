package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func newInspectCommand() *Command {
	return &Command{
		Name:    "inspect",
		Summary: "Inspect local wiring and repository surfaces",
		Usage:   "sros inspect [--path <relative-path>]",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("inspect", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			relPath := fs.String("path", "", "relative path to inspect")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}

			if *relPath != "" {
				target := filepath.Join(ctx.Config.WorkspaceRoot, *relPath)
				info, err := os.Stat(target)
				if err != nil {
					return EnvironmentError(err.Error())
				}
				payload := map[string]any{"path": target, "is_dir": info.IsDir(), "size_bytes": info.Size()}
				return writeOutput(ctx, fmt.Sprintf("inspect: %s", target), payload)
			}

			required := []string{"cmd/sros", "internal/core/boot", "internal/core/runtime", "contracts", "docs", "tests"}
			report := map[string]bool{}
			for _, rel := range required {
				_, err := os.Stat(filepath.Join(ctx.Config.WorkspaceRoot, rel))
				report[rel] = err == nil
			}
			payload := map[string]any{"workspace_root": ctx.Config.WorkspaceRoot, "required_paths": report, "mode": ctx.Bundle.Mode}
			return writeOutput(ctx, "inspect: repository wiring snapshot captured", payload)
		},
	}
}
