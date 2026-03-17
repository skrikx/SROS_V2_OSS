package boot

import (
	"fmt"
	"path/filepath"

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
		{Name: "inspector", Wired: false, DeferredTo: "W07/W08", Summary: "inspect routing boundary only"},
		{Name: "fabric", Wired: false, DeferredTo: "W09", Summary: "fabric routing boundary only"},
	}

	runtimeManager, err := runtime.NewManager(runtime.Options{
		StoreDir: filepath.Join(cfg.ArtifactRoot, "runtime"),
		Mode:     string(cfg.Mode),
		Gate:     sr9.NewGate(sr9.Options{}),
	})
	if err != nil {
		return Bundle{}, err
	}

	return Bundle{
		Mode:       ModeLocalCLI,
		Compiler:   sr8.NewCompiler(sr8.Options{}),
		Runtime:    runtimeManager,
		Boundaries: boundaries,
	}, nil
}
