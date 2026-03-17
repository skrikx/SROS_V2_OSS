package commands

import "context"

func newConnectorsCommand() *Command {
	cmd := &Command{
		Name:    "connectors",
		Summary: "Connector capability surfaces (fabric deferred)",
		Usage:   "sros connectors list",
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
					return DeferredError("connectors list boundary is not wired yet (deferred to W09)")
				}
				data, err := ctx.Bundle.Fabric.ConnectorsList(context.Background())
				if err != nil {
					return EnvironmentError(err.Error())
				}
				return writeOutput(ctx, "connectors listed", data)
			},
		},
	}
	return cmd
}
