package adapters

import (
	"context"
	"os"

	ctools "srosv2/contracts/tools"
)

type LocalFS struct{}

func (LocalFS) Name() string { return "local_fs" }

func (LocalFS) Invoke(_ context.Context, req map[string]any, env ctools.AuthEnvelope) (map[string]any, error) {
	path, _ := req["path"].(string)
	info, err := os.Stat(path)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"connector": env.Connector,
		"path":      path,
		"name":      info.Name(),
		"dir":       info.IsDir(),
	}, nil
}
