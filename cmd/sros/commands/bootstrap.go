package commands

import "fmt"

func newBootstrapCommand() *Command {
	return &Command{
		Name:    "bootstrap",
		Summary: "Show bootstrap mode and service boundaries",
		Usage:   "sros bootstrap",
		Run: func(ctx *Context, args []string) error {
			if err := requireNoArgs(args); err != nil {
				return err
			}
			payload := map[string]any{"mode": ctx.Bundle.Mode, "boundaries": ctx.Bundle.Boundaries}
			text := fmt.Sprintf("mode: %s\n%s", ctx.Bundle.Mode, formatBoundaries(ctx.Bundle.Boundaries))
			return writeOutput(ctx, text, payload)
		},
	}
}
