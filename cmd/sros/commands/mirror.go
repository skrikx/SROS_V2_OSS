package commands

import (
	"flag"
	"strings"
)

func newMirrorCommand() *Command {
	cmd := &Command{
		Name:    "mirror",
		Summary: "Mirror witness and drift surfaces",
		Usage:   "sros mirror <status|witness>",
	}
	cmd.Subcommands = []*Command{
		{
			Name:    "status",
			Summary: "Generate drift and reflection status from a runtime snapshot",
			Usage:   "sros mirror status --input <path>",
			Run: func(ctx *Context, args []string) error {
				fs := flag.NewFlagSet("mirror status", flag.ContinueOnError)
				fs.SetOutput(ioDiscard{})
				input := fs.String("input", "", "runtime snapshot json")
				if err := fs.Parse(args); err != nil {
					return OperatorError(err.Error())
				}
				if strings.TrimSpace(*input) == "" {
					return OperatorError("mirror status requires --input")
				}
				if ctx.Bundle.Mirror == nil {
					return DeferredError("mirror plane is not wired")
				}
				data, err := ctx.Bundle.Mirror.StatusFromFile(strings.TrimSpace(*input))
				if err != nil {
					return EnvironmentError(err.Error())
				}
				return writeOutput(ctx, "mirror status captured", data)
			},
		},
		{
			Name:    "witness",
			Summary: "Emit a semantic witness event from a runtime snapshot",
			Usage:   "sros mirror witness --input <path>",
			Run: func(ctx *Context, args []string) error {
				fs := flag.NewFlagSet("mirror witness", flag.ContinueOnError)
				fs.SetOutput(ioDiscard{})
				input := fs.String("input", "", "witness input json")
				if err := fs.Parse(args); err != nil {
					return OperatorError(err.Error())
				}
				if strings.TrimSpace(*input) == "" {
					return OperatorError("mirror witness requires --input")
				}
				if ctx.Bundle.Mirror == nil {
					return DeferredError("mirror plane is not wired")
				}
				data, err := ctx.Bundle.Mirror.WitnessFromFile(strings.TrimSpace(*input))
				if err != nil {
					return EnvironmentError(err.Error())
				}
				return writeOutput(ctx, "mirror witness emitted", data)
			},
		},
	}
	return cmd
}
