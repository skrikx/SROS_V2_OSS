package provenance

import (
	"encoding/json"

	"srosv2/contracts/evidence"
)

func HashArtifact(ref evidence.ArtifactRef) (string, error) {
	data, err := json.Marshal(ref)
	if err != nil {
		return "", err
	}
	return digestBytes(data), nil
}

func HashChain(refs []evidence.ArtifactRef) ([]string, error) {
	chain := make([]string, 0, len(refs))
	previous := ""
	for _, ref := range refs {
		hash, err := HashArtifact(ref)
		if err != nil {
			return nil, err
		}
		previous = digestBytes([]byte(previous + hash))
		chain = append(chain, previous)
	}
	return chain, nil
}
