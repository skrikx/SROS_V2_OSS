package registry

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	ctools "srosv2/contracts/tools"
)

func manifestPath(root, name string) string {
	return filepath.Join(root, cleanFileName(name)+".json")
}

func loadManifest(path string) (ctools.Manifest, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return ctools.Manifest{}, fmt.Errorf("read manifest %s: %w", path, err)
	}
	var manifest ctools.Manifest
	if err := json.Unmarshal(data, &manifest); err != nil {
		return ctools.Manifest{}, fmt.Errorf("decode manifest %s: %w", path, err)
	}
	return manifest, nil
}

func writeManifest(root string, manifest ctools.Manifest) error {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return fmt.Errorf("create manifest root: %w", err)
	}
	data, err := json.MarshalIndent(manifest, "", "  ")
	if err != nil {
		return fmt.Errorf("marshal manifest: %w", err)
	}
	return os.WriteFile(manifestPath(root, manifest.Name), append(data, '\n'), 0o644)
}

func listManifestFiles(root string) ([]string, error) {
	entries, err := os.ReadDir(root)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	files := []string{}
	for _, entry := range entries {
		if entry.IsDir() || filepath.Ext(entry.Name()) != ".json" {
			continue
		}
		files = append(files, filepath.Join(root, entry.Name()))
	}
	sort.Strings(files)
	return files, nil
}

func cleanFileName(v string) string {
	v = strings.ToLower(strings.TrimSpace(v))
	v = strings.ReplaceAll(v, " ", "_")
	v = strings.ReplaceAll(v, "/", "_")
	return v
}
