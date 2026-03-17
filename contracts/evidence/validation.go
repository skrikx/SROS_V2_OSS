package evidence

import "srosv2/internal/shared/validation"

func ValidateReceipt(receipt Receipt) []error {
	var errs []error
	appendErr := func(err error) {
		if err != nil {
			errs = append(errs, err)
		}
	}

	appendErr(validation.RequiredString("contract_version", receipt.ContractVersion))
	appendErr(validation.RequiredString("receipt_id", string(receipt.ReceiptID)))
	appendErr(validation.RequiredString("run_id", string(receipt.RunID)))
	appendErr(validation.Enum("kind", string(receipt.Kind), []string{"terminal", "stage", "policy", "closure"}))
	appendErr(validation.RequiredString("evidence_bundle_id", string(receipt.EvidenceBundleID)))
	appendErr(validation.RequiredString("status", receipt.Status))
	appendErr(validation.RequiredString("summary", receipt.Summary))
	appendErr(validation.RequiredTime("created_at", receipt.CreatedAt))
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
	appendErr(validation.RequiredString("run_id", string(bundle.RunID)))
	appendErr(validation.RequiredSlice("artifact_refs", bundle.ArtifactRefs))
	return errs
}

func ValidateArtifactRef(ref ArtifactRef) []error {
	var errs []error
	appendErr := func(err error) {
		if err != nil {
			errs = append(errs, err)
		}
	}

	appendErr(validation.RequiredString("artifact_id", string(ref.ArtifactID)))
	appendErr(validation.RequiredString("path", ref.Path))
	appendErr(validation.Enum("digest_algo", string(ref.DigestAlgo), []string{"sha256", "blake3"}))
	appendErr(validation.RequiredString("digest", ref.Digest))
	return errs
}
