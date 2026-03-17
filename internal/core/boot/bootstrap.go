package boot

import (
	"fmt"
	"path/filepath"

	cmemory "srosv2/contracts/memory"
	"srosv2/contracts/policy"
	"srosv2/internal/core/gov"
	"srosv2/internal/core/mem"
	"srosv2/internal/core/mirror"
	"srosv2/internal/core/orch"
	coreprov "srosv2/internal/core/provenance"
	"srosv2/internal/core/runtime"
	"srosv2/internal/core/sr8"
	"srosv2/internal/core/sr9"
	coretrace "srosv2/internal/core/trace"
	"srosv2/internal/shared/config"
	"srosv2/internal/shared/ids"
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
		{Name: "trace", Wired: true, Summary: "append-only evidence lineage active"},
		{Name: "provenance", Wired: true, Summary: "receipts, bundles, and closure proofs active"},
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
	traceService, err := coretrace.New(filepath.Join(cfg.ArtifactRoot, "trace"), nil)
	if err != nil {
		return Bundle{}, err
	}
	provenanceService, err := coreprov.New(filepath.Join(cfg.ArtifactRoot, "provenance"), nil)
	if err != nil {
		return Bundle{}, err
	}
	orchestrator.SetEventHook(func(kind string, payload map[string]any) {
		runID, _ := payload["run_id"].(string)
		if runID == "" {
			return
		}
		_, _ = traceService.Emit(ids.RunID(runID), ids.TraceID("trace_"+runID), ids.SpanID(""), ids.SpanID(""), coretrace.EventWorkUnit, map[string]any{"hook_kind": kind, "payload": payload})
	})
	governor.SetDecisionHook(func(decision policy.PolicyDecision) {
		_, _ = traceService.Emit(decision.RunID, decision.TraceID, ids.SpanID(""), ids.SpanID(""), coretrace.EventPolicyDecision, map[string]any{"verdict": decision.Verdict, "capability": decision.Capability, "reason": decision.Reason})
	})
	memoryStore.SetMutationHook(func(mutation cmemory.MemoryMutation) {
		if mutation.RunID == "" {
			return
		}
		_, _ = traceService.Emit(mutation.RunID, ids.TraceID("trace_"+string(mutation.RunID)), ids.SpanID(""), ids.SpanID(""), coretrace.EventMemoryMutation, map[string]any{"mutation_id": mutation.MutationID, "lineage_ref": mutation.LineageRef, "key": mutation.Key})
	})
	mirrorEngine.SetWitnessHook(func(event mirror.WitnessEvent) {
		if event.RunID == "" {
			return
		}
		_, _ = traceService.Emit(event.RunID, ids.TraceID("trace_"+string(event.RunID)), ids.SpanID(""), ids.SpanID(""), coretrace.EventMirrorWitness, map[string]any{"witness_id": event.WitnessID, "severity": event.Severity, "basis": event.Basis})
	})

	runtimeManager, err := runtime.NewManager(runtime.Options{
		StoreDir:     filepath.Join(cfg.ArtifactRoot, "runtime"),
		Mode:         string(cfg.Mode),
		Gate:         sr9.NewGate(sr9.Options{}),
		Orchestrator: orchestrator,
		Governor:     governor,
		Memory:       memoryStore,
		Mirror:       mirrorEngine,
		Trace:        traceService,
		Provenance:   provenanceService,
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
		Trace:        traceService,
		Provenance:   provenanceService,
		Boundaries:   boundaries,
	}, nil
}
