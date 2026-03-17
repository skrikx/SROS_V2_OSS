package connectors

import (
	"context"

	ctools "srosv2/contracts/tools"
)

func Invoke(ctx context.Context, registry *Registry, name string, req map[string]any, env ctools.AuthEnvelope) (map[string]any, error) {
	adapter, err := registry.Adapter(name)
	if err != nil {
		return nil, err
	}
	return adapter.Invoke(ctx, req, env)
}
