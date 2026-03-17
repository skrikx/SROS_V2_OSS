package runcontract

import (
	"strconv"

	"srosv2/internal/shared/validation"
)

func Validate(contract RunContract) []error {
	var errs []error

	appendErr := func(err error) {
		if err != nil {
			errs = append(errs, err)
		}
	}

	appendErr(validation.RequiredString("contract_version", contract.ContractVersion))
	appendErr(validation.RequiredString("run_id", string(contract.RunID)))
	appendErr(validation.RequiredString("trace_id", string(contract.TraceID)))
	appendErr(validation.RequiredString("operator_id", string(contract.OperatorID)))
	appendErr(validation.RequiredString("tenant_id", string(contract.TenantID)))
	appendErr(validation.RequiredString("workspace_id", string(contract.WorkspaceID)))
	appendErr(validation.RequiredString("intent_summary", contract.IntentSummary))
	appendErr(validation.RequiredString("normalized_request", contract.NormalizedRequest))
	appendErr(validation.Enum("risk_class", string(contract.RiskClass), []string{"low", "medium", "high", "critical"}))
	appendErr(validation.Enum("route_class", string(contract.RouteClass), []string{"local_cli", "local_run"}))
	appendErr(validation.RequiredTime("created_at", contract.CreatedAt))

	for i, req := range contract.RequestedReceipts {
		index := strconv.Itoa(i)
		appendErr(validation.Enum("requested_receipts["+index+"].mode", string(req.Mode), []string{"none", "summary", "full"}))
	}

	for i, ref := range contract.CheckpointRefs {
		index := strconv.Itoa(i)
		appendErr(validation.RequiredString("checkpoint_refs["+index+"].checkpoint_id", string(ref.CheckpointID)))
		appendErr(validation.RequiredString("checkpoint_refs["+index+"].stage", ref.Stage))
	}

	return errs
}
