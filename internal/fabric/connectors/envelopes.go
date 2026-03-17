package connectors

import (
	"encoding/json"
	"fmt"
	"os"
	"time"

	ctools "srosv2/contracts/tools"
)

func LoadEnvelope(path string) (ctools.AuthEnvelope, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ctools.AuthEnvelope{}, fmt.Errorf("read envelope %s: %w", path, err)
	}
	var env ctools.AuthEnvelope
	if err := json.Unmarshal(data, &env); err != nil {
		return ctools.AuthEnvelope{}, fmt.Errorf("decode envelope %s: %w", path, err)
	}
	if errs := ctools.ValidateEnvelope(env); len(errs) > 0 {
		return ctools.AuthEnvelope{}, errs[0]
	}
	if env.ExpiresAt.Before(time.Now().UTC()) {
		return ctools.AuthEnvelope{}, fmt.Errorf("envelope %s is expired", env.EnvelopeID)
	}
	return env, nil
}

func RedactEnvelope(env ctools.AuthEnvelope) map[string]any {
	return map[string]any{
		"envelope_version": env.EnvelopeVersion,
		"envelope_id":      env.EnvelopeID,
		"connector":        env.Connector,
		"auth_type":        env.AuthType,
		"class":            env.Class,
		"display_name":     env.DisplayName,
		"metadata":         env.Metadata,
		"expires_at":       env.ExpiresAt,
		"trace_link":       env.TraceLink,
		"secret_material":  "[REDACTED]",
	}
}
