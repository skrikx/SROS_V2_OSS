package boot

import "srosv2/internal/core/runtime"

type BundleBuilder interface {
	Build() (Bundle, error)
}

type ServiceBoundary struct {
	Name       string `json:"name"`
	Wired      bool   `json:"wired"`
	DeferredTo string `json:"deferred_to,omitempty"`
	Summary    string `json:"summary"`
}

type ServiceSet struct {
	Compiler  runtime.Compiler
	Runtime   runtime.Runtime
	Inspector runtime.Inspector
	Fabric    runtime.Fabric
}
