package boot

import (
	"fmt"

	"srosv2/internal/shared/config"
)

func Bootstrap(cfg config.Config) (Bundle, error) {
	if cfg.Mode != config.ModeLocalCLI {
		return Bundle{}, fmt.Errorf("unsupported mode %q for v2 local bootstrap", cfg.Mode)
	}

	boundaries := []ServiceBoundary{
		{Name: "compiler", Wired: false, DeferredTo: "W04", Summary: "compile routing boundary only"},
		{Name: "runtime", Wired: false, DeferredTo: "W05", Summary: "runtime routing boundary only"},
		{Name: "inspector", Wired: false, DeferredTo: "W07/W08", Summary: "inspect routing boundary only"},
		{Name: "fabric", Wired: false, DeferredTo: "W09", Summary: "fabric routing boundary only"},
	}

	return Bundle{
		Mode:       ModeLocalCLI,
		Boundaries: boundaries,
	}, nil
}
