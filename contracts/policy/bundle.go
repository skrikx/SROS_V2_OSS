package policy

import "srosv2/internal/shared/ids"

type Bundle struct {
	BundleID      ids.PolicyBundleID `json:"bundle_id"`
	Name          string             `json:"name"`
	Version       string             `json:"version"`
	RulesetDigest string             `json:"ruleset_digest"`
}
