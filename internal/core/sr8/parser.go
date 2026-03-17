package sr8

import (
	"fmt"
	"os"
	"strings"

	"srosv2/internal/shared/ids"
)

func ParseRequest(req CompileRequest) (ParsedIntent, error) {
	inline := strings.TrimSpace(req.Intent)
	input := strings.TrimSpace(req.InputPath)

	if inline == "" && input == "" {
		return ParsedIntent{}, fmt.Errorf("compile requires intent text or input path")
	}
	if inline != "" && input != "" {
		return ParsedIntent{}, fmt.Errorf("provide either --intent or --input, not both")
	}

	out := ParsedIntent{
		OperatorID:  normalizeOperator(req.OperatorID),
		TenantID:    normalizeTenant(req.TenantID),
		WorkspaceID: normalizeWorkspace(req.WorkspaceID),
	}

	if inline != "" {
		out.Source = IntentSourceInline
		out.RawIntent = inline
		return out, nil
	}

	data, err := os.ReadFile(input)
	if err != nil {
		return ParsedIntent{}, fmt.Errorf("read compile input %q: %w", input, err)
	}
	text := strings.TrimSpace(string(data))
	if text == "" {
		return ParsedIntent{}, fmt.Errorf("compile input %q is empty", input)
	}

	out.Source = IntentSourceFile
	out.InputPath = input
	out.RawIntent = text
	return out, nil
}

func normalizeOperator(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return "operator_local"
	}
	if err := ids.Validate(v); err != nil {
		return "operator_local"
	}
	return v
}

func normalizeTenant(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return string(ids.DefaultTenantID)
	}
	return v
}

func normalizeWorkspace(v string) string {
	v = strings.TrimSpace(v)
	if v == "" {
		return string(ids.DefaultWorkspaceID)
	}
	return v
}
