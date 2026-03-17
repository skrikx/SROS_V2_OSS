package commands

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
)

func newExamplesCommand() *Command {
	cmd := &Command{
		Name:    "examples",
		Summary: "Example catalog, killer examples, and showcase runner surfaces",
		Usage:   "sros examples <run|catalog>",
		Examples: []string{
			"sros examples run",
			"sros examples catalog",
		},
	}
	cmd.Subcommands = []*Command{
		{
			Name:    "run",
			Summary: "Load the example catalog and summarize the strongest runnable paths",
			Usage:   "sros examples run",
			Run: func(ctx *Context, args []string) error {
				if err := requireNoArgs(args); err != nil {
					return err
				}
				payload, err := loadExampleCatalog(filepath.Join(ctx.CWD, "examples", "catalog.json"))
				if err != nil {
					return EnvironmentError(err.Error())
				}
				text := "examples catalog loaded\nfocus: first-run, replay, fabric, and release-ready examples"
				return writeOutput(ctx, text, payload)
			},
		},
		{
			Name:    "catalog",
			Summary: "Show the example catalog with category counts and showcase paths",
			Usage:   "sros examples catalog",
			Run: func(ctx *Context, args []string) error {
				if err := requireNoArgs(args); err != nil {
					return err
				}
				payload, err := loadExampleCatalog(filepath.Join(ctx.CWD, "examples", "catalog.json"))
				if err != nil {
					return EnvironmentError(err.Error())
				}
				return writeOutput(ctx, "example catalog summarized", payload)
			},
		},
	}
	return cmd
}

func loadExampleCatalog(path string) (map[string]any, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var catalog map[string][]string
	if err := json.Unmarshal(data, &catalog); err != nil {
		return nil, err
	}
	keys := make([]string, 0, len(catalog))
	for key := range catalog {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	counts := map[string]int{}
	showcase := []string{}
	for _, key := range keys {
		counts[key] = len(catalog[key])
		if key == "runs" || key == "tools" || key == "traces" || key == "releases" {
			showcase = append(showcase, catalog[key]...)
		}
	}
	return map[string]any{
		"catalog_path": path,
		"counts":       counts,
		"showcase":     showcase,
		"catalog":      catalog,
	}, nil
}
