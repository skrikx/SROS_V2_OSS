package commands

import (
	"context"
	"flag"
	"strings"

	"srosv2/internal/core/runtime"
)

func newCheckpointCommand() *Command {
	return &Command{
		Name:    "checkpoint",
		Summary: "Create a runtime checkpoint record",
		Usage:   "sros checkpoint [--session <id>] [--stage <stage>]",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("checkpoint", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			sessionID := fs.String("session", "", "session identifier")
			stage := fs.String("stage", "validated", "checkpoint stage")
			latest := fs.Bool("latest", false, "use latest runtime session")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if fs.NArg() != 0 {
				return OperatorError("checkpoint does not accept positional arguments")
			}
			if ctx.Bundle.Runtime == nil {
				return DeferredError("checkpoint boundary is not wired")
			}
			if *latest {
				*sessionID = ""
			}
			resp, err := ctx.Bundle.Runtime.Checkpoint(context.Background(), runtime.CheckpointRequest{
				SessionID: strings.TrimSpace(*sessionID),
				Stage:     strings.TrimSpace(*stage),
			})
			if err != nil {
				return EnvironmentError(err.Error())
			}
			return writeOutput(ctx, "checkpoint created", resp)
		},
	}
}
