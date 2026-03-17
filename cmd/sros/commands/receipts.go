package commands

import (
	"encoding/json"
	"flag"
	"os"
	"strings"
)

func newReceiptsCommand() *Command {
	cmd := &Command{
		Name:    "receipts",
		Summary: "Receipt export and closure proof inspection surfaces",
		Usage:   "sros receipts <export|closure>",
		Examples: []string{
			"sros receipts export --input examples/provenance/receipt_bundle_min.json",
			"sros receipts closure --input examples/provenance/closure_proof_min.json",
		},
	}
	cmd.Subcommands = []*Command{
		{
			Name:    "export",
			Summary: "Export a provenance bundle for audit or replay support",
			Usage:   "sros receipts export --input <path>",
			Run: func(ctx *Context, args []string) error {
				fs := flag.NewFlagSet("receipts export", flag.ContinueOnError)
				fs.SetOutput(ioDiscard{})
				input := fs.String("input", "", "receipt bundle json")
				if err := fs.Parse(args); err != nil {
					return OperatorError(err.Error())
				}
				if strings.TrimSpace(*input) == "" {
					return missingFlagError("receipts export", "--input", "run 'sros receipts export --help'")
				}
				if ctx.Bundle.Provenance == nil {
					return DeferredError("provenance plane is not wired")
				}
				data, err := ctx.Bundle.Provenance.ExportBundle(strings.TrimSpace(*input))
				if err != nil {
					return EnvironmentError(err.Error())
				}
				return writeOutput(ctx, "receipt bundle exported\nfocus: audit and replay handoff artifact", data)
			},
		},
		{
			Name:    "closure",
			Summary: "Inspect a closure proof artifact",
			Usage:   "sros receipts closure --input <path>",
			Run: func(ctx *Context, args []string) error {
				fs := flag.NewFlagSet("receipts closure", flag.ContinueOnError)
				fs.SetOutput(ioDiscard{})
				input := fs.String("input", "", "closure proof json")
				if err := fs.Parse(args); err != nil {
					return OperatorError(err.Error())
				}
				if strings.TrimSpace(*input) == "" {
					return missingFlagError("receipts closure", "--input", "run 'sros receipts closure --help'")
				}
				data, err := os.ReadFile(strings.TrimSpace(*input))
				if err != nil {
					return EnvironmentError(err.Error())
				}
				var payload map[string]any
				if err := json.Unmarshal(data, &payload); err != nil {
					return EnvironmentError(err.Error())
				}
				return writeOutput(ctx, "closure proof inspected\nfocus: terminal proof and linked evidence", payload)
			},
		},
	}
	return cmd
}
