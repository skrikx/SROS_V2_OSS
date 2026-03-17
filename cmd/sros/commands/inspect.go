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
		Summary: "Inspect local wiring, readiness, and repository surfaces with operator-friendly snapshots",
		Usage:   "sros inspect [--path <relative-path>]",
		Examples: []string{
			"sros inspect",
			"sros inspect --path examples/showcase",
		},
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
				return writeOutput(ctx, fmt.Sprintf("inspect target: %s\nkind: %s\nsize_bytes: %d", target, fileKind(info.IsDir()), info.Size()), payload)
			}

			required := []string{"cmd/sros", "internal/core/boot", "internal/core/runtime", "internal/core/mem", "internal/core/mirror", "contracts", "docs", "tests"}
			report := map[string]bool{}
			for _, rel := range required {
				_, err := os.Stat(filepath.Join(ctx.Config.WorkspaceRoot, rel))
				report[rel] = err == nil
			}
			payload := map[string]any{
				"workspace_root":   ctx.Config.WorkspaceRoot,
				"required_paths":   report,
				"mode":             ctx.Bundle.Mode,
				"memory_wired":     ctx.Bundle.Memory != nil,
				"mirror_wired":     ctx.Bundle.Mirror != nil,
				"trace_wired":      ctx.Bundle.Trace != nil,
				"provenance_wired": ctx.Bundle.Provenance != nil,
			}
			return writeOutput(ctx, "inspect: repository wiring snapshot captured\nfocus: repository wiring and operator-visible boundaries", payload)
		},
	}
}

func fileKind(dir bool) string {
	if dir {
		return "directory"
	}
	return "file"
}
