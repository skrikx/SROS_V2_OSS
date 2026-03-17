package commands

import (
	"context"
	"flag"

	"srosv2/internal/core/runtime"
)

func newPlanCommand() *Command {
	return &Command{
		Name:    "plan",
		Summary: "Dispatch plan boundary (runtime planner deferred)",
		Usage:   "sros plan [--run-id <id>] [--plan <path>]",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("plan", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			runID := fs.String("run-id", "", "run identifier")
			plan := fs.String("plan", "", "plan path")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if ctx.Bundle.Runtime == nil {
				return DeferredError("plan boundary is not wired yet (deferred to W05)")
			}
			resp, err := ctx.Bundle.Runtime.Plan(context.Background(), runtime.RunRequest{RunID: *runID, Plan: *plan})
			if err != nil {
				return EnvironmentError(err.Error())
			}
			return writeOutput(ctx, "plan dispatched", resp)
		},
	}
}
