package commands

import (
	"flag"
	"strings"

	"srosv2/internal/shared/ids"
)

func newTraceCommand() *Command {
	cmd := &Command{
		Name:    "trace",
		Summary: "Trace evidence surfaces",
		Usage:   "sros trace <inspect|replay>",
	}
	cmd.Subcommands = []*Command{
		{
			Name:    "inspect",
			Summary: "Inspect a trace example or stored run lineage",
			Usage:   "sros trace inspect [--input <path>] [--run-id <run-id>]",
			Run: func(ctx *Context, args []string) error {
				fs := flag.NewFlagSet("trace inspect", flag.ContinueOnError)
				fs.SetOutput(ioDiscard{})
				input := fs.String("input", "", "trace input json")
				runID := fs.String("run-id", "", "stored run id")
				if err := fs.Parse(args); err != nil {
					return OperatorError(err.Error())
				}
				if ctx.Bundle.Trace == nil {
					return DeferredError("trace plane is not wired")
				}
				if strings.TrimSpace(*input) != "" {
					data, err := ctx.Bundle.Trace.InspectFromFile(strings.TrimSpace(*input))
					if err != nil {
						return EnvironmentError(err.Error())
					}
					return writeOutput(ctx, "trace inspected", data)
				}
				if strings.TrimSpace(*runID) == "" {
					return OperatorError("trace inspect requires --input or --run-id")
				}
				events, err := ctx.Bundle.Trace.Reader.Events(ids.RunID(strings.TrimSpace(*runID)))
				if err != nil {
					return EnvironmentError(err.Error())
				}
				return writeOutput(ctx, "trace inspected", map[string]any{"run_id": strings.TrimSpace(*runID), "events": events})
			},
		},
		{
			Name:    "replay",
			Summary: "Replay a run lineage from append-only trace events",
			Usage:   "sros trace replay [--input <path>] [--run-id <run-id>]",
			Run: func(ctx *Context, args []string) error {
				fs := flag.NewFlagSet("trace replay", flag.ContinueOnError)
				fs.SetOutput(ioDiscard{})
				input := fs.String("input", "", "replay case json")
				runID := fs.String("run-id", "", "stored run id")
				if err := fs.Parse(args); err != nil {
					return OperatorError(err.Error())
				}
				if ctx.Bundle.Trace == nil {
					return DeferredError("trace plane is not wired")
				}
				if strings.TrimSpace(*input) != "" {
					data, err := ctx.Bundle.Trace.InspectFromFile(strings.TrimSpace(*input))
					if err != nil {
						return EnvironmentError(err.Error())
					}
					if run, ok := data["run_id"].(string); ok && run != "" {
						runID = &run
					}
				}
				if strings.TrimSpace(*runID) == "" {
					return OperatorError("trace replay requires --input or --run-id")
				}
				result, err := ctx.Bundle.Trace.Replay.Replay(ids.RunID(strings.TrimSpace(*runID)))
				if err != nil {
					return EnvironmentError(err.Error())
				}
				return writeOutput(ctx, "trace replay completed", map[string]any{"replay": result})
			},
		},
	}
	return cmd
}
