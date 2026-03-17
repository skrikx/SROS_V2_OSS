package commands

import (
	"context"
	"flag"
	"strings"

	"srosv2/internal/core/runtime"
)

func newRunCommand() *Command {
	return &Command{
		Name:    "run",
		Summary: "Admit a canonical run contract into SR9 runtime",
		Usage:   "sros run --contract <run-contract.json>",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("run", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			contractPath := fs.String("contract", "", "path to canonical run contract json")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if fs.NArg() != 0 {
				return OperatorError("run does not accept positional arguments")
			}
			if strings.TrimSpace(*contractPath) == "" {
				return OperatorError("run requires --contract")
			}
			if ctx.Bundle.Runtime == nil {
				return DeferredError("run boundary is not wired")
			}
			resp, err := ctx.Bundle.Runtime.Run(context.Background(), runtime.RunRequest{ContractPath: strings.TrimSpace(*contractPath)})
			if err != nil {
				return EnvironmentError(err.Error())
			}
			return writeOutput(ctx, "runtime session admitted", resp)
		},
	}
}
