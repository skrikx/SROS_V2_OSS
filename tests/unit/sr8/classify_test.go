package sr8_test

import (
	"testing"

	"srosv2/contracts/runcontract"
	"srosv2/internal/core/sr8"
)

func TestClassifyFileTask(t *testing.T) {
	class := sr8.Classify(sr8.NormalizedIntent{NormalizedText: "patch file and update config"})
	if class.Domain != sr8.DomainFileTask {
		t.Fatalf("expected file task domain, got %s", class.Domain)
	}
}

func TestClassifyRisk(t *testing.T) {
	class := sr8.Classify(sr8.NormalizedIntent{NormalizedText: "delete production table"})
	if class.Risk != runcontract.RiskClassHigh && class.Risk != runcontract.RiskClassCritical {
		t.Fatalf("expected high/critical risk, got %s", class.Risk)
	}
}
