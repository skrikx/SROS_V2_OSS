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
	_, releaseHelp, _ := runCLI(t, []string{"release", "--help"}, nil)
	_, testHelp, _ := runCLI(t, []string{"test", "--help"}, nil)
	_, examplesHelp, _ := runCLI(t, []string{"examples", "--help"}, nil)

	combined := strings.Join([]string{rootHelp, toolsHelp, connectorsHelp, mcpHelp, releaseHelp, testHelp, examplesHelp}, "\n")
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
		"replay",
		"verify",
		"release",
		"pack",
		"test",
		"smoke",
		"examples",
		"scaffold",
	}
	for _, req := range required {
		if !strings.Contains(combined, req) {
			t.Fatalf("missing command token in help output: %s", req)
		}
	}
}
