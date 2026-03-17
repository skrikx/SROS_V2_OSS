package boot

import (
	"fmt"
	"path/filepath"

	"srosv2/internal/core/gov"
	"srosv2/internal/core/mem"
	"srosv2/internal/core/mirror"
	"srosv2/internal/core/orch"
	"srosv2/internal/core/runtime"
	"srosv2/internal/core/sr8"
	"srosv2/internal/core/sr9"
	"srosv2/internal/shared/config"
)

func Bootstrap(cfg config.Config) (Bundle, error) {
	if cfg.Mode != config.ModeLocalCLI {
		return Bundle{}, fmt.Errorf("unsupported mode %q for v2 local bootstrap", cfg.Mode)
	}

	boundaries := []ServiceBoundary{
		{Name: "compiler", Wired: true, Summary: "sr8 compiler plane active"},
		{Name: "runtime", Wired: true, Summary: "sr9 runtime gate and state machine active"},
		{Name: "orch", Wired: true, Summary: "orchestration sequencing and checkpoint routing active"},
		{Name: "gov", Wired: true, Summary: "allow ask deny, sandbox, and permission gating active"},
		{Name: "mem", Wired: true, Summary: "workspace memory, lineage, and branch surfaces active"},
		{Name: "mirror", Wired: true, Summary: "semantic drift and witness surfaces active"},
		{Name: "inspector", Wired: true, Summary: "inspect routing boundary active for memory and mirror"},
		{Name: "fabric", Wired: true, DeferredTo: "W09", Summary: "governed fabric semantics only; execution deferred"},
	}

	policyPath := cfg.PolicyBundlePath
	if policyPath == "" {
		policyPath = filepath.Join("examples", "policy", "local_default_policy.json")
	}
	governor, err := gov.NewEngine(gov.Options{
		BundlePath:   policyPath,
		ArtifactRoot: cfg.ArtifactRoot,
	})
	if err != nil {
		return Bundle{}, err
	}
	orchestrator, err := orch.New(orch.Options{
		ArtifactRoot: filepath.Join(cfg.ArtifactRoot, "runtime", "orch"),
	})
	if err != nil {
		return Bundle{}, err
	}
	memoryStore, err := mem.NewStore(filepath.Join(cfg.ArtifactRoot, "memory"), nil)
	if err != nil {
		return Bundle{}, err
	}
	mirrorEngine, err := mirror.New(filepath.Join(cfg.ArtifactRoot, "mirror"), nil)
	if err != nil {
		return Bundle{}, err
	}

	runtimeManager, err := runtime.NewManager(runtime.Options{
		StoreDir:     filepath.Join(cfg.ArtifactRoot, "runtime"),
		Mode:         string(cfg.Mode),
		Gate:         sr9.NewGate(sr9.Options{}),
		Orchestrator: orchestrator,
		Governor:     governor,
		Memory:       memoryStore,
		Mirror:       mirrorEngine,
	})
	if err != nil {
		return Bundle{}, err
	}

	return Bundle{
		Mode:         ModeLocalCLI,
		Compiler:     sr8.NewCompiler(sr8.Options{}),
		Runtime:      runtimeManager,
		Inspector:    runtimeManager,
		Fabric:       runtimeManager,
		Orchestrator: orchestrator,
		Governor:     governor,
		Memory:       memoryStore,
		Mirror:       mirrorEngine,
		Boundaries:   boundaries,
	}, nil
}
