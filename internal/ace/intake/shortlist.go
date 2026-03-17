package intake

type Shortlist struct {
	Skills      []SkillHint      `json:"skills"`
	PromptUnits []PromptUnitHint `json:"prompt_units"`
}

func BuildShortlist(class Classification) Shortlist {
	return Shortlist{
		Skills:      SuggestSkills(class.Domain),
		PromptUnits: SuggestPromptUnits(class.Domain, class.Risk),
	}
}
