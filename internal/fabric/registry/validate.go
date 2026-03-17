package registry

import (
	"strings"

	ctools "srosv2/contracts/tools"
	"srosv2/internal/fabric/harness"
)

func ValidateManifest(manifest ctools.Manifest, hr *harness.Harness) ctools.ValidationResult {
	result := ctools.ValidationResult{Manifest: manifest}
	for _, err := range ctools.ValidateManifest(manifest) {
		result.Errors = append(result.Errors, err.Error())
	}
	result.PolicyBindingPresent = strings.TrimSpace(manifest.PolicyClass) != ""
	if hr != nil {
		resolution, err := hr.Resolve(manifest)
		result.HarnessCompatible = err == nil && resolution.Compatible
		if err != nil {
			result.Errors = append(result.Errors, err.Error())
		} else if !resolution.Compatible {
			result.Errors = append(result.Errors, resolution.Reason)
		}
	}
	result.Valid = len(result.Errors) == 0
	return result
}
