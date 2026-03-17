package commands

import (
	"flag"
	"fmt"

	"srosv2/internal/shared/config"
)

func newConfigCommand() *Command {
	return &Command{
		Name:    "config",
		Summary: "Show resolved config and source",
		Usage:   "sros config [--validate]",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("config", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			validateOnly := fs.Bool("validate", false, "validate config and return")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if fs.NArg() != 0 {
				return OperatorError("config does not accept positional arguments")
			}

			if *validateOnly {
				if err := config.Validate(ctx.Config); err != nil {
					return ConfigError(err.Error())
				}
				return writeOutput(ctx, "config: valid", map[string]any{"valid": true})
			}

			payload := map[string]any{"source": ctx.ConfigSource, "warnings": ctx.Warnings, "config": ctx.Config}
			text := fmt.Sprintf("source: %s\nmode: %s\nworkspace_root: %s\nartifact_root: %s\npolicy_bundle_path: %s\nmemory_store_path: %s\ntrace_store_path: %s\noutput_format: %s",
				ctx.ConfigSource,
				ctx.Config.Mode,
				ctx.Config.WorkspaceRoot,
				ctx.Config.ArtifactRoot,
				ctx.Config.PolicyBundlePath,
				ctx.Config.MemoryStorePath,
				ctx.Config.TraceStorePath,
				ctx.Config.OutputFormat,
			)
			return writeOutput(ctx, text, payload)
		},
	}
}
