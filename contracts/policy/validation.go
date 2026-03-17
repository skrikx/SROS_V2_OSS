package policy

import "srosv2/internal/shared/validation"

func ValidateDecision(decision PolicyDecision) []error {
	var errs []error
	appendErr := func(err error) {
		if err != nil {
			errs = append(errs, err)
		}
	}

	appendErr(validation.RequiredString("contract_version", decision.ContractVersion))
	appendErr(validation.RequiredString("decision_id", string(decision.DecisionID)))
	appendErr(validation.RequiredString("run_id", string(decision.RunID)))
	appendErr(validation.RequiredString("trace_id", string(decision.TraceID)))
	appendErr(validation.Enum("verdict", string(decision.Verdict), []string{"allow", "ask", "deny"}))
	appendErr(validation.Enum("boundary", string(decision.Boundary), []string{"local_process", "local_fs", "external_net"}))
	appendErr(validation.RequiredString("sandbox_profile", decision.SandboxProfile))
	appendErr(validation.RequiredString("bundle_ref", string(decision.BundleRef)))
	appendErr(validation.RequiredString("reason", decision.Reason))
	appendErr(validation.RequiredTime("decided_at", decision.DecidedAt))
	return errs
}

func ValidateBundle(bundle Bundle) []error {
	var errs []error
	appendErr := func(err error) {
		if err != nil {
			errs = append(errs, err)
		}
	}

	appendErr(validation.RequiredString("bundle_id", string(bundle.BundleID)))
	appendErr(validation.RequiredString("name", bundle.Name))
	appendErr(validation.RequiredString("version", bundle.Version))
	appendErr(validation.RequiredString("ruleset_digest", bundle.RulesetDigest))
	return errs
}
