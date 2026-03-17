package sr8

import (
	"strings"

	"srosv2/internal/ace/intake"
)

func Classify(intent NormalizedIntent) Classification {
	plan := intake.BuildPlan(intent.NormalizedText)

	return Classification{
		Domain:          DomainClass(plan.Classification.Domain),
		Risk:            plan.Classification.Risk,
		ArtifactPosture: plan.ArtifactPosture,
		Signals:         append([]string{}, plan.Classification.Signals...),
	}
}

func likelyArtifacts(text string) string {
	t := strings.ToLower(text)
	switch {
	case strings.Contains(t, "patch"), strings.Contains(t, "file"), strings.Contains(t, "edit"):
		return "file_delta"
	case strings.Contains(t, "research"), strings.Contains(t, "investigate"):
		return "research_brief"
	default:
		return "run_contract_only"
	}
}
