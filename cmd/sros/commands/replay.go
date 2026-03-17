package commands

import (
	"flag"
	"os"
	"path/filepath"
	"strings"

	"srosv2/internal/shared/ids"
)

func newReplayCommand() *Command {
	return &Command{
		Name:    "replay",
		Summary: "Replay stored evidence or replay fixtures",
		Usage:   "sros replay [--run-id <id>] [--input <path>]",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("replay", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			runID := fs.String("run-id", "", "stored run id")
			input := fs.String("input", "", "replay input file")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if ctx.Bundle.Trace == nil {
				return DeferredError("trace plane is not wired")
			}
			if strings.TrimSpace(*input) != "" {
				data, err := os.ReadFile(strings.TrimSpace(*input))
				if err != nil {
					return EnvironmentError(err.Error())
				}
				out := filepath.Join(ctx.Config.ArtifactRoot, "replays", filepath.Base(*input))
				if err := os.WriteFile(out, data, 0o644); err != nil {
					return EnvironmentError(err.Error())
				}
				return writeOutput(ctx, "replay artifact copied", map[string]any{"input": *input, "replay_artifact": out})
			}
			if strings.TrimSpace(*runID) == "" {
				return OperatorError("replay requires --run-id or --input")
			}
			result, err := ctx.Bundle.Trace.Replay.Replay(ids.RunID(strings.TrimSpace(*runID)))
			if err != nil {
				return EnvironmentError(err.Error())
			}
			return writeOutput(ctx, "replay completed", map[string]any{"replay": result})
		},
	}
}
