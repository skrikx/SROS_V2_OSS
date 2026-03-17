package harness

type Profile struct {
	Name          string `json:"name"`
	Description   string `json:"description"`
	AllowRead     bool   `json:"allow_read,omitempty"`
	AllowWrite    bool   `json:"allow_write,omitempty"`
	AllowPatch    bool   `json:"allow_patch,omitempty"`
	AllowShell    bool   `json:"allow_shell,omitempty"`
	AllowNetwork  bool   `json:"allow_network,omitempty"`
	BreakGlass    bool   `json:"break_glass,omitempty"`
	Containerized bool   `json:"containerized,omitempty"`
}

func Defaults() map[string]Profile {
	return map[string]Profile{
		"read_only": {
			Name:        "read_only",
			Description: "read-only local inspection",
			AllowRead:   true,
		},
		"patch_only": {
			Name:        "patch_only",
			Description: "filesystem patch operations without shell",
			AllowRead:   true,
			AllowWrite:  true,
			AllowPatch:  true,
		},
		"shell_gated": {
			Name:        "shell_gated",
			Description: "shell allowed only under governed gating",
			AllowRead:   true,
			AllowWrite:  true,
			AllowShell:  true,
		},
		"container_sandbox": {
			Name:          "container_sandbox",
			Description:   "containerized local execution",
			AllowRead:     true,
			AllowWrite:    true,
			AllowShell:    true,
			Containerized: true,
		},
		"network_restricted": {
			Name:        "network_restricted",
			Description: "local execution with outbound network disabled",
			AllowRead:   true,
			AllowWrite:  true,
		},
		"elevated_break_glass": {
			Name:         "elevated_break_glass",
			Description:  "elevated local execution with explicit break-glass posture",
			AllowRead:    true,
			AllowWrite:   true,
			AllowPatch:   true,
			AllowShell:   true,
			AllowNetwork: true,
			BreakGlass:   true,
		},
	}
}
