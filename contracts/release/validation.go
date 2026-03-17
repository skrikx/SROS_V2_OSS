package release

import "srosv2/internal/shared/validation"

func ValidateCheckpoint(record CheckpointRecord) []error {
	var errs []error
	appendErr := func(err error) {
		if err != nil {
			errs = append(errs, err)
		}
	}

	appendErr(validation.RequiredString("contract_version", record.ContractVersion))
	appendErr(validation.RequiredString("checkpoint_id", string(record.CheckpointID)))
	appendErr(validation.RequiredString("run_id", string(record.RunID)))
	appendErr(validation.Enum("stage", string(record.Stage), []string{"draft", "validated", "promoted", "rolled_back"}))
	appendErr(validation.RequiredTime("recorded_at", record.RecordedAt))
	return errs
}

func ValidateRelease(record ReleaseRecord) []error {
	var errs []error
	appendErr := func(err error) {
		if err != nil {
			errs = append(errs, err)
		}
	}

	appendErr(validation.RequiredString("contract_version", record.ContractVersion))
	appendErr(validation.RequiredString("release_id", string(record.ReleaseID)))
	appendErr(validation.RequiredString("checkpoint_id", string(record.CheckpointID)))
	appendErr(validation.Enum("target_stage", string(record.TargetStage), []string{"draft", "validated", "promoted", "rolled_back"}))
	appendErr(validation.RequiredTime("created_at", record.CreatedAt))
	return errs
}

func ValidateRollback(record RollbackRecord) []error {
	var errs []error
	appendErr := func(err error) {
		if err != nil {
			errs = append(errs, err)
		}
	}

	appendErr(validation.RequiredString("contract_version", record.ContractVersion))
	appendErr(validation.RequiredString("rollback_id", string(record.RollbackID)))
	appendErr(validation.RequiredString("release_id", string(record.ReleaseID)))
	appendErr(validation.RequiredString("target_checkpoint_id", string(record.TargetCheckpointID)))
	appendErr(validation.RequiredString("reason", record.Reason))
	appendErr(validation.RequiredTime("created_at", record.CreatedAt))
	return errs
}
