package tools

type NegotiationRequest struct {
	Query           SearchQuery `json:"query"`
	PreferredTool   string      `json:"preferred_tool,omitempty"`
	RequireRunnable bool        `json:"require_runnable,omitempty"`
}

type NegotiationResult struct {
	Allowed            bool          `json:"allowed"`
	Selected           *SearchMatch  `json:"selected,omitempty"`
	Alternates         []SearchMatch `json:"alternates,omitempty"`
	Reason             string        `json:"reason"`
	PolicyFit          string        `json:"policy_fit"`
	EnvironmentFit     string        `json:"environment_fit"`
	SelectedCapability string        `json:"selected_capability,omitempty"`
}
