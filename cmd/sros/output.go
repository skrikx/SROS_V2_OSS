package main

import (
	"encoding/json"
	"fmt"
	"io"

	"srosv2/cmd/sros/commands"
)

type ExitCode int

const (
	ExitSuccess          ExitCode = 0
	ExitOperatorError    ExitCode = 2
	ExitConfigError      ExitCode = 3
	ExitEnvironmentError ExitCode = 4
	ExitDeferredError    ExitCode = 5
	ExitInternalError    ExitCode = 10
)

type errorEnvelope struct {
	Error string `json:"error"`
	Kind  string `json:"kind,omitempty"`
}

func renderError(ctx *commands.Context, err error) ExitCode {
	kind := commands.KindInternal
	message := err.Error()

	if cmdErr, ok := err.(*commands.CommandError); ok {
		kind = cmdErr.Kind
		message = cmdErr.Message
	}

	if ctx != nil && ctx.OutputFormat == "json" {
		_ = json.NewEncoder(ctx.Stderr).Encode(errorEnvelope{Error: message, Kind: string(kind)})
	} else {
		writeLine(ctx.Stderr, fmt.Sprintf("error: %s", message))
	}

	switch kind {
	case commands.KindOperator:
		return ExitOperatorError
	case commands.KindConfig:
		return ExitConfigError
	case commands.KindEnvironment:
		return ExitEnvironmentError
	case commands.KindDeferred:
		return ExitDeferredError
	default:
		return ExitInternalError
	}
}

func writeLine(w io.Writer, line string) {
	if w == nil {
		return
	}
	_, _ = fmt.Fprintln(w, line)
}
