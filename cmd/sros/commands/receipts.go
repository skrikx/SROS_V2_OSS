package commands

import (
	"context"
	"flag"
)

func newReceiptsCommand() *Command {
	return &Command{
		Name:    "receipts",
		Summary: "Inspect receipt boundary (evidence plane deferred)",
		Usage:   "sros receipts --run-id <run-id>",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("receipts", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			runID := fs.String("run-id", "", "run identifier")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if ctx.Bundle.Inspector == nil {
				return DeferredError("receipts boundary is not wired yet (deferred to W08)")
			}
			data, err := ctx.Bundle.Inspector.Receipts(context.Background(), *runID)
			if err != nil {
				return EnvironmentError(err.Error())
			}
			return writeOutput(ctx, "receipts inspected", data)
		},
	}
}
