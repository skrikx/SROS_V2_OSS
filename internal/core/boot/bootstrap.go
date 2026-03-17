package boot

import (
	"fmt"

	"srosv2/internal/core/sr8"
	"srosv2/internal/shared/config"
)

func Bootstrap(cfg config.Config) (Bundle, error) {
	if cfg.Mode != config.ModeLocalCLI {
		return Bundle{}, fmt.Errorf("unsupported mode %q for v2 local bootstrap", cfg.Mode)
	}

	boundaries := []ServiceBoundary{
		{Name: "compiler", Wired: true, Summary: "sr8 compiler plane active"},
		{Name: "runtime", Wired: false, DeferredTo: "W05", Summary: "runtime routing boundary only"},
		{Name: "inspector", Wired: false, DeferredTo: "W07/W08", Summary: "inspect routing boundary only"},
		{Name: "fabric", Wired: false, DeferredTo: "W09", Summary: "fabric routing boundary only"},
	}

	return Bundle{
		Mode:       ModeLocalCLI,
		Compiler:   sr8.NewCompiler(sr8.Options{}),
		Boundaries: boundaries,
	}, nil
}
