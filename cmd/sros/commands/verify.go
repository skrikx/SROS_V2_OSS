package commands

import (
	"os"
	"path/filepath"
)

func newVerifyCommand() *Command {
	return &Command{
		Name:    "verify",
		Summary: "Verify first-run readiness, persistence posture, docs, examples, and showcase surfaces",
		Usage:   "sros verify",
		Examples: []string{
			"sros verify",
		},
		Run: func(ctx *Context, args []string) error {
			if err := requireNoArgs(args); err != nil {
				return err
			}
			checks := []map[string]any{
				verifyPath("migrations", filepath.Join(ctx.CWD, "migrations")),
				verifyPath("docs_tree", filepath.Join(ctx.CWD, "docs")),
				verifyPath("examples", filepath.Join(ctx.CWD, "examples")),
				verifyPath("scripts", filepath.Join(ctx.CWD, "scripts")),
				verifyPath("artifacts", ctx.Config.ArtifactRoot),
				verifyPath("showcase", filepath.Join(ctx.CWD, "artifacts", "showcase")),
			}
			payload := map[string]any{
				"checks":      checks,
				"database":    ctx.Config.Database.Summary(),
				"persistence": persistenceSummary(ctx),
			}
			return writeOutput(ctx, "verification completed\nfocus: first-run readiness and shareable operator surfaces", payload)
		},
	}
}

func verifyPath(name, path string) map[string]any {
	_, err := os.Stat(path)
	return map[string]any{"name": name, "path": path, "present": err == nil}
}

func persistenceSummary(ctx *Context) map[string]any {
	if ctx.Bundle.Persistence == nil {
		return map[string]any{"enabled": false, "connected": false}
	}
	return ctx.Bundle.Persistence.Summary()
}
