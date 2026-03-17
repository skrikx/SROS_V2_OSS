package commands

import (
	"context"
	"flag"
	"strings"

	"srosv2/internal/core/runtime"
)

func newPlanCommand() *Command {
	return &Command{
		Name:    "plan",
		Summary: "Run runtime preflight planning on a canonical run contract",
		Usage:   "sros plan --contract <run-contract.json>",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("plan", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			contractPath := fs.String("contract", "", "path to canonical run contract json")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if fs.NArg() != 0 {
				return OperatorError("plan does not accept positional arguments")
			}
			if strings.TrimSpace(*contractPath) == "" {
				return OperatorError("plan requires --contract")
			}
			if ctx.Bundle.Runtime == nil {
				return DeferredError("plan boundary is not wired")
			}
			resp, err := ctx.Bundle.Runtime.Plan(context.Background(), runtime.RunRequest{ContractPath: strings.TrimSpace(*contractPath)})
			if err != nil {
				return EnvironmentError(err.Error())
			}
			return writeOutput(ctx, "runtime preflight planned", resp)
		},
	}
}
