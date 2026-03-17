package commands

import (
	"context"
	"flag"

	"srosv2/internal/core/runtime"
)

func newResumeCommand() *Command {
	return &Command{
		Name:    "resume",
		Summary: "Dispatch resume boundary (runtime deferred)",
		Usage:   "sros resume --session <session-id>",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("resume", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			sessionID := fs.String("session", "", "session identifier")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if *sessionID == "" {
				return OperatorError("resume requires --session")
			}
			if ctx.Bundle.Runtime == nil {
				return DeferredError("resume boundary is not wired yet (deferred to W05)")
			}
			resp, err := ctx.Bundle.Runtime.Resume(context.Background(), runtime.ResumeRequest{SessionID: *sessionID})
			if err != nil {
				return EnvironmentError(err.Error())
			}
			return writeOutput(ctx, "resume dispatched", resp)
		},
	}
}
