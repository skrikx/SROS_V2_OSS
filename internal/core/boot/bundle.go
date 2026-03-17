package boot

import (
	"srosv2/internal/core/gov"
	"srosv2/internal/core/mem"
	"srosv2/internal/core/mirror"
	"srosv2/internal/core/orch"
	"srosv2/internal/core/runtime"
)

type Bundle struct {
	Mode         Mode               `json:"mode"`
	Compiler     Compiler           `json:"-"`
	Runtime      runtime.Runtime    `json:"-"`
	Inspector    runtime.Inspector  `json:"-"`
	Fabric       runtime.Fabric     `json:"-"`
	Orchestrator *orch.Orchestrator `json:"-"`
	Governor     *gov.Engine        `json:"-"`
	Memory       *mem.Store         `json:"-"`
	Mirror       *mirror.Engine     `json:"-"`
	Boundaries   []ServiceBoundary  `json:"boundaries"`
}

func (b Bundle) Boundary(name string) (ServiceBoundary, bool) {
	for _, item := range b.Boundaries {
		if item.Name == name {
			return item, true
		}
	}
	return ServiceBoundary{}, false
}
