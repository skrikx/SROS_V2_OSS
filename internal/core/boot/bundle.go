package boot

import "srosv2/internal/core/runtime"

type Bundle struct {
	Mode       Mode              `json:"mode"`
	Compiler   runtime.Compiler  `json:"-"`
	Runtime    runtime.Runtime   `json:"-"`
	Inspector  runtime.Inspector `json:"-"`
	Fabric     runtime.Fabric    `json:"-"`
	Boundaries []ServiceBoundary `json:"boundaries"`
}

func (b Bundle) Boundary(name string) (ServiceBoundary, bool) {
	for _, item := range b.Boundaries {
		if item.Name == name {
			return item, true
		}
	}
	return ServiceBoundary{}, false
}
