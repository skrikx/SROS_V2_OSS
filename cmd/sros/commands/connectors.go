package commands

import (
	"context"
	"flag"
)

func newConnectorsCommand() *Command {
	cmd := &Command{
		Name:    "connectors",
		Summary: "Governed connector capability surfaces",
		Usage:   "sros connectors <list|envelope>",
	}
	cmd.Subcommands = []*Command{
		{
			Name:    "list",
			Summary: "List registered connectors",
			Usage:   "sros connectors list",
			Run: func(ctx *Context, args []string) error {
				if err := requireNoArgs(args); err != nil {
					return err
				}
				if ctx.Bundle.Fabric == nil {
					return DeferredError("connectors list boundary is not wired")
				}
				data, err := ctx.Bundle.Fabric.ConnectorsList(context.Background())
				if err != nil {
					return EnvironmentError(err.Error())
				}
				return writeOutput(ctx, "connectors governed", data)
			},
		},
		{
			Name:    "envelope",
			Summary: "Inspect a connector auth envelope",
			Usage:   "sros connectors envelope inspect --input <path>",
			Subcommands: []*Command{
				{
					Name:    "inspect",
					Summary: "Inspect a redacted connector auth envelope",
					Usage:   "sros connectors envelope inspect --input <path>",
					Run: func(ctx *Context, args []string) error {
						fs := flag.NewFlagSet("connectors envelope inspect", flag.ContinueOnError)
						fs.SetOutput(ioDiscard{})
						input := fs.String("input", "", "envelope input path")
						if err := fs.Parse(args); err != nil {
							return OperatorError(err.Error())
						}
						if *input == "" {
							return OperatorError("connectors envelope inspect requires --input")
						}
						if ctx.Bundle.Fabric == nil {
							return DeferredError("connectors envelope inspect boundary is not wired")
						}
						data, err := ctx.Bundle.Fabric.ConnectorsInspectEnvelope(context.Background(), *input)
						if err != nil {
							return EnvironmentError(err.Error())
						}
						return writeOutput(ctx, "connector envelope inspected", data)
					},
				},
			},
		},
	}
	return cmd
}
