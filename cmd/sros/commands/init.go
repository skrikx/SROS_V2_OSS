package commands

import (
	"encoding/json"
	"fmt"
	"io"
	"sort"
	"strings"

	"srosv2/internal/core/boot"
	"srosv2/internal/shared/config"
)

type Context struct {
	Config       config.Config
	ConfigSource string
	Warnings     []string
	Bundle       boot.Bundle
	Stdout       io.Writer
	Stderr       io.Writer
	OutputFormat string
	CWD          string
}

type Handler func(*Context, []string) error

type Command struct {
	Name        string
	Summary     string
	Usage       string
	Examples    []string
	Aliases     []string
	Subcommands []*Command
	Run         Handler
}

type ErrorKind string

const (
	KindOperator    ErrorKind = "operator"
	KindConfig      ErrorKind = "config"
	KindEnvironment ErrorKind = "environment"
	KindDeferred    ErrorKind = "deferred"
	KindInternal    ErrorKind = "internal"
)

type CommandError struct {
	Kind    ErrorKind
	Message string
}

func (e *CommandError) Error() string { return e.Message }

func OperatorError(message string) error { return &CommandError{Kind: KindOperator, Message: message} }
func ConfigError(message string) error   { return &CommandError{Kind: KindConfig, Message: message} }
func EnvironmentError(message string) error {
	return &CommandError{Kind: KindEnvironment, Message: message}
}
func DeferredError(message string) error { return &CommandError{Kind: KindDeferred, Message: message} }
func InternalError(message string) error { return &CommandError{Kind: KindInternal, Message: message} }

func NewRootCommand() *Command {
	root := &Command{
		Name:    "sros",
		Summary: "SROS V2 local-only operator CLI",
		Usage:   "sros <command> [flags]",
	}
	root.Subcommands = []*Command{
		newInitCommand(),
		newBootstrapCommand(),
		newDoctorCommand(),
		newSeedCommand(),
		newConfigCommand(),
		newCompileCommand(),
		newRunCommand(),
		newPlanCommand(),
		newResumeCommand(),
		newPauseCommand(),
		newCheckpointCommand(),
		newRollbackCommand(),
		newTraceCommand(),
		newReceiptsCommand(),
		newMemoryCommand(),
		newMirrorCommand(),
		newInspectCommand(),
		newStatusCommand(),
		newToolsCommand(),
		newConnectorsCommand(),
		newMCPCommand(),
	}
	return root
}

func Dispatch(root *Command, ctx *Context, args []string) error {
	if len(args) == 0 {
		WriteHelp(ctx, root, nil)
		return nil
	}
	if args[0] == "help" {
		return dispatchHelp(root, ctx, args[1:])
	}

	current := root
	path := []string{}
	idx := 0
	for idx < len(args) {
		tok := args[idx]
		if isHelpToken(tok) {
			WriteHelp(ctx, current, path)
			return nil
		}
		next := findSubcommand(current, tok)
		if next == nil {
			break
		}
		current = next
		path = append(path, next.Name)
		idx++
	}

	rest := args[idx:]
	if len(rest) > 0 && isHelpToken(rest[0]) {
		WriteHelp(ctx, current, path)
		return nil
	}
	if current.Run == nil {
		if len(rest) > 0 {
			return OperatorError(fmt.Sprintf("unknown command %q", strings.Join(args, " ")))
		}
		WriteHelp(ctx, current, path)
		return nil
	}
	return current.Run(ctx, rest)
}

func dispatchHelp(root *Command, ctx *Context, args []string) error {
	if len(args) == 0 {
		WriteHelp(ctx, root, nil)
		return nil
	}
	current := root
	path := []string{}
	for _, tok := range args {
		next := findSubcommand(current, tok)
		if next == nil {
			return OperatorError(fmt.Sprintf("unknown command path for help: %q", strings.Join(args, " ")))
		}
		current = next
		path = append(path, next.Name)
	}
	WriteHelp(ctx, current, path)
	return nil
}

