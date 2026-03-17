package tools

type SearchQuery struct {
	Class             string           `json:"class,omitempty"`
	Domain            string           `json:"domain,omitempty"`
	PolicyClass       string           `json:"policy_class,omitempty"`
	SandboxProfile    string           `json:"sandbox_profile,omitempty"`
	MCPMode           string           `json:"mcp_mode,omitempty"`
	AuthType          string           `json:"auth_type,omitempty"`
	Status            []LifecycleState `json:"status,omitempty"`
	TrustBoundary     string           `json:"trust_boundary,omitempty"`
	IncludeUnsafe     bool             `json:"include_unsafe,omitempty"`
	IncludeHistorical bool             `json:"include_historical,omitempty"`
	Limit             int              `json:"limit,omitempty"`
}

type SearchMatch struct {
	Rank                 int             `json:"rank"`
	Manifest             ManifestSummary `json:"manifest"`
	PolicyFit            string          `json:"policy_fit"`
	EnvironmentFit       string          `json:"environment_fit"`
	AuthEnvelopeRequired bool            `json:"auth_envelope_required"`
	SelectedSandbox      string          `json:"selected_sandbox,omitempty"`
	Unsafe               bool            `json:"unsafe"`
	Selectable           bool            `json:"selectable"`
	Reason               string          `json:"reason"`
}

type ValidationResult struct {
	Valid                bool     `json:"valid"`
	Manifest             Manifest `json:"manifest"`
	HarnessCompatible    bool     `json:"harness_compatible"`
	PolicyBindingPresent bool     `json:"policy_binding_present"`
	Errors               []string `json:"errors,omitempty"`
}

type InvocationResult struct {
	Allowed        bool           `json:"allowed"`
	Capability     string         `json:"capability"`
	Status         string         `json:"status"`
	SandboxProfile string         `json:"sandbox_profile,omitempty"`
	TraceLinked    bool           `json:"trace_linked"`
	ReceiptLinked  bool           `json:"receipt_linked"`
	Output         map[string]any `json:"output,omitempty"`
	Reason         string         `json:"reason,omitempty"`
}
