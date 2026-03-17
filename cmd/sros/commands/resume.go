package commands

import (
	"context"
	"flag"
	"strings"

	"srosv2/internal/core/runtime"
)

func newResumeCommand() *Command {
	return &Command{
		Name:    "resume",
		Summary: "Resume a runtime session",
		Usage:   "sros resume [--session <session-id>] [--approval-file <path>]",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("resume", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			sessionID := fs.String("session", "", "session identifier")
			approvalFile := fs.String("approval-file", "", "path to local approval artifact json")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if fs.NArg() != 0 {
				return OperatorError("resume does not accept positional arguments")
			}
			if ctx.Bundle.Runtime == nil {
				return DeferredError("resume boundary is not wired")
			}
			resp, err := ctx.Bundle.Runtime.Resume(context.Background(), runtime.ResumeRequest{
				SessionID:    strings.TrimSpace(*sessionID),
				ApprovalFile: strings.TrimSpace(*approvalFile),
			})
			if err != nil {
				return EnvironmentError(err.Error())
			}
			return writeOutput(ctx, "runtime resumed", resp)
		},
	}
}
