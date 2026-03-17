package memory

import "srosv2/internal/shared/validation"

func ValidateMutation(mutation MemoryMutation) []error {
	var errs []error
	appendErr := func(err error) {
		if err != nil {
			errs = append(errs, err)
		}
	}

	appendErr(validation.RequiredString("contract_version", mutation.ContractVersion))
	appendErr(validation.RequiredString("mutation_id", string(mutation.MutationID)))
	appendErr(validation.RequiredString("run_id", string(mutation.RunID)))
	appendErr(validation.RequiredString("session_id", string(mutation.SessionID)))
	appendErr(validation.Enum("scope", string(mutation.Scope), []string{"session", "workspace", "run", "global"}))
	appendErr(validation.Enum("kind", string(mutation.Kind), []string{"upsert", "delete", "link", "annotate", "prune_recommend", "compact_recommend"}))
	appendErr(validation.RequiredString("branch.branch_id", string(mutation.Branch.BranchID)))
	appendErr(validation.RequiredString("key", mutation.Key))
	appendErr(validation.RequiredTime("occurred_at", mutation.OccurredAt))
	return errs
}
