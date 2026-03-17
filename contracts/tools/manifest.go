package tools

import (
	"strings"

	"srosv2/internal/shared/validation"
)

type LifecycleState string

const (
	StateDraft        LifecycleState = "draft"
	StateValidated    LifecycleState = "validated"
	StateAdmitted     LifecycleState = "admitted"
	StateExperimental LifecycleState = "experimental"
	StateActive       LifecycleState = "active"
	StateQuarantined  LifecycleState = "quarantined"
	StateDeprecated   LifecycleState = "deprecated"
	StateDisabled     LifecycleState = "disabled"
)

type Manifest struct {
	ManifestVersion      string            `json:"manifest_version"`
	Name                 string            `json:"name"`
	Title                string            `json:"title"`
	Description          string            `json:"description"`
	Version              string            `json:"version"`
	Class                string            `json:"class"`
	Domain               string            `json:"domain"`
	PolicyClass          string            `json:"policy_class"`
	Status               LifecycleState    `json:"status"`
	TrustBoundary        string            `json:"trust_boundary"`
	SandboxProfiles      []string          `json:"sandbox_profiles"`
	AuthType             string            `json:"auth_type"`
	AuthEnvelopeRequired bool              `json:"auth_envelope_required"`
	MCPIngressCapable    bool              `json:"mcp_ingress_capable"`
	MCPEgressCapable     bool              `json:"mcp_egress_capable"`
	RemoteCapable        bool              `json:"remote_capable"`
	Unsafe               bool              `json:"unsafe"`
	Experimental         bool              `json:"experimental"`
	QuarantineReason     string            `json:"quarantine_reason,omitempty"`
	ConnectorRef         string            `json:"connector_ref,omitempty"`
	Command              []string          `json:"command,omitempty"`
	Args                 []string          `json:"args,omitempty"`
	AllowedPaths         []string          `json:"allowed_paths,omitempty"`
	Capabilities         []string          `json:"capabilities,omitempty"`
	Metadata             map[string]string `json:"metadata,omitempty"`
}

type ManifestSummary struct {
	Name                 string         `json:"name"`
	Title                string         `json:"title"`
	Class                string         `json:"class"`
	Domain               string         `json:"domain"`
	PolicyClass          string         `json:"policy_class"`
	Status               LifecycleState `json:"status"`
	TrustBoundary        string         `json:"trust_boundary"`
	AuthType             string         `json:"auth_type"`
	AuthEnvelopeRequired bool           `json:"auth_envelope_required"`
	Unsafe               bool           `json:"unsafe"`
}

func (m Manifest) Summary() ManifestSummary {
	return ManifestSummary{
		Name:                 m.Name,
		Title:                m.Title,
		Class:                m.Class,
		Domain:               m.Domain,
		PolicyClass:          m.PolicyClass,
		Status:               m.Status,
		TrustBoundary:        m.TrustBoundary,
		AuthType:             m.AuthType,
		AuthEnvelopeRequired: m.AuthEnvelopeRequired,
		Unsafe:               m.Unsafe,
	}
}

func ValidateManifest(manifest Manifest) []error {
	var errs []error
	appendErr := func(err error) {
		if err != nil {
			errs = append(errs, err)
		}
	}

	appendErr(validation.RequiredString("manifest_version", manifest.ManifestVersion))
	appendErr(validation.RequiredString("name", manifest.Name))
	appendErr(validation.RequiredString("title", manifest.Title))
	appendErr(validation.RequiredString("description", manifest.Description))
	appendErr(validation.RequiredString("version", manifest.Version))
	appendErr(validation.RequiredString("class", manifest.Class))
	appendErr(validation.RequiredString("domain", manifest.Domain))
	appendErr(validation.RequiredString("policy_class", manifest.PolicyClass))
	appendErr(validation.RequiredString("trust_boundary", manifest.TrustBoundary))
	appendErr(validation.RequiredSlice("sandbox_profiles", manifest.SandboxProfiles))
	appendErr(validation.Enum("status", string(manifest.Status), []string{
		string(StateDraft),
		string(StateValidated),
		string(StateAdmitted),
		string(StateExperimental),
		string(StateActive),
		string(StateQuarantined),
		string(StateDeprecated),
		string(StateDisabled),
	}))
	if manifest.AuthEnvelopeRequired && strings.TrimSpace(manifest.AuthType) == "" {
		appendErr(validation.RequiredString("auth_type", manifest.AuthType))
	}
	return errs
}
