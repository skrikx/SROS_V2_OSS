package commands

import (
	"os"
	"path/filepath"
)

func newScaffoldCommand() *Command {
	return &Command{
		Name:    "scaffold",
		Summary: "Scaffold local artifact and persistence directories",
		Usage:   "sros scaffold",
		Run: func(ctx *Context, args []string) error {
			if err := requireNoArgs(args); err != nil {
				return err
			}
			dirs := []string{
				filepath.Join(ctx.Config.ArtifactRoot, "receipts"),
				filepath.Join(ctx.Config.ArtifactRoot, "bundles"),
				filepath.Join(ctx.Config.ArtifactRoot, "releases"),
				filepath.Join(ctx.Config.ArtifactRoot, "replays"),
			}
			for _, dir := range dirs {
				if err := os.MkdirAll(dir, 0o755); err != nil {
					return EnvironmentError(err.Error())
				}
			}
			return writeOutput(ctx, "local scaffold prepared", map[string]any{"dirs": dirs})
		},
	}
}
