package main

import (
	"srosv2/cmd/sros/commands"
)

func Execute(args []string, opts RunOptions) int {
	global, remaining, err := parseGlobalOptions(args)
	if err != nil {
		ctx := &commands.Context{Stderr: opts.Stderr, OutputFormat: "text"}
		return int(renderError(ctx, commands.OperatorError(err.Error())))
	}

	ctx, err := buildCommandContext(global, opts)
	if err != nil {
		errCtx := &commands.Context{Stderr: opts.Stderr, OutputFormat: "text"}
		return int(renderError(errCtx, commands.ConfigError(err.Error())))
	}

	root := commands.NewRootCommand()
	if global.Help && len(remaining) == 0 {
		commands.WriteHelp(ctx, root, nil)
		return int(ExitSuccess)
	}

	if len(remaining) == 0 {
		commands.WriteHelp(ctx, root, nil)
		return int(ExitSuccess)
	}

	if err := commands.Dispatch(root, ctx, remaining); err != nil {
		return int(renderError(ctx, err))
	}

	return int(ExitSuccess)
}

func CommandTree() *commands.Command {
	return commands.NewRootCommand()
}

func CommandPaths() []string {
	return commands.CommandPaths(NewRootCommand())
}

func NewRootCommand() *commands.Command {
	return commands.NewRootCommand()
}
