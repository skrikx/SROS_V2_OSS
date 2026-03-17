package mcpclient

import (
	"encoding/json"
	"fmt"
	"os"

	ctools "srosv2/contracts/tools"
	"srosv2/internal/fabric/registry"
)

func Ingest(path string, reg *registry.Registry) (ctools.Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ctools.Manifest{}, fmt.Errorf("read MCP input %s: %w", path, err)
	}
	var remote RemoteCapability
	if err := json.Unmarshal(data, &remote); err != nil {
		return ctools.Manifest{}, fmt.Errorf("decode MCP input %s: %w", path, err)
	}
	manifest := Normalize(remote)
	manifest, _, err = reg.Register(manifest)
	if err != nil {
		return ctools.Manifest{}, err
	}
	return reg.Admit(manifest.Name)
}