func WriteHelp(ctx *Context, cmd *Command, path []string) {
	if cmd == nil {
		return
	}
	fullPath := "sros"
	if len(path) > 0 {
		fullPath += " " + strings.Join(path, " ")
	}

	lines := []string{fmt.Sprintf("%s - %s", fullPath, cmd.Summary), "", "Usage:"}
	usage := strings.TrimSpace(cmd.Usage)
	if usage == "" {
		usage = fullPath + " [flags]"
	}
	lines = append(lines, "  "+usage)

	if len(cmd.Subcommands) > 0 {
		lines = append(lines, "", "Commands:")
		for _, sub := range cmd.Subcommands {
			lines = append(lines, fmt.Sprintf("  %-12s %s", sub.Name, sub.Summary))
		}
	}
	if len(cmd.Examples) > 0 {
		lines = append(lines, "", "Examples:")
		for _, ex := range cmd.Examples {
			lines = append(lines, "  "+ex)
		}
	}

	if cmd.Name == "sros" {
		lines = append(lines, "", "Global Flags:")
		lines = append(lines,
			"  --config <path>      Use explicit config file",
			"  --format <text|json> Set output format",
			"  --workspace <path>   Override workspace root",
			"  -h, --help           Show help",
		)
	}
	lines = append(lines, "", "Run 'sros help <command>' for more details.")
	_, _ = fmt.Fprintln(ctx.Stdout, strings.Join(lines, "\n"))
}

func findSubcommand(cmd *Command, token string) *Command {
	for _, sub := range cmd.Subcommands {
		if token == sub.Name {
			return sub
		}
		for _, alias := range sub.Aliases {
			if token == alias {
				return sub
			}
		}
	}
	return nil
}

func isHelpToken(token string) bool { return token == "-h" || token == "--help" }

func writeOutput(ctx *Context, text string, payload any) error {
	if strings.ToLower(ctx.OutputFormat) == "json" {
		encoded, err := json.MarshalIndent(payload, "", "  ")
		if err != nil {
			return InternalError("encode output")
		}
		_, _ = fmt.Fprintln(ctx.Stdout, string(encoded))
		return nil
	}
	_, _ = fmt.Fprintln(ctx.Stdout, text)
	return nil
}

func formatBoundaries(boundaries []boot.ServiceBoundary) string {
	if len(boundaries) == 0 {
		return "(none)"
	}
	rows := make([]string, 0, len(boundaries))
	for _, b := range boundaries {
		state := "deferred"
		if b.Wired {
			state = "wired"
		}
		if b.DeferredTo != "" {
			rows = append(rows, fmt.Sprintf("- %s: %s (%s)", b.Name, state, b.DeferredTo))
		} else {
			rows = append(rows, fmt.Sprintf("- %s: %s", b.Name, state))
		}
	}
	sort.Strings(rows)
	return strings.Join(rows, "\n")
}

func CommandPaths(root *Command) []string {
	paths := []string{}
	var walk func(prefix string, cmd *Command)
	walk = func(prefix string, cmd *Command) {
		path := prefix
		if prefix == "" {
			path = cmd.Name
		} else if cmd.Name != "sros" {
			path = prefix + " " + cmd.Name
		}
		if cmd.Name != "sros" {
			paths = append(paths, path)
		}
		for _, sub := range cmd.Subcommands {
			walk(path, sub)
		}
	}
	walk("", root)
	sort.Strings(paths)
	return paths
}

func requireNoArgs(args []string) error {
	if len(args) != 0 {
		return OperatorError("this command does not accept positional arguments")
	}
	return nil
}

type ioDiscard struct{}

func (ioDiscard) Write(p []byte) (int, error) { return len(p), nil }

func stringsJoin(lines []string) string {
	if len(lines) == 0 {
		return ""
	}
	return strings.Join(lines, "\n")
}
