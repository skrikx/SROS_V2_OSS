package commands

import (
	"context"
	"flag"
)

func newMirrorCommand() *Command {
	return &Command{
		Name:    "mirror",
		Summary: "Inspect mirror boundary (mirror plane deferred)",
		Usage:   "sros mirror --run-id <run-id>",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("mirror", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			runID := fs.String("run-id", "", "run identifier")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if ctx.Bundle.Inspector == nil {
				return DeferredError("mirror boundary is not wired yet (deferred to W07)")
			}
			data, err := ctx.Bundle.Inspector.Mirror(context.Background(), *runID)
			if err != nil {
				return EnvironmentError(err.Error())
			}
			return writeOutput(ctx, "mirror inspected", data)
		},
	}
}
