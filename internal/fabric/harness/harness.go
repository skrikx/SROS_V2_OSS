package harness

import (
	"fmt"

	ctools "srosv2/contracts/tools"
)

type Harness struct {
	profiles map[string]Profile
}

func New() *Harness {
	return &Harness{profiles: Defaults()}
}

func (h *Harness) Profiles() map[string]Profile {
	out := map[string]Profile{}
	for k, v := range h.profiles {
		out[k] = v
	}
	return out
}

func (h *Harness) Resolve(manifest ctools.Manifest) (Resolution, error) {
	if len(manifest.SandboxProfiles) == 0 {
		return Resolution{}, fmt.Errorf("manifest %s declares no sandbox profiles", manifest.Name)
	}
	checked := make([]string, 0, len(manifest.SandboxProfiles))
	for _, name := range manifest.SandboxProfiles {
		checked = append(checked, name)
		profile, ok := h.profiles[name]
		if !ok {
			continue
		}
		if Compatible(profile, manifest.Class) {
			return Resolution{Compatible: true, Profile: profile, Checked: checked, Reason: "profile compatible"}, nil
		}
	}
	return Resolution{Compatible: false, Checked: checked, Reason: "no declared sandbox profile is compatible"}, nil
}
