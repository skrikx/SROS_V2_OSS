package gov

import (
	"fmt"
	"strings"
	"time"

	"srosv2/contracts/policy"
	"srosv2/contracts/runcontract"
	"srosv2/internal/shared/ids"
)

type Request struct {
	RunID       ids.RunID
	TraceID     ids.TraceID
	RiskClass   runcontract.RiskClass
	Capability  string
	BreakGlass  bool
	Description string
}

type Result struct {
	Decision policy.PolicyDecision `json:"decision"`
}

func newDecision(bundle policy.Bundle, req Request, verdict policy.Verdict, boundary policy.TrustBoundary, sandbox, reason string, now time.Time) policy.PolicyDecision {
	return policy.PolicyDecision{
		ContractVersion: "v2.0",
		DecisionID:      ids.PolicyDecisionID("pd_" + decisionSuffix(req, now)),
		RunID:           req.RunID,
		TraceID:         req.TraceID,
		Capability:      req.Capability,
		Verdict:         verdict,
		Boundary:        boundary,
		SandboxProfile:  sandbox,
		BundleRef:       bundle.BundleID,
		Reason:          strings.TrimSpace(reason),
		BreakGlass:      req.BreakGlass,
		DecidedAt:       now.UTC(),
	}
}

func decisionSuffix(req Request, now time.Time) string {
	return fmt.Sprintf("%s_%d", strings.ReplaceAll(req.Capability, ".", "_"), now.UTC().UnixNano())
}
