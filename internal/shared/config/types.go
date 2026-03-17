package config

type Mode string

const (
	ModeLocalCLI Mode = "local_cli"
)

type Config struct {
	Mode             Mode           `json:"mode"`
	WorkspaceRoot    string         `json:"workspace_root"`
	ArtifactRoot     string         `json:"artifact_root"`
	PolicyBundlePath string         `json:"policy_bundle_path"`
	MemoryStorePath  string         `json:"memory_store_path"`
	TraceStorePath   string         `json:"trace_store_path"`
	OutputFormat     string         `json:"output_format"`
	Database         DatabaseConfig `json:"database"`
}

type LoadOptions struct {
	CWD          string
	ExplicitPath string
	LookupEnv    func(string) string
}

type LoadResult struct {
	Config   Config   `json:"config"`
	Source   string   `json:"source"`
	Warnings []string `json:"warnings,omitempty"`
}
