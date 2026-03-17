package boot

import (
	"context"

	"srosv2/internal/core/gov"
	"srosv2/internal/core/mem"
	"srosv2/internal/core/mirror"
	"srosv2/internal/core/orch"
	coreprov "srosv2/internal/core/provenance"
	"srosv2/internal/core/runtime"
	"srosv2/internal/core/sr8"
	coretrace "srosv2/internal/core/trace"
)

type BundleBuilder interface {
	Build() (Bundle, error)
}

type ServiceBoundary struct {
	Name       string `json:"name"`
	Wired      bool   `json:"wired"`
	DeferredTo string `json:"deferred_to,omitempty"`
	Summary    string `json:"summary"`
}

type Compiler interface {
	Compile(context.Context, sr8.CompileRequest) (sr8.CompileResult, error)
}

type ServiceSet struct {
	Compiler     Compiler
	Runtime      runtime.Runtime
	Inspector    runtime.Inspector
	Fabric       runtime.Fabric
	Orchestrator *orch.Orchestrator
	Governor     *gov.Engine
	Memory       *mem.Store
	Mirror       *mirror.Engine
	Trace        *coretrace.Service
	Provenance   *coreprov.Service
}
