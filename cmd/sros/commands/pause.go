package commands

import (
	"context"
	"flag"
	"strings"

	"srosv2/internal/core/runtime"
)

func newPauseCommand() *Command {
	return &Command{
		Name:    "pause",
		Summary: "Pause a runtime session",
		Usage:   "sros pause [--session <session-id>] [--reason <text>]",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("pause", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			sessionID := fs.String("session", "", "session identifier")
			reason := fs.String("reason", "", "pause reason")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if fs.NArg() != 0 {
				return OperatorError("pause does not accept positional arguments")
			}
			if ctx.Bundle.Runtime == nil {
				return DeferredError("pause boundary is not wired")
			}
			resp, err := ctx.Bundle.Runtime.Pause(context.Background(), runtime.PauseRequest{
				SessionID: strings.TrimSpace(*sessionID),
				Reason:    strings.TrimSpace(*reason),
			})
			if err != nil {
				return EnvironmentError(err.Error())
			}
			return writeOutput(ctx, "runtime paused", resp)
		},
	}
}
