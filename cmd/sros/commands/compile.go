package commands

import (
	"context"
	"flag"

	"srosv2/internal/core/runtime"
)

func newCompileCommand() *Command {
	return &Command{
		Name:    "compile",
		Summary: "Parse inputs and dispatch compile boundary (SR8 deferred)",
		Usage:   "sros compile --intent <text> [--input <path>]",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("compile", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			intent := fs.String("intent", "", "operator intent text")
			input := fs.String("input", "", "path to intent input file")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if *intent == "" && *input == "" {
				return OperatorError("compile requires --intent or --input")
			}
			if ctx.Bundle.Compiler == nil {
				return DeferredError("compile boundary is not wired yet (deferred to W04)")
			}

			resp, err := ctx.Bundle.Compiler.Compile(context.Background(), runtime.CompileRequest{Intent: *intent, InputPath: *input})
			if err != nil {
				return EnvironmentError(err.Error())
			}
			return writeOutput(ctx, "compile dispatched", resp)
		},
	}
}
