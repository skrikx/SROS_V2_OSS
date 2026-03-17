package boot

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"srosv2/contracts/evidence"
	"srosv2/contracts/policy"
	ctools "srosv2/contracts/tools"
	"srosv2/internal/core/gov"
	coreprov "srosv2/internal/core/provenance"
	coretrace "srosv2/internal/core/trace"
	"srosv2/internal/fabric/connectors"
	"srosv2/internal/fabric/connectors/adapters"
	"srosv2/internal/fabric/harness"
	"srosv2/internal/fabric/mcpclient"
	"srosv2/internal/fabric/registry"
	"srosv2/internal/shared/ids"
)

type fabricService struct {
	registry   *registry.Registry
	harness    *harness.Harness
	connectors *connectors.Registry
	mcp        *mcpclient.Client
	governor   *gov.Engine
	trace      *coretrace.Service
	provenance *coreprov.Service
}

func newFabricService(artifactRoot string, governor *gov.Engine, trace *coretrace.Service, provenance *coreprov.Service) (*fabricService, error) {
	hr := harness.New()
	reg, err := registry.New(filepath.Join(artifactRoot, "fabric", "registry"), hr, builtinManifests())
	if err != nil {
		return nil, err
	}
	conn := connectors.NewRegistry(adapters.LocalFS{}, adapters.LocalHTTP{})
	return &fabricService{
		registry:   reg,
		harness:    hr,
		connectors: conn,
		mcp:        mcpclient.New(reg),
		governor:   governor,
		trace:      trace,
		provenance: provenance,
	}, nil
}

func builtinManifests() []ctools.Manifest {
	return []ctools.Manifest{
		{
			ManifestVersion: "v2.0",
			Name:            "local.patch",
			Title:           "Local Patch",
			Description:     "Governed local patch application",
			Version:         "1.0.0",
			Class:           "tool.local.patch",
			Domain:          "workspace",
			PolicyClass:     "patch.apply",
			Status:          ctools.StateActive,
			TrustBoundary:   string(policy.TrustBoundaryLocalFS),
			SandboxProfiles: []string{"patch_only"},
			Capabilities:    []string{"patch.apply"},
		},
		{
			ManifestVersion: "v2.0",
			Name:            "local.shell",
			Title:           "Local Shell",
			Description:     "Governed shell execution",
			Version:         "1.0.0",
			Class:           "tool.local.shell",
			Domain:          "workspace",
			PolicyClass:     "shell.exec",
			Status:          ctools.StateExperimental,
			TrustBoundary:   string(policy.TrustBoundaryLocalProcess),
			SandboxProfiles: []string{"shell_gated"},
			Capabilities:    []string{"shell.exec"},
			Unsafe:          true,
			Experimental:    true,
		},
		{
			ManifestVersion: "v2.0",
			Name:            "local.fileio",
			Title:           "Local File IO",
			Description:     "Governed file read and write",
			Version:         "1.0.0",
			Class:           "tool.local.fileio",
			Domain:          "workspace",
			PolicyClass:     "fileio.read",
			Status:          ctools.StateActive,
			TrustBoundary:   string(policy.TrustBoundaryLocalFS),
			SandboxProfiles: []string{"read_only", "patch_only"},
			Capabilities:    []string{"fileio.read", "fileio.write"},
		},
		{
			ManifestVersion:      "v2.0",
			Name:                 "connector.local_fs",
			Title:                "Local FS Connector",
			Description:          "Governed connector for local filesystem metadata",
			Version:              "1.0.0",
			Class:                "connector.local_fs",
			Domain:               "workspace",
			PolicyClass:          "connector.invoke",
			Status:               ctools.StateActive,
			TrustBoundary:        string(policy.TrustBoundaryLocalFS),
			SandboxProfiles:      []string{"read_only"},
			AuthType:             "local_token",
			AuthEnvelopeRequired: true,
			ConnectorRef:         "local_fs",
			Capabilities:         []string{"connector.invoke"},
		},
		{
			ManifestVersion:      "v2.0",
			Name:                 "connector.local_http",
			Title:                "Local HTTP Connector",
			Description:          "Governed connector for HTTP fetches",
			Version:              "1.0.0",
			Class:                "connector.local_http",
			Domain:               "network",
			PolicyClass:          "connector.invoke",
			Status:               ctools.StateExperimental,
			TrustBoundary:        string(policy.TrustBoundaryExternalNet),
			SandboxProfiles:      []string{"network_restricted"},
			AuthType:             "bearer",
			AuthEnvelopeRequired: true,
			ConnectorRef:         "local_http",
			Capabilities:         []string{"connector.invoke"},
			Experimental:         true,
		},
		{
			ManifestVersion:   "v2.0",
			Name:              "mcp.ingest",
			Title:             "MCP Ingest",
			Description:       "Governed MCP normalization and admission",
			Version:           "1.0.0",
			Class:             "mcp.ingress",
			Domain:            "integration",
			PolicyClass:       "mcp.ingest",
			Status:            ctools.StateActive,
			TrustBoundary:     string(policy.TrustBoundaryExternalNet),
			SandboxProfiles:   []string{"network_restricted"},
			MCPIngressCapable: true,
			Capabilities:      []string{"mcp.ingest"},
		},
	}
}

