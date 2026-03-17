package mirror

import (
	"crypto/sha256"
	"encoding/hex"
)

func shortHash(v string) string {
	h := sha256.Sum256([]byte(v))
	return hex.EncodeToString(h[:])[:12]
}
