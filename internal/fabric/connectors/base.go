package connectors

import (
	"context"

	ctools "srosv2/contracts/tools"
)

type Adapter interface {
	Name() string
	Invoke(context.Context, map[string]any, ctools.AuthEnvelope) (map[string]any, error)
}
