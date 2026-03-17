package cli_test

import (
	"strings"
	"testing"
)

func TestCommandTreeContainsRequiredCommands(t *testing.T) {
	_, rootHelp, rootErr := runCLI(t, []string{"--help"}, nil)
	if rootErr != "" {
		// help should not emit stderr
		t.Fatalf("unexpected stderr for root help: %s", rootErr)
	}
	_, toolsHelp, _ := runCLI(t, []string{"tools", "--help"}, nil)
	_, connectorsHelp, _ := runCLI(t, []string{"connectors", "--help"}, nil)
	_, mcpHelp, _ := runCLI(t, []string{"mcp", "--help"}, nil)

	combined := strings.Join([]string{rootHelp, toolsHelp, connectorsHelp, mcpHelp}, "\n")
	required := []string{
		"init",
		"bootstrap",
		"doctor",
		"seed",
		"config",
		"compile",
		"run",
		"plan",
		"resume",
		"pause",
		"checkpoint",
		"rollback",
		"trace",
		"receipts",
		"memory",
		"mirror",
		"inspect",
		"status",
		"tools",
		"list",
		"show",
		"validate",
		"register",
		"connectors",
		"mcp",
		"ingest",
	}
	for _, req := range required {
		if !strings.Contains(combined, req) {
			t.Fatalf("missing command token in help output: %s", req)
		}
	}
}
