package commands

import (
	"context"
	"flag"
	"strings"

	"srosv2/internal/core/runtime"
)

func newRollbackCommand() *Command {
	return &Command{
		Name:    "rollback",
		Summary: "Rollback a runtime session to a checkpoint",
		Usage:   "sros rollback [--session <id>] [--checkpoint <id>] [--reason <text>]",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("rollback", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			sessionID := fs.String("session", "", "session identifier")
			checkpointID := fs.String("checkpoint", "", "checkpoint identifier")
			reason := fs.String("reason", "", "rollback reason")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if fs.NArg() != 0 {
				return OperatorError("rollback does not accept positional arguments")
			}
			if ctx.Bundle.Runtime == nil {
				return DeferredError("rollback boundary is not wired")
			}
			resp, err := ctx.Bundle.Runtime.Rollback(context.Background(), runtime.RollbackRequest{
				SessionID:    strings.TrimSpace(*sessionID),
				CheckpointID: strings.TrimSpace(*checkpointID),
				Reason:       strings.TrimSpace(*reason),
			})
			if err != nil {
				return EnvironmentError(err.Error())
			}
			return writeOutput(ctx, "runtime rolled back", resp)
		},
	}
}
