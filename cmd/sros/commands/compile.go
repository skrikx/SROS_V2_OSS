package commands

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"srosv2/internal/core/sr8"
)

func newCompileCommand() *Command {
	return &Command{
		Name:    "compile",
		Summary: "Compile intent through SR8 into canonical run contract",
		Usage:   "sros compile [--intent <text>|--input <path>] [--emit-srxml] [--format <text|json>]",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("compile", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})

			intent := fs.String("intent", "", "operator intent text")
			input := fs.String("input", "", "path to intent input file")
			operator := fs.String("operator", "", "override operator id")
			tenant := fs.String("tenant", "", "override tenant id")
			workspace := fs.String("workspace", "", "override workspace id")
			emitSRXML := fs.Bool("emit-srxml", false, "emit SRXML artifact")
			format := fs.String("format", "", "command output format: text|json")
			receiptPath := fs.String("receipt-path", "", "optional explicit path to copy compile receipt")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if fs.NArg() != 0 {
				return OperatorError("compile does not accept positional arguments")
			}
			if strings.TrimSpace(*intent) == "" && strings.TrimSpace(*input) == "" {
				return OperatorError("compile requires --intent or --input")
			}

			if ctx.Bundle.Compiler == nil {
				return DeferredError("compile boundary is not wired")
			}

			req := sr8.CompileRequest{
				Intent:       strings.TrimSpace(*intent),
				InputPath:    strings.TrimSpace(*input),
				OperatorID:   strings.TrimSpace(*operator),
				TenantID:     strings.TrimSpace(*tenant),
				WorkspaceID:  strings.TrimSpace(*workspace),
				ArtifactRoot: ctx.Config.ArtifactRoot,
				EmitSRXML:    *emitSRXML,
			}

			result, err := ctx.Bundle.Compiler.Compile(context.Background(), req)
			if err != nil {
				return EnvironmentError(err.Error())
			}

			if strings.TrimSpace(*receiptPath) != "" {
				if err := writeReceiptFile(*receiptPath, result.Receipt); err != nil {
					return EnvironmentError(err.Error())
				}
			}

			cmdFormat := strings.ToLower(strings.TrimSpace(*format))
			if cmdFormat != "" {
				if cmdFormat != "text" && cmdFormat != "json" {
					return OperatorError("--format must be text or json")
				}
				ctx.OutputFormat = cmdFormat
			}

			payload := map[string]any{
				"accepted":       result.Accepted,
				"summary":        result.Summary,
				"normalized":     result.Normalized,
				"classification": result.Classification,
				"topology":       result.Topology,
				"run_contract":   result.RunContract,
				"receipt":        result.Receipt,
				"srxml":          result.SRXML,
				"artifacts":      result.Artifacts,
				"output_dir":     result.OutputDir,
			}

			text := fmt.Sprintf(
				"compile accepted\nrun_id: %s\ndomain: %s\nrisk: %s\ntopology: %s\noutput_dir: %s",
				result.RunContract.RunID,
				result.Classification.Domain,
				result.Classification.Risk,
				result.Topology.Topology,
				result.OutputDir,
			)
			if *emitSRXML {
				text += "\nsrxml: emitted"
			}
			if strings.TrimSpace(*receiptPath) != "" {
				text += "\nreceipt_copy: " + filepath.Clean(*receiptPath)
			}

			return writeOutput(ctx, text, payload)
		},
	}
}

func writeReceiptFile(path string, receipt sr8.CompileReceipt) error {
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return fmt.Errorf("create receipt directory: %w", err)
	}
	data, err := json.MarshalIndent(receipt, "", "  ")
	if err != nil {
		return fmt.Errorf("encode receipt: %w", err)
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return fmt.Errorf("write receipt file: %w", err)
	}
	return nil
}
