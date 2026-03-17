package evidence

import "srosv2/internal/shared/ids"

type DigestAlgorithm string

const (
	DigestAlgorithmSHA256 DigestAlgorithm = "sha256"
	DigestAlgorithmBLAKE3 DigestAlgorithm = "blake3"
)

type ArtifactRef struct {
	ArtifactID ids.ArtifactID  `json:"artifact_id"`
	Path       string          `json:"path"`
	DigestAlgo DigestAlgorithm `json:"digest_algo"`
	Digest     string          `json:"digest"`
	MediaType  string          `json:"media_type,omitempty"`
	Bytes      int64           `json:"bytes,omitempty"`
}
