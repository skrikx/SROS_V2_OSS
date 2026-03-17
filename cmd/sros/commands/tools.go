package commands

import (
	"context"
	"flag"
)

func newToolsCommand() *Command {
	cmd := &Command{
		Name:    "tools",
		Summary: "Tool capability surfaces (fabric deferred)",
		Usage:   "sros tools <list|show|validate|register>",
	}

	cmd.Subcommands = []*Command{
		{
			Name:    "list",
			Summary: "List known tool manifests",
			Usage:   "sros tools list",
			Run: func(ctx *Context, args []string) error {
				if err := requireNoArgs(args); err != nil {
					return err
				}
				if ctx.Bundle.Fabric == nil {
					return DeferredError("tools list boundary is not wired yet (deferred to W09)")
				}
				data, err := ctx.Bundle.Fabric.ToolsList(context.Background())
				if err != nil {
					return EnvironmentError(err.Error())
				}
				return writeOutput(ctx, "tools listed", data)
			},
		},
		{
			Name:    "show",
			Summary: "Show a single tool manifest",
			Usage:   "sros tools show --name <tool>",
			Run: func(ctx *Context, args []string) error {
				fs := flag.NewFlagSet("tools show", flag.ContinueOnError)
				fs.SetOutput(ioDiscard{})
				name := fs.String("name", "", "tool name")
				if err := fs.Parse(args); err != nil {
					return OperatorError(err.Error())
				}
				if *name == "" {
					return OperatorError("tools show requires --name")
				}
				if ctx.Bundle.Fabric == nil {
					return DeferredError("tools show boundary is not wired yet (deferred to W09)")
				}
				data, err := ctx.Bundle.Fabric.ToolsShow(context.Background(), *name)
				if err != nil {
					return EnvironmentError(err.Error())
				}
				return writeOutput(ctx, "tool shown", data)
			},
		},
		{
			Name:    "validate",
			Summary: "Validate a tool manifest file",
			Usage:   "sros tools validate --manifest <path>",
			Run: func(ctx *Context, args []string) error {
				fs := flag.NewFlagSet("tools validate", flag.ContinueOnError)
				fs.SetOutput(ioDiscard{})
				manifest := fs.String("manifest", "", "manifest path")
				if err := fs.Parse(args); err != nil {
					return OperatorError(err.Error())
				}
				if *manifest == "" {
					return OperatorError("tools validate requires --manifest")
				}
				if ctx.Bundle.Fabric == nil {
					return DeferredError("tools validate boundary is not wired yet (deferred to W09)")
				}
				data, err := ctx.Bundle.Fabric.ToolsValidate(context.Background(), *manifest)
				if err != nil {
					return EnvironmentError(err.Error())
				}
				return writeOutput(ctx, "tool manifest validated", data)
			},
		},
		{
			Name:    "register",
			Summary: "Register a tool manifest",
			Usage:   "sros tools register --manifest <path>",
			Run: func(ctx *Context, args []string) error {
				fs := flag.NewFlagSet("tools register", flag.ContinueOnError)
				fs.SetOutput(ioDiscard{})
				manifest := fs.String("manifest", "", "manifest path")
				if err := fs.Parse(args); err != nil {
					return OperatorError(err.Error())
				}
				if *manifest == "" {
					return OperatorError("tools register requires --manifest")
				}
				if ctx.Bundle.Fabric == nil {
					return DeferredError("tools register boundary is not wired yet (deferred to W09)")
				}
				data, err := ctx.Bundle.Fabric.ToolsRegister(context.Background(), *manifest)
				if err != nil {
					return EnvironmentError(err.Error())
				}
				return writeOutput(ctx, "tool manifest registered", data)
			},
		},
	}

	return cmd
}
