package commands

import (
	"context"
	"flag"

	"srosv2/internal/core/runtime"
)

func newPauseCommand() *Command {
	return &Command{
		Name:    "pause",
		Summary: "Dispatch pause boundary (runtime deferred)",
		Usage:   "sros pause --session <session-id> [--reason <text>]",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("pause", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			sessionID := fs.String("session", "", "session identifier")
			reason := fs.String("reason", "", "pause reason")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if *sessionID == "" {
				return OperatorError("pause requires --session")
			}
			if ctx.Bundle.Runtime == nil {
				return DeferredError("pause boundary is not wired yet (deferred to W05)")
			}
			resp, err := ctx.Bundle.Runtime.Pause(context.Background(), runtime.PauseRequest{SessionID: *sessionID, Reason: *reason})
			if err != nil {
				return EnvironmentError(err.Error())
			}
			return writeOutput(ctx, "pause dispatched", resp)
		},
	}
}
