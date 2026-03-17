package intake

import ctools "srosv2/contracts/tools"

type Shortlist struct {
	Skills       []SkillHint          `json:"skills"`
	PromptUnits  []PromptUnitHint     `json:"prompt_units"`
	Capabilities []ctools.SearchMatch `json:"capabilities,omitempty"`
}

func BuildShortlist(class Classification) Shortlist {
	return Shortlist{
		Skills:      SuggestSkills(class.Domain),
		PromptUnits: SuggestPromptUnits(class.Domain, class.Risk),
	}
}

func BuildCapabilityShortlist(class Classification, matches []ctools.SearchMatch) Shortlist {
	shortlist := BuildShortlist(class)
	shortlist.Capabilities = matches
	return shortlist
}
