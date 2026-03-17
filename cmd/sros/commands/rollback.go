package commands

import (
	"context"
	"flag"

	"srosv2/internal/core/runtime"
)

func newRollbackCommand() *Command {
	return &Command{
		Name:    "rollback",
		Summary: "Dispatch rollback boundary (runtime deferred)",
		Usage:   "sros rollback --session <id> --checkpoint <id>",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("rollback", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			sessionID := fs.String("session", "", "session identifier")
			checkpointID := fs.String("checkpoint", "", "checkpoint identifier")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if *sessionID == "" || *checkpointID == "" {
				return OperatorError("rollback requires --session and --checkpoint")
			}
			if ctx.Bundle.Runtime == nil {
				return DeferredError("rollback boundary is not wired yet (deferred to W05)")
			}
			resp, err := ctx.Bundle.Runtime.Rollback(context.Background(), runtime.RollbackRequest{SessionID: *sessionID, CheckpointID: *checkpointID})
			if err != nil {
				return EnvironmentError(err.Error())
			}
			return writeOutput(ctx, "rollback dispatched", resp)
		},
	}
}
