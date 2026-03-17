package commands

import (
	"context"
	"flag"
)

func newMemoryCommand() *Command {
	return &Command{
		Name:    "memory",
		Summary: "Inspect memory boundary (memory plane deferred)",
		Usage:   "sros memory --run-id <run-id>",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("memory", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			runID := fs.String("run-id", "", "run identifier")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if ctx.Bundle.Inspector == nil {
				return DeferredError("memory boundary is not wired yet (deferred to W07)")
			}
			data, err := ctx.Bundle.Inspector.Memory(context.Background(), *runID)
			if err != nil {
				return EnvironmentError(err.Error())
			}
			return writeOutput(ctx, "memory inspected", data)
		},
	}
}
