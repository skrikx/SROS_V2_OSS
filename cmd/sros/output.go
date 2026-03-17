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
		if hint := operatorHint(kind, message); hint != "" {
			writeLine(ctx.Stderr, "next: "+hint)
		}
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

func operatorHint(kind commands.ErrorKind, message string) string {
	switch kind {
	case commands.KindOperator:
		if message != "" {
			return "run the matching '--help' command for the exact flags and examples"
		}
	case commands.KindConfig:
		return "check .env.example, local config, or run 'sros verify' for bootstrap hints"
	case commands.KindEnvironment:
		return "check local runtime dependencies, file paths, or run 'sros verify' for readiness details"
	case commands.KindDeferred:
		return "use a wired command path or inspect status to confirm the active boundary set"
	}
	return ""
}

func writeLine(w io.Writer, line string) {
	if w == nil {
		return
	}
	_, _ = fmt.Fprintln(w, line)
}
