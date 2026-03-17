package commands

import "fmt"

func newStatusCommand() *Command {
	return &Command{
		Name:    "status",
		Summary: "Show CLI mode, config source, and boundary wiring",
		Usage:   "sros status",
		Run: func(ctx *Context, args []string) error {
			if err := requireNoArgs(args); err != nil {
				return err
			}
			payload := map[string]any{
				"mode":          ctx.Bundle.Mode,
				"config_source": ctx.ConfigSource,
				"workspace":     ctx.Config.WorkspaceRoot,
				"boundaries":    ctx.Bundle.Boundaries,
			}
			text := fmt.Sprintf("mode: %s\nconfig_source: %s\nworkspace: %s\n%s", ctx.Bundle.Mode, ctx.ConfigSource, ctx.Config.WorkspaceRoot, formatBoundaries(ctx.Bundle.Boundaries))
			return writeOutput(ctx, text, payload)
		},
	}
}
