package intake

type SkillHint struct {
	Name   string `json:"name"`
	Reason string `json:"reason"`
}

func SuggestSkills(domain DomainClass) []SkillHint {
	switch domain {
	case DomainFileTask:
		return []SkillHint{{Name: "patch-planning", Reason: "file-oriented intent detected"}}
	case DomainResearch:
		return []SkillHint{{Name: "research-brief", Reason: "research intent detected"}}
	default:
		return []SkillHint{{Name: "operator-baseline", Reason: "general local operation"}}
	}
}
