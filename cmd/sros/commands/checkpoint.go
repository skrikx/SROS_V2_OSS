package commands

import (
	"context"
	"flag"

	"srosv2/internal/core/runtime"
)

func newCheckpointCommand() *Command {
	return &Command{
		Name:    "checkpoint",
		Summary: "Dispatch checkpoint boundary (runtime deferred)",
		Usage:   "sros checkpoint --session <id> --stage <stage>",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("checkpoint", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			sessionID := fs.String("session", "", "session identifier")
			stage := fs.String("stage", "", "checkpoint stage")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if *sessionID == "" || *stage == "" {
				return OperatorError("checkpoint requires --session and --stage")
			}
			if ctx.Bundle.Runtime == nil {
				return DeferredError("checkpoint boundary is not wired yet (deferred to W05)")
			}
			resp, err := ctx.Bundle.Runtime.Checkpoint(context.Background(), runtime.CheckpointRequest{SessionID: *sessionID, Stage: *stage})
			if err != nil {
				return EnvironmentError(err.Error())
			}
			return writeOutput(ctx, "checkpoint dispatched", resp)
		},
	}
}
