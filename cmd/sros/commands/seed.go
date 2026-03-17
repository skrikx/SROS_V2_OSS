package commands

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func newSeedCommand() *Command {
	return &Command{
		Name:    "seed",
		Summary: "Create a local seed artifact for smoke workflows",
		Usage:   "sros seed [--name <seed-name>]",
		Run: func(ctx *Context, args []string) error {
			fs := flag.NewFlagSet("seed", flag.ContinueOnError)
			fs.SetOutput(ioDiscard{})
			name := fs.String("name", "default", "seed file name")
			if err := fs.Parse(args); err != nil {
				return OperatorError(err.Error())
			}
			if fs.NArg() != 0 {
				return OperatorError("seed does not accept positional arguments")
			}

			seedDir := filepath.Join(ctx.Config.ArtifactRoot, "seeds")
			if err := os.MkdirAll(seedDir, 0o755); err != nil {
				return EnvironmentError(fmt.Sprintf("create seed directory: %v", err))
			}

			seedPath := filepath.Join(seedDir, *name+".txt")
			if err := os.WriteFile(seedPath, []byte("sros v2 seed artifact\n"), 0o644); err != nil {
				return EnvironmentError(fmt.Sprintf("write seed file: %v", err))
			}

			return writeOutput(ctx, "seed created: "+seedPath, map[string]any{"seed_path": seedPath})
		},
	}
}
