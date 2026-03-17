package tools

import (
	"time"

	"srosv2/internal/shared/validation"
)

type AuthEnvelope struct {
	EnvelopeVersion string            `json:"envelope_version"`
	EnvelopeID      string            `json:"envelope_id"`
	Connector       string            `json:"connector"`
	AuthType        string            `json:"auth_type"`
	Class           string            `json:"class"`
	DisplayName     string            `json:"display_name"`
	Metadata        map[string]string `json:"metadata,omitempty"`
	SecretMaterial  string            `json:"secret_material,omitempty"`
	ExpiresAt       time.Time         `json:"expires_at"`
	TraceLink       string            `json:"trace_link,omitempty"`
}

func ValidateEnvelope(env AuthEnvelope) []error {
	var errs []error
	appendErr := func(err error) {
		if err != nil {
			errs = append(errs, err)
		}
	}

	appendErr(validation.RequiredString("envelope_version", env.EnvelopeVersion))
	appendErr(validation.RequiredString("envelope_id", env.EnvelopeID))
	appendErr(validation.RequiredString("connector", env.Connector))
	appendErr(validation.RequiredString("auth_type", env.AuthType))
	appendErr(validation.RequiredString("class", env.Class))
	appendErr(validation.RequiredString("display_name", env.DisplayName))
	appendErr(validation.RequiredTime("expires_at", env.ExpiresAt))
	return errs
}