func (f *fabricService) ToolsList(context.Context) (map[string]any, error) {
	items := []ctools.ManifestSummary{}
	for _, manifest := range f.registry.List() {
		items = append(items, manifest.Summary())
	}
	return map[string]any{
		"registry_root": f.registry.Root(),
		"count":         len(items),
		"manifests":     items,
	}, nil
}

func (f *fabricService) ToolsShow(_ context.Context, name string) (map[string]any, error) {
	manifest, ok := f.registry.Get(name)
	if !ok {
		return nil, fmt.Errorf("manifest %s not found", name)
	}
	return map[string]any{"manifest": manifest}, nil
}

func (f *fabricService) ToolsSearch(_ context.Context, query ctools.SearchQuery) (map[string]any, error) {
	return map[string]any{"results": f.registry.Search(query)}, nil
}

func (f *fabricService) ToolsValidate(_ context.Context, path string) (map[string]any, error) {
	manifest, err := readManifest(path)
	if err != nil {
		return nil, err
	}
	result := f.registry.Validate(manifest)
	return map[string]any{"validation": result}, nil
}

func (f *fabricService) ToolsRegister(_ context.Context, path string) (map[string]any, error) {
	manifest, err := readManifest(path)
	if err != nil {
		return nil, err
	}
	registered, result, err := f.registry.Register(manifest)
	if err != nil {
		return map[string]any{"validation": result}, err
	}
	admitted, err := f.registry.Admit(registered.Name)
	if err != nil {
		return nil, err
	}
	return map[string]any{"registered": true, "manifest": admitted, "validation": result}, nil
}

func (f *fabricService) ConnectorsList(context.Context) (map[string]any, error) {
	return map[string]any{"connectors": f.connectors.List()}, nil
}

func (f *fabricService) ConnectorsInspectEnvelope(_ context.Context, path string) (map[string]any, error) {
	env, err := connectors.LoadEnvelope(path)
	if err != nil {
		return nil, err
	}
	_ = f.traceCapability("connector.envelope.inspect", map[string]any{"connector": env.Connector, "envelope_id": env.EnvelopeID})
	return connectors.RedactEnvelope(env), nil
}

func (f *fabricService) MCPIngest(_ context.Context, path string) (map[string]any, error) {
	result, err := f.mcp.Ingest(path)
	if err != nil {
		return nil, err
	}
	_ = f.traceCapability("mcp.ingest", map[string]any{"input": path})
	return result, nil
}

func (f *fabricService) traceCapability(capability string, payload map[string]any) error {
	if f.trace == nil {
		return nil
	}
	runID := ids.RunID("run_fabric_" + sanitize(capability))
	traceID := ids.TraceID("trace_fabric_" + sanitize(capability))
	_, err := f.trace.Emit(runID, traceID, ids.SpanID(""), ids.SpanID(""), coretrace.EventWorkUnit, map[string]any{
		"capability": capability,
		"payload":    payload,
	})
	if err != nil {
		return err
	}
	if f.provenance == nil {
		return nil
	}
	_, err = f.provenance.EmitReceipt(runID, evidence.ReceiptKindStage, "completed", "fabric capability invoked", nil, "")
	return err
}

func readManifest(path string) (ctools.Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ctools.Manifest{}, err
	}
	var manifest ctools.Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return ctools.Manifest{}, err
	}
	return manifest, nil
}

func sanitize(v string) string {
	v = strings.TrimSpace(v)
	v = strings.ReplaceAll(v, ".", "_")
	return strings.ReplaceAll(v, " ", "_")
}
