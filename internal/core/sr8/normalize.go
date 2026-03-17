package sr8

import (
	"crypto/sha256"
	"encoding/hex"
	"strings"
	"time"
)

func Normalize(parsed ParsedIntent, now time.Time) NormalizedIntent {
	text := normalizeWhitespace(parsed.RawIntent)
	summary := text
	if len(summary) > 120 {
		summary = summary[:120]
	}

	seed := text + "|" + parsed.OperatorID + "|" + now.UTC().Format(time.RFC3339Nano)
	requestID := "cmp_" + shortHash(seed)

	return NormalizedIntent{
		Source:           parsed.Source,
		RawIntent:        parsed.RawIntent,
		NormalizedText:   text,
		IntentSummary:    summary,
		InputPath:        parsed.InputPath,
		OperatorID:       parsed.OperatorID,
		TenantID:         parsed.TenantID,
		WorkspaceID:      parsed.WorkspaceID,
		CompileRequestID: requestID,
		RequestedAt:      now.UTC(),
	}
}

func normalizeWhitespace(v string) string {
	v = strings.ReplaceAll(v, "\r\n", "\n")
	lines := strings.Split(v, "\n")
	for i := range lines {
		lines[i] = strings.TrimSpace(lines[i])
	}
	v = strings.Join(lines, " ")
	return strings.Join(strings.Fields(v), " ")
}

func shortHash(v string) string {
	h := sha256.Sum256([]byte(v))
	return hex.EncodeToString(h[:])[:12]
}
