package policy

import "srosv2/internal/shared/ids"

type Bundle struct {
	BundleID              ids.PolicyBundleID        `json:"bundle_id"`
	Name                  string                    `json:"name"`
	Version               string                    `json:"version"`
	RulesetDigest         string                    `json:"ruleset_digest"`
	DefaultVerdict        Verdict                   `json:"default_verdict,omitempty"`
	DefaultSandboxProfile string                    `json:"default_sandbox_profile,omitempty"`
	BreakGlassAllowed     bool                      `json:"break_glass_allowed,omitempty"`
	Capabilities          []CapabilityPolicy        `json:"capabilities,omitempty"`
	Sandboxes             map[string]SandboxProfile `json:"sandboxes,omitempty"`
	Metadata              map[string]string         `json:"metadata,omitempty"`
}

type CapabilityPolicy struct {
	Name              string          `json:"name"`
	Verdict           Verdict         `json:"verdict"`
	SandboxProfile    string          `json:"sandbox_profile,omitempty"`
	AllowedBoundaries []TrustBoundary `json:"allowed_boundaries,omitempty"`
	MaxRiskClass      string          `json:"max_risk_class,omitempty"`
	RequireApproval   bool            `json:"require_approval,omitempty"`
	AllowBreakGlass   bool            `json:"allow_break_glass,omitempty"`
}

type SandboxProfile struct {
	Name             string `json:"name"`
	AllowShell       bool   `json:"allow_shell,omitempty"`
	AllowPatch       bool   `json:"allow_patch,omitempty"`
	AllowExternalNet bool   `json:"allow_external_net,omitempty"`
	Description      string `json:"description,omitempty"`
}
