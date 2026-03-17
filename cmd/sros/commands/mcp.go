package commands

import (
	"context"
	"flag"
)

func newMCPCommand() *Command {
	cmd := &Command{
		Name:    "mcp",
		Summary: "Governed MCP capability surface",
		Usage:   "sros mcp ingest --file <path>",
	}
	cmd.Subcommands = []*Command{
		{
			Name:    "ingest",
			Summary: "Ingest MCP payload file",
			Usage:   "sros mcp ingest --file <path>",
			Run: func(ctx *Context, args []string) error {
				fs := flag.NewFlagSet("mcp ingest", flag.ContinueOnError)
				fs.SetOutput(ioDiscard{})
				file := fs.String("file", "", "input file path")
				if err := fs.Parse(args); err != nil {
					return OperatorError(err.Error())
				}
				if *file == "" {
					return OperatorError("mcp ingest requires --file")
				}
				if ctx.Bundle.Fabric == nil {
					return DeferredError("mcp ingest boundary is not wired")
				}
				data, err := ctx.Bundle.Fabric.MCPIngest(context.Background(), *file)
				if err != nil {
					return EnvironmentError(err.Error())
				}
				return writeOutput(ctx, "mcp capability governed", data)
			},
		},
	}
	return cmd
}
