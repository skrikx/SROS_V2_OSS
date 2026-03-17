package commands

import (
	"context"
	"flag"
	"fmt"
	"strings"

	"srosv2/internal/core/runtime"
)

func newStatusCommand() *Command {
	return &Command{
		Name:    "status",
		Summary: "Show runtime status and boundary wiring",
		Usage:   "sros status [--session <id>] [--latest]",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("status", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			sessionID := fs.String("session", "", "runtime session identifier")
			latest := fs.Bool("latest", false, "show latest runtime session status")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if fs.NArg() != 0 {
				return OperatorError("status does not accept positional arguments")
			}

			snapshot := runtime.StatusSnapshot{
				Mode:    string(ctx.Bundle.Mode),
				Summary: "runtime boundary unavailable",
			}
			if ctx.Bundle.Runtime != nil {
				var err error
				snapshot, err = ctx.Bundle.Runtime.Status(context.Background(), runtime.StatusRequest{
					SessionID: strings.TrimSpace(*sessionID),
					Latest:    *latest,
				})
				if err != nil {
					return EnvironmentError(err.Error())
				}
			}

			payload := map[string]any{
				"mode":             ctx.Bundle.Mode,
				"config_source":    ctx.ConfigSource,
				"workspace":        ctx.Config.WorkspaceRoot,
				"boundaries":       ctx.Bundle.Boundaries,
				"runtime":          snapshot,
				"memory_wired":     ctx.Bundle.Memory != nil,
				"mirror_wired":     ctx.Bundle.Mirror != nil,
				"trace_wired":      ctx.Bundle.Trace != nil,
				"provenance_wired": ctx.Bundle.Provenance != nil,
			}
			text := fmt.Sprintf(
				"mode: %s\nconfig_source: %s\nworkspace: %s\nruntime_summary: %s\nruntime_session: %s\nruntime_state: %s\nlatest_checkpoint: %s\nlatest_rollback: %s\nwaiting_approval: %s\nlatest_mutation: %s\nlatest_witness: %s\nmemory_wired: %t\nmirror_wired: %t\ntrace_wired: %t\nprovenance_wired: %t\n%s",
				ctx.Bundle.Mode,
				ctx.ConfigSource,
				ctx.Config.WorkspaceRoot,
				emptyFallback(snapshot.Summary, "(none)"),
				sessionIDFromSnapshot(snapshot),
				stateFromSnapshot(snapshot),
				emptyFallback(snapshot.LatestCheckpointID, "(none)"),
				emptyFallback(snapshot.LatestRollbackID, "(none)"),
				emptyFallback(snapshot.WaitingApproval, "(none)"),
				emptyFallback(snapshot.LatestMutationID, "(none)"),
				emptyFallback(snapshot.LatestWitnessID, "(none)"),
				ctx.Bundle.Memory != nil,
				ctx.Bundle.Mirror != nil,
				ctx.Bundle.Trace != nil,
				ctx.Bundle.Provenance != nil,
				formatBoundaries(ctx.Bundle.Boundaries),
			)
			return writeOutput(ctx, text, payload)
		},
	}
}

func sessionIDFromSnapshot(snapshot runtime.StatusSnapshot) string {
	if snapshot.Session == nil || snapshot.Session.SessionID == "" {
		return "(none)"
	}
	return snapshot.Session.SessionID
}

func stateFromSnapshot(snapshot runtime.StatusSnapshot) string {
	if snapshot.Session == nil || snapshot.Session.State == "" {
		return "(none)"
	}
	return string(snapshot.Session.State)
}

func emptyFallback(v, fallback string) string {
	if strings.TrimSpace(v) == "" {
		return fallback
	}
	return v
}
