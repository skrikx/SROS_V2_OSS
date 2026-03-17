package commands

import (
	"flag"
	"strings"

	"srosv2/internal/shared/ids"
)

func newMemoryCommand() *Command {
	cmd := &Command{
		Name:    "memory",
		Summary: "Memory continuity surfaces",
		Usage:   "sros memory <recall|branch|rewind|status>",
	}
	cmd.Subcommands = []*Command{
		{
			Name:    "recall",
			Summary: "Import seed data and query the recall index",
			Usage:   "sros memory recall --input <path> [--query <text>]",
			Run: func(ctx *Context, args []string) error {
				fs := flag.NewFlagSet("memory recall", flag.ContinueOnError)
				fs.SetOutput(ioDiscard{})
				input := fs.String("input", "", "seed json")
				query := fs.String("query", "workspace", "recall query")
				if err := fs.Parse(args); err != nil {
					return OperatorError(err.Error())
				}
				if strings.TrimSpace(*input) == "" {
					return OperatorError("memory recall requires --input")
				}
				if ctx.Bundle.Memory == nil {
					return DeferredError("memory plane is not wired")
				}
				if _, err := ctx.Bundle.Memory.ImportSeed(strings.TrimSpace(*input)); err != nil {
					return EnvironmentError(err.Error())
				}
				data, err := ctx.Bundle.Memory.Recall(strings.TrimSpace(*query))
				if err != nil {
					return EnvironmentError(err.Error())
				}
				return writeOutput(ctx, "memory recall completed", data)
			},
		},
		{
			Name:    "branch",
			Summary: "Create a real branch lineage object and apply branch mutations",
			Usage:   "sros memory branch --input <path>",
			Run: func(ctx *Context, args []string) error {
				fs := flag.NewFlagSet("memory branch", flag.ContinueOnError)
				fs.SetOutput(ioDiscard{})
				input := fs.String("input", "", "branch lineage json")
				if err := fs.Parse(args); err != nil {
					return OperatorError(err.Error())
				}
				if strings.TrimSpace(*input) == "" {
					return OperatorError("memory branch requires --input")
				}
				if ctx.Bundle.Memory == nil {
					return DeferredError("memory plane is not wired")
				}
				data, err := ctx.Bundle.Memory.ApplyBranchFile(strings.TrimSpace(*input))
				if err != nil {
					return EnvironmentError(err.Error())
				}
				return writeOutput(ctx, "memory branch applied", data)
			},
		},
		{
			Name:    "rewind",
			Summary: "Rewind a branch to an explicit mutation id",
			Usage:   "sros memory rewind --branch <id> --mutation <id> --operator <id> [--reason <text>]",
			Run: func(ctx *Context, args []string) error {
				fs := flag.NewFlagSet("memory rewind", flag.ContinueOnError)
				fs.SetOutput(ioDiscard{})
				branchID := fs.String("branch", "", "branch id")
				mutationID := fs.String("mutation", "", "target mutation id")
				operatorID := fs.String("operator", "op_local", "operator id")
				reason := fs.String("reason", "manual rewind", "rewind reason")
				if err := fs.Parse(args); err != nil {
					return OperatorError(err.Error())
				}
				if strings.TrimSpace(*branchID) == "" || strings.TrimSpace(*mutationID) == "" {
					return OperatorError("memory rewind requires --branch and --mutation")
				}
				if ctx.Bundle.Memory == nil {
					return DeferredError("memory plane is not wired")
				}
				if err := ctx.Bundle.Memory.Rewind(
					ids.BranchID(strings.TrimSpace(*branchID)),
					ids.MemoryMutationID(strings.TrimSpace(*mutationID)),
					ids.OperatorID(strings.TrimSpace(*operatorID)),
					strings.TrimSpace(*reason),
				); err != nil {
					return EnvironmentError(err.Error())
				}
				branches, err := ctx.Bundle.Memory.Branches()
				if err != nil {
					return EnvironmentError(err.Error())
				}
				return writeOutput(ctx, "memory rewind applied", map[string]any{"branches": branches})
			},
		},
		{
			Name:    "status",
			Summary: "Show memory store, lineage, and branch status",
			Usage:   "sros memory status [--query <text>]",
			Run: func(ctx *Context, args []string) error {
				fs := flag.NewFlagSet("memory status", flag.ContinueOnError)
				fs.SetOutput(ioDiscard{})
				query := fs.String("query", "workspace", "recall query")
				if err := fs.Parse(args); err != nil {
					return OperatorError(err.Error())
				}
				if ctx.Bundle.Memory == nil {
					return DeferredError("memory plane is not wired")
				}
				data, err := ctx.Bundle.Memory.Recall(strings.TrimSpace(*query))
				if err != nil {
					return EnvironmentError(err.Error())
				}
				branches, err := ctx.Bundle.Memory.Branches()
				if err != nil {
					return EnvironmentError(err.Error())
				}
				tree, err := ctx.Bundle.Memory.SessionTree()
				if err != nil {
					return EnvironmentError(err.Error())
				}
				payload := map[string]any{"recall": data, "branches": branches, "session_tree": tree}
				return writeOutput(ctx, "memory status captured", payload)
			},
		},
	}
	return cmd
}
