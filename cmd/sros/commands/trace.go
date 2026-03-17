package commands

import (
	"context"
	"flag"
)

func newTraceCommand() *Command {
	return &Command{
		Name:    "trace",
		Summary: "Inspect trace boundary (trace plane deferred)",
		Usage:   "sros trace --run-id <run-id>",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("trace", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			runID := fs.String("run-id", "", "run identifier")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if ctx.Bundle.Inspector == nil {
				return DeferredError("trace boundary is not wired yet (deferred to W08)")
			}
			data, err := ctx.Bundle.Inspector.Trace(context.Background(), *runID)
			if err != nil {
				return EnvironmentError(err.Error())
			}
			return writeOutput(ctx, "trace inspected", data)
		},
	}
}
