package commands

import (
	"context"
	"flag"

	"srosv2/contracts/release"
	"srosv2/internal/shared/ids"
)

func newReleaseCommand() *Command {
	cmd := &Command{
		Name:    "release",
		Summary: "Local release baseline surfaces",
		Usage:   "sros release pack --checkpoint <id>",
	}
	cmd.Subcommands = []*Command{
		{
			Name:    "pack",
			Summary: "Create a local release baseline artifact",
			Usage:   "sros release pack --checkpoint <id> [--stage promoted]",
			Run: func(ctx *Context, args []string) error {
				fs := flag.NewFlagSet("release pack", flag.ContinueOnError)
				fs.SetOutput(ioDiscard{})
				checkpoint := fs.String("checkpoint", "", "checkpoint id")
				stage := fs.String("stage", string(release.StagePromoted), "target stage")
				if err := fs.Parse(args); err != nil {
					return OperatorError(err.Error())
				}
				if *checkpoint == "" {
					return OperatorError("release pack requires --checkpoint")
				}
				if ctx.Bundle.Provenance == nil {
					return DeferredError("provenance plane is not wired")
				}
				data, err := ctx.Bundle.Provenance.PackRelease(context.Background(), ids.CheckpointID(*checkpoint), release.Stage(*stage))
				if err != nil {
					return EnvironmentError(err.Error())
				}
				return writeOutput(ctx, "release baseline packed", data)
			},
		},
	}
	return cmd
}
