package registry

import (
	"fmt"

	ctools "srosv2/contracts/tools"
)

var allowedTransitions = map[ctools.LifecycleState][]ctools.LifecycleState{
	ctools.StateDraft:        {ctools.StateValidated},
	ctools.StateValidated:    {ctools.StateAdmitted},
	ctools.StateAdmitted:     {ctools.StateExperimental, ctools.StateActive, ctools.StateQuarantined, ctools.StateDisabled},
	ctools.StateExperimental: {ctools.StateActive, ctools.StateQuarantined, ctools.StateDeprecated, ctools.StateDisabled},
	ctools.StateActive:       {ctools.StateQuarantined, ctools.StateDeprecated, ctools.StateDisabled},
	ctools.StateQuarantined:  {ctools.StateDisabled},
	ctools.StateDeprecated:   {ctools.StateDisabled},
	ctools.StateDisabled:     {},
}

func CanTransition(from, to ctools.LifecycleState) bool {
	for _, next := range allowedTransitions[from] {
		if next == to {
			return true
		}
	}
	return false
}

func Transition(manifest ctools.Manifest, to ctools.LifecycleState) (ctools.Manifest, error) {
	if !CanTransition(manifest.Status, to) {
		return ctools.Manifest{}, fmt.Errorf("invalid lifecycle transition %s -> %s", manifest.Status, to)
	}
	manifest.Status = to
	return manifest, nil
}
