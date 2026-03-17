package commands

import (
	"context"
	"flag"

	"srosv2/internal/core/runtime"
)

func newRunCommand() *Command {
	return &Command{
		Name:    "run",
		Summary: "Dispatch run boundary (SR9 deferred)",
		Usage:   "sros run [--run-id <id>] [--plan <path>]",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("run", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			runID := fs.String("run-id", "", "run identifier")
			plan := fs.String("plan", "", "plan path")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if ctx.Bundle.Runtime == nil {
				return DeferredError("run boundary is not wired yet (deferred to W05)")
			}
			resp, err := ctx.Bundle.Runtime.Run(context.Background(), runtime.RunRequest{RunID: *runID, Plan: *plan})
			if err != nil {
				return EnvironmentError(err.Error())
			}
			return writeOutput(ctx, "run dispatched", resp)
		},
	}
}
